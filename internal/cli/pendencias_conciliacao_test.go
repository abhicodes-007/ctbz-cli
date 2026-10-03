package cli

import (
	"strings"
	"testing"
)

const pathConciliacao = "/api/plataforma/conciliacao-fiscal/v2/pendencias?pagina=%d&periodoFinal=2026-10-03&periodoInicial=2025-10-03" +
	"&status=PENDENTE&tipoPendencia=RECEITA_SEM_RECEBIMENTO&totalPagina=10"

func TestPendenciasConciliacao(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/conciliacao-fiscal/v2/init": fixture(t, "conciliacao_resumo")}))

	out, stderr, code := execCLI(t, "", "pendencias", "conciliacao", "--fail-on-pendencias", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("sem pendências deveria sair com 0, veio %d: %s", code, stderr)
	}
	want := "competencia_referencia,notas_fiscais_pendentes,recebimentos_pendentes,conciliacoes_automaticas\n09/2026,0,0,0\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	if _, _, code := execCLI(t, "", "pendencias", "conciliacao", "--listar", "boletos"); code != ExitUsage {
		t.Errorf("--listar inválido: código %d, quero %d", code, ExitUsage)
	}
}

func TestPendenciasConciliacaoListar(t *testing.T) {
	fixNow(t, "2026-10-03")
	pagina := func(n int) string { return strings.Replace(pathConciliacao, "%d", string(rune('0'+n)), 1) }
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/conciliacao-fiscal/v2/init": `{"qtdNotasFiscaisPendentes":2,"qtdRecebimentosPendentes":0,` +
			`"qtdConciliacoesAutomaticasMesAnterior":1,"competenciaMesAnterior":"2026-09"}`,
		pagina(1): `{"pagina":[{"id":1,"numeroNota":"10","valor":100.5}],"totalRegistros":2,"totalPaginas":2}`,
		pagina(2): `{"pagina":[{"id":2,"numeroNota":"11","valor":50}],"totalRegistros":2,"totalPaginas":2}`,
	}))

	out, stderr, code := execCLI(t, "", "pendencias", "conciliacao", "--listar", "notas", "--fail-on-pendencias", "-o", "csv")
	if code != ExitAttention {
		t.Errorf("com pendências, código %d, quero %d: %s", code, ExitAttention, stderr)
	}
	if want := "id,numeroNota,valor\n1,10,100.5\n2,11,50\n"; out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	if _, _, code := execCLI(t, "", "pendencias", "conciliacao", "--fail-on-pendencias"); code != ExitAttention {
		t.Errorf("resumo com pendências: código %d, quero %d", code, ExitAttention)
	}
}
