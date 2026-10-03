package cli

import (
	"strings"
	"testing"

	"github.com/edusouza/ctbz-cli/internal/ctbz"
)

func TestEmpresasList(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/dadosempresa/get": dadosEmpresaJSON}))

	out, stderr, code := execCLI(t, "", "empresas", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "atual,cnpj,razao_social,situacao\n" +
		"true,22222222000122,FULANO TECNOLOGIA LTDA,ATIVO\n" +
		"false,11111111000111,FULANO MEI,INATIVO\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	out, _, _ = execCLI(t, "", "empresas")
	if !strings.Contains(out, "11.111.111/0001-11") || !strings.Contains(out, "Atual") {
		t.Errorf("tabela:\n%s", out)
	}
}

func TestEmpresaUsar(t *testing.T) {
	t.Run("CNPJ inválido", func(t *testing.T) {
		t.Setenv("CTBZ_HOME", t.TempDir())
		if _, _, code := execCLI(t, "", "empresa", "usar", "123"); code != ExitUsage {
			t.Errorf("código %d, quero %d", code, ExitUsage)
		}
	})

	t.Run("empresa já é a atual", func(t *testing.T) {
		store := &ctbz.Store{Dir: t.TempDir()}
		t.Setenv("CTBZ_HOME", store.Dir)
		store.SaveSession(&ctbz.Session{BaseURL: "http://nao-deve-ser-chamado.invalid", CNPJ: "22222222000122"})
		_, stderr, code := execCLI(t, "", "empresa", "usar", "22.222.222/0001-22")
		if code != ExitOK || !strings.Contains(stderr, "já está") {
			t.Errorf("código %d, stderr %q", code, stderr)
		}
	})

	t.Run("troca refazendo o login", func(t *testing.T) {
		srv := fakeContabilizei(t)
		store := &ctbz.Store{Dir: t.TempDir()}
		t.Setenv("CTBZ_HOME", store.Dir)
		t.Setenv("CTBZ_BASE_URL", srv.URL)
		t.Setenv("CTBZ_USER", "00000000000")
		t.Setenv("CTBZ_PASSWORD", "x")
		t.Setenv("CTBZ_OTP_CMD", "echo 123456")
		store.SaveSession(&ctbz.Session{BaseURL: srv.URL, CNPJ: "11111111000111"})

		_, stderr, code := execCLI(t, "", "empresa", "usar", "22222222000122")
		if code != ExitOK {
			t.Fatalf("código %d: %s", code, stderr)
		}
		sess, err := store.LoadSession()
		if err != nil || sess.CNPJ != "22222222000122" {
			t.Errorf("sessão não trocou de empresa: %+v %v", sess, err)
		}
	})
}
