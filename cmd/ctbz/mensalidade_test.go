package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/edusouza/ctbz-cli/internal/ctbz"
)

const fixtureInit = `{"status":"ATIVO","competencia":"10/2026","dataProximoPagamento":"2026-11-10",
	"habilitado":true,"ativado":true,"adyenClientKey":"SEGREDO-CHAVE","cartao":{"final":"4242"}}`

const fixtureHistorico = `{"pagamentos":[
	{"dataPagamento":"2026-09-10","competencia":"09/2026","valor":189.9,"status":"PAGO","cartao":"4242","adyenClientKey":"SEGREDO-CHAVE"},
	{"dataPagamento":1788998400000,"competencia":"08/2026","valor":"1.189,90","status":"PAGO"}]}`

func setupMensalidade(t *testing.T, historico string) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/plataforma/payments/recorrencia/init":
			io.WriteString(w, fixtureInit)
		case "/api/plataforma/payments/recorrencia/historico":
			io.WriteString(w, historico)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	store := &ctbz.Store{Dir: t.TempDir()}
	t.Setenv("CTBZ_HOME", store.Dir)
	t.Setenv("CTBZ_OUTPUT", "")
	if err := store.SaveSession(&ctbz.Session{BaseURL: srv.URL, CreatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
}

func TestMensalidadeHistoricoFormats(t *testing.T) {
	setupMensalidade(t, fixtureHistorico)
	ctx := context.Background()

	out, errOut, err := capture(t, func() error { return cmdMensalidade(ctx, []string{"historico", "-o", "json"}) })
	if err != nil {
		t.Fatal(err)
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
	if errOut != "" {
		t.Errorf("stderr deveria estar vazio, veio %q", errOut)
	}
	for _, segredo := range []string{"SEGREDO-CHAVE", "adyen", "4242", "cartao"} {
		if strings.Contains(out, segredo) {
			t.Errorf("saída JSON vazou %q:\n%s", segredo, out)
		}
	}

	out, _, err = capture(t, func() error { return cmdMensalidade(ctx, nil) })
	if err == nil {
		t.Error("esperava erro sem subcomando")
	}

	out, _, err = capture(t, func() error { return cmdMensalidade(ctx, []string{"historico"}) })
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Próxima cobrança:", "10/11/2026", "Pagamentos:", "R$ 189,90", "R$ 1.189,90", "PAGO"} {
		if !strings.Contains(out, want) {
			t.Errorf("tabela sem %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "SEGREDO-CHAVE") {
		t.Errorf("tabela vazou segredo:\n%s", out)
	}

	out, _, err = capture(t, func() error { return cmdMensalidade(ctx, []string{"historico", "-o", "csv"}) })
	if err != nil {
		t.Fatal(err)
	}
	if lines := strings.Split(strings.TrimSpace(out), "\n"); len(lines) != 2 || !strings.HasPrefix(lines[0], "debito_automatico_ativo,") {
		t.Errorf("CSV inesperado:\n%s", out)
	}
}

func TestMensalidadeHistoricoSemPagamentos(t *testing.T) {
	setupMensalidade(t, `{"pagamentos":[]}`)
	out, _, err := capture(t, func() error {
		return cmdMensalidade(context.Background(), []string{"historico", "-o", "json"})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"pagamentos": []`) {
		t.Errorf("esperava lista vazia:\n%s", out)
	}
}
