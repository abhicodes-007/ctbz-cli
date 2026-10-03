package cli

import "testing"

func TestContas(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/movimentacao-financeira/contasUsuario": `[
{"id":1,"descricao":"Aluguel","descricaoContaContabil":"Aluguéis e Arrendamentos","classificacao":"DESPESA","situacao":"ATIVO"},
{"id":2,"descricao":"Aluguel de máquinas","descricaoContaContabil":"Aluguéis e Arrendamentos","classificacao":"DESPESA","situacao":"INATIVO"},
{"id":3,"descricao":"Prestação de serviços","descricaoContaContabil":"Receita de Serviços","classificacao":"RECEITA","situacao":"ATIVO"}]`}))

	out, stderr, code := execCLI(t, "", "contas", "--busca", "ALUG", "--situacao", "ativo", "-o", "csv")
	want := "descricao,conta_contabil,classificacao,situacao,id\nAluguel,Aluguéis e Arrendamentos,DESPESA,ATIVO,1\n"
	if code != ExitOK || out != want {
		t.Errorf("código %d (%s):\n%s\nesperado:\n%s", code, stderr, out, want)
	}
	if out, _, _ := execCLI(t, "", "contas", "--busca", "receita", "-o", "csv"); out != "descricao,conta_contabil,classificacao,situacao,id\nPrestação de serviços,Receita de Serviços,RECEITA,ATIVO,3\n" {
		t.Errorf("busca na conta contábil:\n%s", out)
	}
	if _, _, code := execCLI(t, "", "contas", "--situacao", "todas"); code != ExitUsage {
		t.Errorf("situação inválida: código %d, quero %d", code, ExitUsage)
	}
}

func TestContasFixture(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/movimentacao-financeira/contasUsuario": fixture(t, "contas_usuario")}))
	if _, stderr, code := execCLI(t, "", "contas", "-o", "json"); code != ExitOK {
		t.Errorf("código %d: %s", code, stderr)
	}
}
