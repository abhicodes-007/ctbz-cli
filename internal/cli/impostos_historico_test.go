package cli

import (
	"strings"
	"testing"
)

func TestImpostosHistorico(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/impostos/v2/historico-impostos/guias?pagina=1":                            fixture(t, "historico_guias"),
		"/api/plataforma/impostos/v2/historico-impostos/init":                                      fixture(t, "historico_resumo"),
		"/api/plataforma/impostos/v2/historico-impostos/guias?ano=2026&mes=7&pagina=1&status=PAGO": `{"paginaAtual":1,"totalPaginas":2,"competencias":[{"competencia":"Julho / 2026","guias":[{"id":1,"imposto":"DAS","impostoDescricao":"","competencia":{"mes":7,"ano":2026},"dataVencimento":"20/08/2026","valorPrincipal":10,"valorPago":10,"status":"PAGO","tipo":"GUIA"}]}]}`,
		"/api/plataforma/impostos/v2/historico-impostos/guias?ano=2026&mes=7&pagina=2&status=PAGO": `{"paginaAtual":2,"totalPaginas":2,"competencias":[{"competencia":"Julho / 2026","guias":[{"id":2,"imposto":"DARF","impostoDescricao":"DARF Unificado","competencia":{"mes":7,"ano":2026},"dataVencimento":"20/08/2026","valorPrincipal":20,"valorPago":null,"status":"PAGO","tipo":"GUIA"}]}]}`,
	}))

	out, stderr, code := execCLI(t, "", "impostos", "historico", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "competencia,id,imposto,vencimento,valor,valor_pago,situacao,tipo\n" +
		"07/2026,1000000000000001,DARF Unificado,2026-10-06,1000.00,1234.56,PENDENTE,guia\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	_, stderr, _ = execCLI(t, "", "impostos", "historico")
	if !strings.Contains(stderr, "Em dia: sim · Guias vencidas: 0") {
		t.Errorf("resumo no stderr: %q", stderr)
	}

	// Filtros e paginação: duas páginas, status em minúsculas é normalizado.
	out, stderr, code = execCLI(t, "", "impostos", "historico", "--ano", "2026", "--mes", "7", "--status", "pago", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	if !strings.Contains(out, "07/2026,1,DAS,") || !strings.Contains(out, "07/2026,2,DARF Unificado,2026-08-20,20.00,,PAGO") {
		t.Errorf("paginação/filtros:\n%s", out)
	}
	if _, _, code := execCLI(t, "", "impostos", "historico", "--mes", "13"); code != ExitUsage {
		t.Errorf("--mes 13: código %d", code)
	}
}
