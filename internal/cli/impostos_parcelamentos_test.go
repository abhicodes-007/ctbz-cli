package cli

import (
	"strings"
	"testing"
)

func TestImpostosParcelamentos(t *testing.T) {
	// Exemplo sintético com os campos vistos no front (a fixture real tem listas vazias).
	comItens := `{"abaParcelamentos":{"emAndamento":[],"ativos":[{"idParcelamento":7,"titulo":"Simples Nacional",
"tipoParcelamento":"SIMPLES_NACIONAL","status":"ATIVO","parcelaAtual":{"numeroParcela":3}}],
"historico":[{"idParcelamento":8,"titulo":"PGFN","tipoParcelamento":"PGFN_PREVIDENCIARIO","status":"QUITADO","parcelaAtual":null}]}}`
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/impostos/v3/impostos-a-pagar/init":            comItens,
		"/api/plataforma/impostos/parcelamento/detalhes/7":             `{"descricao":"Parcelamento do Simples","saldoDevedor":1234.5,"totalParcelas":60}`,
		"/api/plataforma/dashboard/informerendimento/debitos-federais": `{"possuiDebitosFederais":true}`,
	}))

	out, stderr, code := execCLI(t, "", "impostos", "parcelamentos", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "grupo,id,titulo,tipo,situacao,parcela_atual\n" +
		"ativo,7,Simples Nacional,simples_nacional,ATIVO,3\n" +
		"historico,8,PGFN,pgfn_previdenciario,QUITADO,\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}

	out, _, code = execCLI(t, "", "impostos", "parcelamento", "7", "-o", "json")
	if code != ExitOK || !strings.Contains(out, `"saldoDevedor": 1234.5`) {
		t.Errorf("detalhe (código %d):\n%s", code, out)
	}

	out, _, code = execCLI(t, "", "impostos", "debitos", "--fail-on-debitos", "-o", "csv")
	if code != ExitAttention || out != "possui_debitos_federais\ntrue\n" {
		t.Errorf("débitos (código %d):\n%s", code, out)
	}
}

func TestImpostosParcelamentosVazio(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/impostos/v3/impostos-a-pagar/init":            fixture(t, "parcelamentos"),
		"/api/plataforma/dashboard/informerendimento/debitos-federais": fixture(t, "debitos_federais"),
	}))
	out, _, code := execCLI(t, "", "impostos", "parcelamentos", "-o", "json")
	if code != ExitOK || out != "[]\n" {
		t.Errorf("sem parcelamentos (código %d): %q", code, out)
	}
	if _, _, code := execCLI(t, "", "impostos", "debitos", "--fail-on-debitos"); code != ExitOK {
		t.Errorf("sem débitos deveria sair com 0, veio %d", code)
	}
}
