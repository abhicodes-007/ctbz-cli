package cli

import (
	"strings"
	"testing"
)

const guiasEmAtrasoJSON = `{"emAtraso":[{"id":9,"nome":{"label":"DAS"},"tipo":"GUIA","identificadorImposto":"DAS",
"vencimentoOriginal":"2026-08-20","valor":{"label":150.5},"competencia":"Jul de 2026","status":{"label":"VENCIDA"}},
{"id":10,"nome":{"label":"Parcela 2/10"},"tipo":"PARCELA","identificadorImposto":"PARC","vencimentoOriginal":"2026-08-31",
"valor":{"label":null},"competencia":"Jul de 2026","status":{"label":"CALCULANDO"}}],"esteMes":[],"proximoMes":[]}`

func TestImpostos(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/impostos/v5/impostos-a-pagar/guias": fixture(t, "guias_a_pagar")}))

	out, stderr, code := execCLI(t, "", "impostos", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "grupo,id,imposto,competencia,vencimento,valor,situacao,tipo\n" +
		"este_mes,1000000000000001,DARF Unificado,07/2026,2026-10-06,1000.00,RECALCULADA,guia\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}
	if stderr != "" {
		t.Errorf("CSV não deveria escrever totais no stderr: %q", stderr)
	}

	_, stderr, _ = execCLI(t, "", "impostos")
	if !strings.Contains(stderr, "Este mês: R$ 1.000,00 (1)") || !strings.Contains(stderr, "Em atraso: R$ 0,00 (0)") {
		t.Errorf("totais no stderr: %q", stderr)
	}
	if _, _, code := execCLI(t, "", "impostos", "--fail-on-atraso"); code != ExitOK {
		t.Errorf("sem atraso, --fail-on-atraso deveria sair com 0, veio %d", code)
	}
}

func TestImpostosEmAtraso(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/impostos/v5/impostos-a-pagar/guias": guiasEmAtrasoJSON}))

	out, _, code := execCLI(t, "", "impostos", "--atrasadas", "--fail-on-atraso", "-o", "json")
	if code != ExitAttention {
		t.Errorf("código %d, quero %d", code, ExitAttention)
	}
	for _, want := range []string{`"grupo": "em_atraso"`, `"valor": 150.50`, `"valor": null`, `"tipo": "parcela"`} {
		if !strings.Contains(out, want) {
			t.Errorf("JSON sem %s:\n%s", want, out)
		}
	}
}

func TestCompetenciaDeTexto(t *testing.T) {
	for in, want := range map[string]any{
		"Jul de 2026":    "07/2026",
		"Julho / 2026":   "07/2026",
		"set./26":        "09/2026",
		"Dezembro/2025":  "12/2025",
		"texto qualquer": "texto qualquer",
		"":               nil,
	} {
		if got := competenciaDeTexto(in); got != want {
			t.Errorf("competenciaDeTexto(%q) = %v, quero %v", in, got, want)
		}
	}
}

func TestImpostosGuia(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/impostos/v5/impostos-a-pagar/guia/1000000000000001": fixture(t, "guia_detalhe"),
	}))
	out, stderr, code := execCLI(t, "", "impostos", "guia", "1000000000000001", "-o", "json")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	for _, want := range []string{`"competencia": "07/2026"`, `"vencimento": "2026-10-06"`, `"valor_total": 1234.56`,
		`"valor_original": 1000.00`, `"valor_estimado": null`, `"situacao": [`, `"tipo": "guia"`} {
		if !strings.Contains(out, want) {
			t.Errorf("JSON sem %s:\n%s", want, out)
		}
	}
	if _, _, code := execCLI(t, "", "impostos", "guia", "abc"); code != ExitUsage {
		t.Errorf("ID inválido: código %d", code)
	}
}

func TestImpostosCalculoETabelaIRRF(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/impostos/como-imposto-foi-calculado/init":        fixture(t, "calculo_imposto"),
		"/api/plataforma/impostos/como-imposto-foi-calculado/tabela-irrf": fixture(t, "tabela_irrf"),
	}))
	out, stderr, code := execCLI(t, "", "impostos", "calculo", "-o", "json")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	for _, want := range []string{`"mes": "Setembro"`, `"faturamento": 1000.00`, `"das_total": null`, `"fator_r": 100`} {
		if !strings.Contains(out, want) {
			t.Errorf("JSON sem %s:\n%s", want, out)
		}
	}
	out, _, _ = execCLI(t, "", "impostos", "tabela-irrf", "-o", "csv")
	if !strings.HasPrefix(out, "base_calculo,aliquota,deducao\n") || !strings.Contains(out, `"7,5%","R$ 182,16"`) {
		t.Errorf("CSV:\n%s", out)
	}
}
