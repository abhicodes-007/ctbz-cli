package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "api", "testdata", name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestEmpresaSocios(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/legado/socio/list": fixture(t, "socios")}))
	out, stderr, code := execCLI(t, "", "empresa", "socios", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "nome,cpf,administrador,responsavel_receita,categoria,salario_base,entrada,situacao\n" +
		"FULANO DE TAL,00000000000,true,true,Sócio com Prolabore,1000.00,2026-07-01,Ativo\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	out, _, _ = execCLI(t, "", "empresa", "socios")
	if !strings.Contains(out, "000.000.000-00") || !strings.Contains(out, "R$ 1.000,00") {
		t.Errorf("tabela:\n%s", out)
	}
}

func TestEmpresaAtividades(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/legado/notafiscal/cnaeanexosmultiplos/list": fixture(t, "cnaes")}))
	out, stderr, code := execCLI(t, "", "empresa", "atividades", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	for _, want := range []string{
		"cnae,descricao,principal,ramo,anexos_simples\n",
		"6209-1/00,\"Suporte técnico, manutenção e outros serviços em tecnologia da informação\",true,servico,V; III\n",
		"6319-4/00,",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("CSV sem %q:\n%s", want, out)
		}
	}
}

func TestFormatHelpers(t *testing.T) {
	if formatCNAE("6209100") != "6209-1/00" || formatCNAE("123") != "123" {
		t.Error("formatCNAE")
	}
	if anexoRomano(3) != "III" || anexoRomano(9) != "?" {
		t.Error("anexoRomano")
	}
}
