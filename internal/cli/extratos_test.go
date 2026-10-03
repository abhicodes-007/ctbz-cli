package cli

import "testing"

func TestExtratos(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/movimentacao-financeira/v2/extratos": `[
{"ano":2026,"mes":10,"idContaBancaria":7,"banco":"Contabilizei Conta PJ","agencia":"0001","numeroConta":"12345","situacao":"ABERTO","statusIntegracao":null},
{"ano":2025,"mes":12,"idContaBancaria":7,"banco":"Contabilizei Conta PJ","agencia":"0001","numeroConta":"12345","situacao":"FECHADO","statusIntegracao":"INTEGRADO"}]`,
		"/api/plataforma/contabancaria/list": `{"bancos":[{"id":1}],"contasBancarias":[{"id":7,"nomeBanco":"Contabilizei","codigoBanco":"301-1",
"agencia":"0001","contaCorrente":"12345","vlrSaldoInicial":0,"dataSaldoInicial":1784635200000,"statusIntegracao":"INTEGRADA","fluxoIntegracao":"CONTABILIZEI_BANK"}]}`,
	}))
	out, stderr, code := execCLI(t, "", "extratos", "--ano", "2026", "-o", "csv")
	want := "competencia,banco,agencia,conta,situacao,integracao,id_conta\n10/2026,Contabilizei Conta PJ,0001,12345,ABERTO,,7\n"
	if code != ExitOK || out != want {
		t.Errorf("extratos (código %d, %s):\n%s\nesperado:\n%s", code, stderr, out, want)
	}
	out, stderr, code = execCLI(t, "", "contas-bancarias", "-o", "csv")
	want = "banco,codigo_banco,agencia,conta,saldo_inicial,data_saldo_inicial,integracao,fluxo_integracao,id\n" +
		"Contabilizei,301-1,0001,12345,0.00,2026-07-21,INTEGRADA,CONTABILIZEI_BANK,7\n"
	if code != ExitOK || out != want {
		t.Errorf("contas (código %d, %s):\n%s\nesperado:\n%s", code, stderr, out, want)
	}
}

func TestExtratosFixtures(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/movimentacao-financeira/v2/extratos": fixture(t, "extratos"),
		"/api/plataforma/contabancaria/list":                  fixture(t, "contas_bancarias"),
	}))
	for _, cmd := range []string{"extratos", "contas-bancarias"} {
		if _, stderr, code := execCLI(t, "", cmd, "-o", "json"); code != ExitOK {
			t.Errorf("%s: código %d: %s", cmd, code, stderr)
		}
	}
}
