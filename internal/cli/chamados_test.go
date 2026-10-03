package cli

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChamadosEmAndamento(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/atendimento/chamados?em-andamento=true": `[
{"id":123,"assunto":"Dúvida sobre o DAS","canal":"chat","status":"open","atualizado":"02/10/2026","previsaoRetorno":"em até 2 dias úteis"}]`,
	}))
	out, stderr, code := execCLI(t, "", "chamados", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "id,assunto,status,canal,criado,atualizado,previsao_retorno,link\n" +
		"123,Dúvida sobre o DAS,open,chat,,2026-10-02,em até 2 dias úteis,https://suporte.contabilizei.com.br/hc/pt-br/requests/123\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
}

func TestChamadosEmAndamentoVazio(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/atendimento/chamados?em-andamento=true": fixture(t, "chamados_em_andamento"),
	}))
	if out, _, code := execCLI(t, "", "chamados", "-o", "json"); code != ExitOK || out != "[]\n" {
		t.Errorf("sem chamados (código %d): %q", code, out)
	}
}

func TestChamadosFinalizadosComFonteAlternativa(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RequestURI() {
		case "/api/plataforma/atendimento/chamados?finalizados=true":
			http.Error(w, "Memcache put: Item may not be more than 1048503 bytes in length", http.StatusBadRequest)
		case "/api/plataforma/dadosempresa/get":
			w.Write([]byte(`{"empresaAtual":{},"empresas":[],"chamados":[
{"id":"1","status":"open","subject":"Aberto","created_at":"01/10/2026","url":"x"},
{"id":"2","status":"solved","subject":"Resolvido","created_at":"22/07/2026","url":"x"},
{"id":"3","status":"closed","subject":"Fechado","created_at":"03/05/2023","url":"x"}]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	withSession(t, srv)

	out, stderr, code := execCLI(t, "", "chamados", "--finalizados", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	if !strings.Contains(stderr, "HTTP 400") || !strings.Contains(stderr, "100 mais recentes") {
		t.Errorf("aviso no stderr: %q", stderr)
	}
	want := "id,assunto,status,canal,criado,atualizado,previsao_retorno,link\n" +
		"2,Resolvido,solved,,2026-07-22,,,https://suporte.contabilizei.com.br/hc/pt-br/requests/2\n" +
		"3,Fechado,closed,,2023-05-03,,,https://suporte.contabilizei.com.br/hc/pt-br/requests/3\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
}

func TestChamadosFinalizadosFixture(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/dadosempresa/get": fixture(t, "dadosempresa")}))
	out, _, code := execCLI(t, "", "chamados", "--finalizados", "-o", "json")
	if code != ExitOK || strings.Count(out, `"status": "closed"`) != 3 {
		t.Errorf("fallback com a fixture (código %d):\n%s", code, out)
	}
}
