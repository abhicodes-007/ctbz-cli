package cli

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/edusouza/ctbz-cli/internal/ctbz"
)

// fakeContabilizei responde com telas mínimas do fluxo real de login.
func fakeContabilizei(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(contabilizeiMux())
	t.Cleanup(srv.Close)
	return srv
}

// contabilizeiMux imita as telas de login (usuário/senha, OTP 123456, seleção de empresa).
func contabilizeiMux() *http.ServeMux {
	storage := base64.StdEncoding.EncodeToString([]byte(url.PathEscape(
		`{"email":"fulano@example.com","empresa":{"cnpj":"22222222000122","razaoSocial":"FULANO LTDA"}}`)))
	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		switch {
		case r.Method == http.MethodGet:
			io.WriteString(w, `<form id="form-login"><input type="hidden" name="token" value="T"/></form>`)
		case string(body) == "v=2&otp=123456":
			w.WriteHeader(http.StatusPartialContent)
			io.WriteString(w, `<input type="radio" name="cnpj" value="11111111000111"><div><span>A</span></div>
<input type="radio" name="cnpj" value="22222222000122"><div><span>FULANO LTDA</span></div>
<input type="hidden" name="token" value="T" />`)
		case strings.HasPrefix(string(body), "v="):
			w.WriteHeader(http.StatusNotFound)
		default:
			http.SetCookie(w, &http.Cookie{Name: "__C", Value: "c", Path: "/"})
			io.WriteString(w, `seu e-mail: <strong>fu***@ex***com</strong><input id="otp">`)
		}
	})
	mux.HandleFunc("/selecionarempresa", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "oauth-token", Value: "ok", Path: "/"})
		io.WriteString(w, `<script>localStorage.setItem("l","`+storage+`");</script>`)
	})
	return mux
}

func TestLoginWithOTPCommand(t *testing.T) {
	srv := fakeContabilizei(t)
	t.Setenv("CTBZ_BASE_URL", srv.URL)
	t.Setenv("CTBZ_USER", "00000000000")
	t.Setenv("CTBZ_PASSWORD", "x")
	store := &ctbz.Store{Dir: t.TempDir()}

	sess, err := login(context.Background(), store, testStreams(), loginOpts{
		otpCmd:     `echo "Seu código: 123456"`,
		otpTimeout: 5 * time.Second,
		cnpj:       "22.222.222/0001-22",
		quiet:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if sess.CNPJ != "22222222000122" || sess.Cookies["oauth-token"] == nil {
		t.Fatalf("sessão inesperada: %+v", sess)
	}
	if _, err := store.LoadSession(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.LoadPending(); err == nil {
		t.Fatal("login pendente deveria ter sido removido")
	}
}

func TestLoginPendingThenResume(t *testing.T) {
	srv := fakeContabilizei(t)
	t.Setenv("CTBZ_BASE_URL", srv.URL)
	t.Setenv("CTBZ_USER", "00000000000")
	t.Setenv("CTBZ_PASSWORD", "x")
	store := &ctbz.Store{Dir: t.TempDir()}
	ctx := context.Background()

	// Sem OTP e sem terminal: o login fica pendente.
	if _, err := login(ctx, store, testStreams(), loginOpts{quiet: true}); !errors.Is(err, errPending) {
		t.Fatalf("esperava errPending, veio %v", err)
	}
	// OTP certo, mas há duas empresas e nenhuma escolhida: continua pendente.
	if _, err := login(ctx, store, testStreams(), loginOpts{otpCode: "123456", quiet: true}); !errors.Is(err, errPending) {
		t.Fatalf("esperava errPending na seleção de empresa, veio %v", err)
	}
	if p, err := store.LoadPending(); err != nil || p.State.Stage != ctbz.StageSelectCompany {
		t.Fatalf("pendência inesperada: %+v %v", p, err)
	}
	sess, err := login(ctx, store, testStreams(), loginOpts{cnpj: "22222222000122", quiet: true})
	if err != nil {
		t.Fatal(err)
	}
	if sess.CNPJ != "22222222000122" {
		t.Fatalf("CNPJ = %q", sess.CNPJ)
	}
}

func testStreams() streams {
	return streams{in: strings.NewReader(""), out: io.Discard, err: io.Discard}
}

func TestReloginUnicoComChamadasEmParalelo(t *testing.T) {
	mux := contabilizeiMux()
	mux.HandleFunc("/api/plataforma/", func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie("oauth-token"); err != nil || c.Value != "ok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		io.WriteString(w, `{}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	withSession(t, srv) // sessão sem cookies: toda chamada recebe 401 até o re-login
	store := &ctbz.Store{Dir: os.Getenv("CTBZ_HOME")}
	sess, err := store.LoadSession()
	if err != nil {
		t.Fatal(err)
	}
	sess.CNPJ = "22222222000122" // o re-login escolhe a mesma empresa
	if err := store.SaveSession(sess); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CTBZ_BASE_URL", srv.URL)
	t.Setenv("CTBZ_USER", "00000000000")
	t.Setenv("CTBZ_PASSWORD", "x")
	execucoes := filepath.Join(t.TempDir(), "otp")
	t.Setenv("CTBZ_OTP_CMD", "echo >> "+execucoes+"; echo 123456")

	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var v map[string]any
			if err := getJSON(context.Background(), testStreams(), "appbar/get", &v); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	data, _ := os.ReadFile(execucoes)
	if n := strings.Count(string(data), "\n"); n != 1 {
		t.Errorf("CTBZ_OTP_CMD rodou %d vezes, quero 1 (um único re-login)", n)
	}
}
