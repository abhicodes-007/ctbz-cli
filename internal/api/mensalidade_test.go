package api

import (
	"context"
	"testing"
)

type textoFixo string

func (t textoFixo) GetText(context.Context, string) (string, error) { return string(t), nil }

// A situação da mensalidade é texto puro, fora dos contratos JSON: o formato visto na API
// real é "OK" sem aspas.
func TestBuscarSituacaoMensalidade(t *testing.T) {
	for in, want := range map[string]string{"OK": "OK", " OK\n": "OK", `"OK"`: "OK", "INADIMPLENTE": "INADIMPLENTE"} {
		got, err := BuscarSituacaoMensalidade(context.Background(), textoFixo(in))
		if err != nil || got != want {
			t.Errorf("BuscarSituacaoMensalidade(%q) = %q, %v; quero %q", in, got, err, want)
		}
	}
}
