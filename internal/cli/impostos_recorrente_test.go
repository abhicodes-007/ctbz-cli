package cli

import (
	"strings"
	"testing"
)

func TestImpostosRecorrente(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/payments/recorrencia/init": fixture(t, "recorrencia"),
		"/api/plataforma/payments/recorrencia/historico": `{"pagamentos":[{"mes":9,"ano":2026,"custoOperacao":1.99,
"guias":[{"nome":"DAS","status":"PAGO","valor":150.5},{"nome":"INSS","status":"RECUSADO","valor":null}]},
{"mes":"10","ano":"2026","guias":[]}]}`,
	}))

	out, stderr, code := execCLI(t, "", "impostos", "recorrente", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "disponivel,ativo,situacao,competencia,proximo_pagamento,proxima_tentativa,cartoes_salvos," +
		"pagamentos_agendados,pagamentos_concluidos,pagamentos_recusados\n" +
		"true,false,,09/2026,,,0,0,0,0\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}

	out, _, code = execCLI(t, "", "impostos", "recorrente", "historico", "-o", "csv")
	want = "competencia,item,situacao,valor\n" +
		"09/2026,DAS,PAGO,150.50\n" +
		"09/2026,INSS,RECUSADO,\n" +
		"09/2026,Custo de operação,,1.99\n"
	if code != ExitOK || out != want {
		t.Errorf("histórico (código %d):\n%s\nesperado:\n%s", code, out, want)
	}
}

func TestImpostosRecorrenteNaoMostraChaves(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/payments/recorrencia/init": `{"habilitado":true,"ativado":true,"status":"ATIVO","competencia":"10-2026",
"adyenClientKey":{"https://app.contabilizei.com.br":"live_SEGREDO"},"cartoesSalvos":[{"final":"1234","bandeira":"visa"}],
"pagamentosAgendados":[],"pagamentosConcluido":[],"pagamentosRecusado":[]}`,
	}))
	out, _, code := execCLI(t, "", "impostos", "recorrente", "-o", "json")
	if code != ExitOK || strings.Contains(out, "SEGREDO") || strings.Contains(out, "1234") || !strings.Contains(out, `"cartoes_salvos": 1`) {
		t.Errorf("código %d:\n%s", code, out)
	}
}
