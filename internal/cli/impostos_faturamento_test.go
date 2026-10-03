package cli

import (
	"strings"
	"testing"
)

func TestImpostosFaturamento(t *testing.T) {
	calculo := `{"nomeMesCompetencia":"Setembro","faturamentoTotal":10,"dasSimples":{"inconsistente":false},"darf":{"inss":{},"irrf":{}},
"valorFaturamentoUltimos12Meses":36000,"valorProlaboreUltimos12Meses":10000,"percentualFatorR":27.78,
"historicoFaturamento":[{"mes":"jan./27","valorFaturamento":3000,"valorProlabore":1000},{"mes":"dez./26","valorFaturamento":2000,"valorProlabore":null}]}`
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/impostos/como-imposto-foi-calculado/init":              calculo,
		"/api/plataforma/impostos/v5/historico-impostos/dados-grafico?ano=2026": `{"meses":{"12":{"totalPago":150.5}}}`,
		"/api/plataforma/impostos/v5/historico-impostos/dados-grafico?ano=2027": fixture(t, "impostos_pagos_no_ano"),
	}))
	out, stderr, code := execCLI(t, "", "impostos", "faturamento", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "competencia,faturamento,prolabore,impostos_pagos\n" +
		"12/2026,2000.00,,150.50\n" +
		"01/2027,3000.00,1000.00,1000.00\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	_, stderr, _ = execCLI(t, "", "impostos", "faturamento")
	if !strings.Contains(stderr, "RBT12): R$ 36.000,00") || !strings.Contains(stderr, "Fator R: 27.78%") {
		t.Errorf("resumo: %q", stderr)
	}
}
