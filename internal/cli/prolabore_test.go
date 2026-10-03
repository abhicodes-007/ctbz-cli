package cli

import (
	"strings"
	"testing"
)

func TestProlabore(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/prolabore/central/init": fixture(t, "prolabore_central"),
		"/api/plataforma/dashboard/prolabore":    fixture(t, "prolabore_dashboard"),
	}))
	out, stderr, code := execCLI(t, "", "prolabore", "-o", "json")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	for _, want := range []string{`"gerenciamento": "INTELIGENTE"`, `"total": 1000.00`, `"competencia_atual": null`,
		`"indisponivel": true`, `"socios": [`, `"cpf": "00000000000"`, `"atualizado": "2026-09-01"`, `"dependentes": 2`} {
		if !strings.Contains(out, want) {
			t.Errorf("JSON sem %s:\n%s", want, out)
		}
	}
	if out, _, _ := execCLI(t, "", "prolabore"); !strings.Contains(out, "Gerenciamento:") || !strings.Contains(out, "FULANO DE TAL") {
		t.Errorf("tabela:\n%s", out)
	}
}

func TestProlaboreHistorico(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/prolabore/central/init": `{"socios":[{"id":1,"nome":"A"},{"id":2,"nome":"B"}]}`,
		"/api/plataforma/prolabore/central/historico/1": `[{"competencia":"Julho/2026","nome":"A","prolabore":"R$ 1.621,00","descontos":"R$ 178,31"},
{"competencia":"Dezembro/2025","nome":"A","prolabore":"R$ 1.518,00","descontos":"R$ 166,98"}]`,
		"/api/plataforma/prolabore/central/historico/2": `[{"competencia":"Julho/2026","nome":"B","prolabore":"-","descontos":""}]`,
	}))
	out, stderr, code := execCLI(t, "", "prolabore", "historico", "--ano", "2026", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "competencia,socio,prolabore,descontos\n07/2026,A,1621.00,178.31\n07/2026,B,-,\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	out, _, code = execCLI(t, "", "prolabore", "historico", "--socio", "1", "-o", "csv")
	if code != ExitOK || strings.Count(out, "\n") != 3 || !strings.Contains(out, "12/2025,A,1518.00,166.98") {
		t.Errorf("--socio 1 (código %d):\n%s", code, out)
	}
}
