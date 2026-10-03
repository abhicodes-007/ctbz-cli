package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

const historicoInitJSON = `{"status":"ATIVO","competencia":"10/2026","dataProximoPagamento":"2026-11-10",
	"habilitado":true,"ativado":true,"adyenClientKey":"SEGREDO-CHAVE","cartao":{"final":"4242"}}`

const historicoPagamentosJSON = `{"pagamentos":[
	{"dataPagamento":"2026-09-10","competencia":"09/2026","valor":189.9,"status":"PAGO","cartao":"4242","adyenClientKey":"SEGREDO-CHAVE"},
	{"dataPagamento":1788998400000,"competencia":"08/2026","valor":"1.189,90","status":"PAGO"}]}`

func withHistorico(t *testing.T, pagamentos string) {
	t.Helper()
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/payments/recorrencia/init":      historicoInitJSON,
		"/api/plataforma/payments/recorrencia/historico": pagamentos,
	}))
}

func TestMensalidadeHistoricoFormats(t *testing.T) {
	withHistorico(t, historicoPagamentosJSON)

	out, stderr, code := execCLI(t, "", "mensalidade", "historico", "-o", "json")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	var v struct {
		Ativo      bool   `json:"debito_automatico_ativo"`
		Proximo    string `json:"proximo_pagamento"`
		Pagamentos []struct {
			Data   string  `json:"data"`
			Valor  float64 `json:"valor"`
			Status string  `json:"status"`
		} `json:"pagamentos"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("stdout não é JSON: %v\n%s", err, out)
	}
	if !v.Ativo || v.Proximo != "2026-11-10" || len(v.Pagamentos) != 2 {
		t.Fatalf("JSON inesperado: %s", out)
	}
	if v.Pagamentos[0].Valor != 189.9 || v.Pagamentos[1].Valor != 1189.9 || v.Pagamentos[1].Data == "" {
		t.Errorf("pagamentos inesperados: %s", out)
	}
	if stderr != "" {
		t.Errorf("stderr deveria estar vazio, veio %q", stderr)
	}
	for _, segredo := range []string{"SEGREDO-CHAVE", "adyen", "4242", "cartao"} {
		if strings.Contains(out, segredo) {
			t.Errorf("saída JSON vazou %q:\n%s", segredo, out)
		}
	}

	out, stderr, code = execCLI(t, "", "mensalidade", "historico")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	for _, want := range []string{"Próxima cobrança:", "10/11/2026", "Pagamentos:", "R$ 189,90", "R$ 1.189,90", "PAGO"} {
		if !strings.Contains(out, want) {
			t.Errorf("tabela sem %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "SEGREDO-CHAVE") {
		t.Errorf("tabela vazou segredo:\n%s", out)
	}

	out, stderr, code = execCLI(t, "", "mensalidade", "historico", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	if lines := strings.Split(strings.TrimSpace(out), "\n"); len(lines) != 2 || !strings.HasPrefix(lines[0], "debito_automatico_ativo,") {
		t.Errorf("CSV inesperado:\n%s", out)
	}
}

func TestMensalidadeHistoricoSemPagamentos(t *testing.T) {
	withHistorico(t, `{"pagamentos":[]}`)
	out, stderr, code := execCLI(t, "", "mensalidade", "historico", "-o", "json")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	if !strings.Contains(out, `"pagamentos": []`) {
		t.Errorf("esperava lista vazia:\n%s", out)
	}
}

func TestMensalidadeHistoricoFormatoNaoReconhecido(t *testing.T) {
	withHistorico(t, `{"pagamentos":[{"xyz":1}]}`)
	_, stderr, code := execCLI(t, "", "mensalidade", "historico", "-o", "json")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	if !strings.Contains(stderr, "formato de pagamentos não reconhecido") {
		t.Errorf("esperava aviso no stderr, veio %q", stderr)
	}
}

func TestMensalidadeHistoricoRejeitaArgumento(t *testing.T) {
	t.Setenv("CTBZ_HOME", t.TempDir())
	_, _, code := execCLI(t, "", "mensalidade", "historico", "extra")
	if code != ExitUsage {
		t.Errorf("código %d, quero %d", code, ExitUsage)
	}
}
