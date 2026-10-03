package cli

import (
	"strings"
	"testing"
)

func TestLucrosInforme(t *testing.T) {
	fixNow(t, "2026-10-03")
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/informerendimento/listSocioInformeRendimentos/2025": `[{"id":7,"nome":"FULANO"},{"id":"8","nome":"BELTRANO"}]`,
		"/api/plataforma/informerendimento/getValoresComprovanteRendimento/7/2025": `{"nome":"FULANO","cpf":"12345678901",
"rendimentos":19452,"previdencia":2139.72,"irrfRetido":0,"decimoTerceiro":0,"irrfdecimoTerceiro":0,"lucro":50000.5}`,
		"/api/plataforma/informerendimento/getValoresComprovanteRendimento/8/2025": `{"nome":"","cpf":"",
"rendimentos":"R$ 1.000,00","previdencia":"110.00","irrfRetido":null,"lucro":"-"}`,
	}))
	out, stderr, code := execCLI(t, "", "lucros", "informe", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "ano,socio,cpf,rendimentos,previdencia,irrf_retido,decimo_terceiro,irrf_decimo_terceiro,lucros_isentos\n" +
		"2025,FULANO,12345678901,19452.00,2139.72,0.00,0.00,0.00,50000.50\n" +
		"2025,BELTRANO,,1000.00,110.00,,,,-\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
}

func TestLucrosInformeVazio(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/informerendimento/listSocioInformeRendimentos/2024": fixture(t, "socios_informe"),
	}))
	out, stderr, code := execCLI(t, "", "lucros", "informe", "--ano", "2024")
	if code != ExitOK || !strings.Contains(stderr, "Nenhum sócio com informe de rendimentos em 2024") {
		t.Errorf("código %d, stdout %q, stderr %q", code, out, stderr)
	}
}
