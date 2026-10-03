package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/edusouza/ctbz-cli/internal/ctbz"
)

func TestParseGlobalFlags(t *testing.T) {
	for _, tc := range []struct {
		args       []string
		wantOutput string
		wantRest   string
	}{
		{[]string{"-o", "json", "status"}, "json", "status"},
		{[]string{"--output=csv", "empresa", "-o", "table"}, "csv", "empresa -o table"},
		{[]string{"status"}, "", "status"},
	} {
		globalOutput = ""
		rest, err := parseGlobalFlags(tc.args)
		if err != nil {
			t.Fatal(err)
		}
		if globalOutput != tc.wantOutput || strings.Join(rest, " ") != tc.wantRest {
			t.Errorf("%v → output=%q rest=%v", tc.args, globalOutput, rest)
		}
	}
	globalOutput = ""
	if _, err := parseGlobalFlags([]string{"-o"}); err == nil {
		t.Error("esperava erro para -o sem valor")
	}
}

func TestDefaultFormatPrecedence(t *testing.T) {
	t.Cleanup(func() { globalOutput = "" })
	globalOutput = ""
	t.Setenv("CTBZ_OUTPUT", "")
	if got := defaultFormat(); got != "table" {
		t.Errorf("padrão = %q", got)
	}
	t.Setenv("CTBZ_OUTPUT", "csv")
	if got := defaultFormat(); got != "csv" {
		t.Errorf("com CTBZ_OUTPUT = %q", got)
	}
	if got := explicitFormatOr("json"); got != "json" {
		t.Errorf("api deveria ignorar CTBZ_OUTPUT, veio %q", got)
	}
	globalOutput = "table"
	if got := defaultFormat(); got != "table" {
		t.Errorf("-o global deveria vencer CTBZ_OUTPUT, veio %q", got)
	}
}

// capture executa fn trocando os.Stdout e os.Stderr por pipes.
func capture(t *testing.T, fn func() error) (stdout, stderr string, err error) {
	t.Helper()
	oldOut, oldErr := os.Stdout, os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout, os.Stderr = wOut, wErr
	outC, errC := make(chan string), make(chan string)
	go func() { b, _ := io.ReadAll(rOut); outC <- string(b) }()
	go func() { b, _ := io.ReadAll(rErr); errC <- string(b) }()
	err = fn()
	wOut.Close()
	wErr.Close()
	os.Stdout, os.Stderr = oldOut, oldErr
	return <-outC, <-errC, err
}

func TestEmpresaFormats(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/plataforma/dadosempresa/get" {
			http.NotFound(w, r)
			return
		}
		io.WriteString(w, `{"empresaAtual":{"cnpj":"22.222.222/0001-22","razaoSocial":"FULANO TECNOLOGIA LTDA",
			"regimeTributario":"Simples Nacional","statusEmpresa":"ATIVO","ramosAtividade":[],"plano":"Básico",
			"certificado":{"status":{"descricao":"Ativo"},"dataValidade":"22/07/2027"}},
			"empresas":[{"cnpj":"22222222000122","razaoSocial":"FULANO TECNOLOGIA LTDA","statusEmpresa":"ATIVO"},
			            {"cnpj":"11111111000111","razaoSocial":"FULANO MEI","statusEmpresa":"INATIVO"}]}`)
	}))
	t.Cleanup(srv.Close)
	store := &ctbz.Store{Dir: t.TempDir()}
	t.Setenv("CTBZ_HOME", store.Dir)
	t.Setenv("CTBZ_OUTPUT", "")
	if err := store.SaveSession(&ctbz.Session{BaseURL: srv.URL, CreatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// --json é atalho para -o json.
	for _, args := range [][]string{{"--json"}, {"-o", "json"}, {"--output=json"}} {
		out, errOut, err := capture(t, func() error { return cmdEmpresa(ctx, args) })
		if err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		var v struct {
			CNPJ                string `json:"cnpj"`
			CertificadoValidade string `json:"certificado_validade"`
			OutrasEmpresas      []struct {
				CNPJ string `json:"cnpj"`
			} `json:"outras_empresas"`
		}
		if err := json.Unmarshal([]byte(out), &v); err != nil {
			t.Fatalf("%v: stdout não é JSON: %v\n%s", args, err, out)
		}
		if v.CNPJ != "22222222000122" || v.CertificadoValidade != "2027-07-22" || len(v.OutrasEmpresas) != 1 {
			t.Errorf("%v: JSON inesperado: %s", args, out)
		}
		if errOut != "" {
			t.Errorf("%v: stderr deveria estar vazio, veio %q", args, errOut)
		}
	}

	out, _, err := capture(t, func() error { return cmdEmpresa(ctx, []string{"-o", "csv"}) })
	if err != nil {
		t.Fatal(err)
	}
	if lines := strings.Split(strings.TrimSpace(out), "\n"); len(lines) != 2 || !strings.HasPrefix(lines[0], "razao_social,cnpj,") {
		t.Errorf("CSV inesperado:\n%s", out)
	}

	out, _, err = capture(t, func() error { return cmdEmpresa(ctx, nil) })
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"CNPJ:", "22.222.222/0001-22", "22/07/2027", "Outras empresas:", "FULANO MEI"} {
		if !strings.Contains(out, want) {
			t.Errorf("tabela sem %q:\n%s", want, out)
		}
	}

	if _, _, err := capture(t, func() error { return cmdEmpresa(ctx, []string{"-o", "xml"}) }); err == nil {
		t.Error("esperava erro para formato inválido")
	}
}
