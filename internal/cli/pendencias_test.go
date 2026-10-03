package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/edusouza/ctbz-cli/internal/output"
)

// fixNow fixa "hoje" (now) numa data AAAA-MM-DD durante o teste.
func fixNow(t *testing.T, dia string) {
	t.Helper()
	d, err := time.ParseInLocation("2006-01-02", dia, time.Local)
	if err != nil {
		t.Fatal(err)
	}
	orig := now
	now = func() time.Time { return d.Add(10 * time.Hour) }
	t.Cleanup(func() { now = orig })
}

func TestPendencias(t *testing.T) {
	fixNow(t, "2026-10-03")
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/home/pendencia/pendenciasEmpresa": fixture(t, "pendencias_empresa")}))

	// O detalhe da fixture ("<texto omitido>") parece uma tag HTML e some, como deve.
	out, stderr, code := execCLI(t, "", "pendencias", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "id,tipo,detalhe,criada,prazo,situacao,alerta\n" +
		"1000000000000001,Cadastre o PIS para ativar o pró-labore automático,,2026-07-21,2026-07-21,Pendente,vencida\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}

	out, _, code = execCLI(t, "", "pendencias", "--todas", "--fail-on-vencidas", "-o", "json")
	if code != ExitAttention {
		t.Errorf("com pendência vencida, código %d, quero %d", code, ExitAttention)
	}
	for _, want := range []string{`"situacao": "Finalizada"`, `"alerta": null`, `"alerta": "vencida"`} {
		if !strings.Contains(out, want) {
			t.Errorf("JSON sem %s:\n%s", want, out)
		}
	}
}

func TestPendenciasProximaEDetalheHTML(t *testing.T) {
	fixNow(t, "2026-07-18")
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/home/pendencia/pendenciasEmpresa": `[
{"id":1,"dataCriacao":0,"dataLimite":"21/07/2026","detalhe":"Cadastre o PIS. <br> Veja <a href=\"x\">aqui</a> &amp; pronto",
 "tipoPendencia":{"codigo":43,"titulo":"PIS"},"situacaoPendencia":{"id":"PENDENTE","descricao":"Pendente"}},
{"id":2,"dataCriacao":0,"dataLimite":"","detalhe":"","tipoPendencia":{"codigo":1,"titulo":"Sem prazo"},
 "situacaoPendencia":{"id":"EM_ANALISE","descricao":"Em análise"}}]`}))

	out, stderr, code := execCLI(t, "", "pendencias", "--fail-on-vencidas", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("sem vencidas deveria sair com 0, veio %d: %s", code, stderr)
	}
	want := "id,tipo,detalhe,criada,prazo,situacao,alerta\n" +
		"1,PIS,Cadastre o PIS. Veja aqui & pronto,,2026-07-21,Pendente,próxima\n" +
		"2,Sem prazo,,,,Em análise,\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
}

func TestAlertaDePrazo(t *testing.T) {
	fixNow(t, "2026-10-03")
	dia := func(s string) output.Date { d, _ := output.ParseDate(s); return d }
	for _, c := range []struct {
		prazo  string
		aberto bool
		want   any
	}{
		{"2026-10-02", true, alertaVencida},
		{"2026-10-03", true, alertaProxima},
		{"2026-10-10", true, alertaProxima},
		{"2026-10-11", true, nil},
		{"2026-10-02", false, nil},
		{"", true, nil},
	} {
		if got := alertaDePrazo(dia(c.prazo), c.aberto); got != c.want {
			t.Errorf("alertaDePrazo(%s, %v) = %v, quero %v", c.prazo, c.aberto, got, c.want)
		}
	}
}
