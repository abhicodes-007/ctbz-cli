package cli

import (
	"strings"
	"testing"
)

func TestBalancete(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/relatorios-ms/gerar-balancete/2026/8": `[
{"id":"1","idContaPai":null,"descricao":"ATIVO","nivel":1,"tipo":"S","saldoAnterior":0,"totalDebito":1815,"totalCredito":1419.83,"saldoExercicio":395.17},
{"id":"1.01.01.01.00","idContaPai":"1.01.01.01","descricao":"Caixa Geral","nivel":5,"tipo":"A","saldoAnterior":0,"totalDebito":1139,"totalCredito":1280.83,"saldoExercicio":-141.83}]`,
		"/api/plataforma/relatorios-ms/gerar-balancete/2026/9": fixture(t, "balancete"),
	}))
	out, stderr, code := execCLI(t, "", "balancete", "2026-08", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "conta,descricao,nivel,saldo_anterior,debitos,creditos,saldo\n" +
		"1,ATIVO,1,0.00,1815.00,1419.83,395.17\n" +
		"1.01.01.01.00,Caixa Geral,5,0.00,1139.00,1280.83,-141.83\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	out, _, _ = execCLI(t, "", "balancete", "2026-08")
	if !strings.Contains(out, "        Caixa Geral") { // nível 5: quatro recuos de dois espaços
		t.Errorf("tabela deveria recuar a conta de nível 5:\n%s", out)
	}
	if out, _, code := execCLI(t, "", "balancete", "2026-09", "-o", "json"); code != ExitOK || !strings.Contains(out, `"descricao": "CIRCULANTE"`) {
		t.Errorf("fixture (código %d):\n%s", code, out)
	}
	if _, _, code := execCLI(t, "", "balancete", "08/2026"); code != ExitUsage {
		t.Errorf("mês inválido: código %d, quero %d", code, ExitUsage)
	}
}

func TestBalanco(t *testing.T) {
	balanco := `[
{"id":"1","descricao":"ATIVO","nivel":1,"classificacaoConta":"ATIVO","saldoExercicio":395.17,"saldoExercicioAnterior":0},
{"id":"2","descricao":"PASSIVO","nivel":1,"classificacaoConta":"PASSIVO","saldoExercicio":-395.17,"saldoExercicioAnterior":0},
{"id":"3","descricao":"RECEITAS","nivel":1,"classificacaoConta":"RESULTADO","saldoExercicio":100,"saldoExercicioAnterior":0}]`
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/relatorios-ms/gerarbalanco/2025/12": balanco,
		"/api/plataforma/relatorios-ms/gerarbalanco/2026/9":  fixture(t, "balanco"),
	}))
	out, stderr, code := execCLI(t, "", "balanco", "2025", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "conta,descricao,nivel,grupo,saldo,saldo_exercicio_anterior\n" +
		"1,ATIVO,1,ATIVO,395.17,0.00\n" +
		"2,PASSIVO,1,PASSIVO,-395.17,0.00\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	if out, _, code := execCLI(t, "", "balanco", "2026-09", "-o", "json"); code != ExitOK || !strings.Contains(out, `"grupo": "ATIVO"`) {
		t.Errorf("fixture (código %d):\n%s", code, out)
	}
	if _, _, code := execCLI(t, "", "balanco", "set/2025"); code != ExitUsage {
		t.Errorf("período inválido: código %d, quero %d", code, ExitUsage)
	}
}
