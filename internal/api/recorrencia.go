package api

import (
	"context"
	"encoding/json"
)

const (
	PathRecorrenciaInit      = "payments/recorrencia/init"
	PathRecorrenciaHistorico = "payments/recorrencia/historico"
)

// Recorrencia é a situação do pagamento recorrente de impostos (cartão ou conta PJ).
// A resposta real também traz adyenClientKey, que a CLI nunca lê nem mostra.
type Recorrencia struct {
	Habilitado           bool              `json:"habilitado"`
	Ativado              bool              `json:"ativado"`
	Aceite               bool              `json:"aceite"`
	Status               string            `json:"status"`
	Competencia          string            `json:"competencia"` // MM-AAAA
	DataProximoPagamento string            `json:"dataProximoPagamento"`
	DataProximaTentativa string            `json:"dataProximaTentativa"`
	CartoesSalvos        []json.RawMessage `json:"cartoesSalvos"` // só a quantidade é usada
	PagamentosAgendados  []json.RawMessage `json:"pagamentosAgendados"`
	PagamentosConcluido  []json.RawMessage `json:"pagamentosConcluido"`
	PagamentosRecusado   []json.RawMessage `json:"pagamentosRecusado"`
}

// BuscarRecorrencia lê payments/recorrencia/init.
func BuscarRecorrencia(ctx context.Context, g Getter) (*Recorrencia, error) {
	return get[Recorrencia](ctx, g, PathRecorrenciaInit)
}

// PagamentoRecorrente é um mês do histórico de pagamento recorrente. Os campos são os lidos
// pelo front (tela de pagamento recorrente); a conta verificada não tinha histórico.
type PagamentoRecorrente struct {
	Mes   any `json:"mes" contract:"optional"` // número ou texto
	Ano   any `json:"ano" contract:"optional"`
	Guias []struct {
		Nome   string   `json:"nome" contract:"optional"`
		Status string   `json:"status" contract:"optional"`
		Valor  *float64 `json:"valor" contract:"optional"`
	} `json:"guias" contract:"optional"`
	CustoOperacao *float64 `json:"custoOperacao" contract:"optional"`
}

// HistoricoRecorrencia é a resposta de payments/recorrencia/historico (sem paginação).
type HistoricoRecorrencia struct {
	Pagamentos []PagamentoRecorrente `json:"pagamentos"`
}

// BuscarHistoricoRecorrencia lê payments/recorrencia/historico.
func BuscarHistoricoRecorrencia(ctx context.Context, g Getter) (*HistoricoRecorrencia, error) {
	return get[HistoricoRecorrencia](ctx, g, PathRecorrenciaHistorico)
}
