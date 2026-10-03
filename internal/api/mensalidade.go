package api

import "context"

const (
	PathFatura      = "dashboard/fatura"
	PathMensalidade = "dashboard/v1/mensalidade"
)

// Fatura é a fatura atual da Contabilizei (card "Mensalidade" do painel).
type Fatura struct {
	Total                 *float64 `json:"total"`
	Competencia           string   `json:"competencia"` // ex.: "Outubro, 2026"
	Status                string   `json:"status"`      // ex.: "Em aberto"
	Vencimento            any      `json:"vencimento"`  // formato não verificado (vinha null)
	MensagemAdicional     string   `json:"mensagemAdicional"`
	PossuiCartaoPrincipal bool     `json:"possuiCartaoPrincipal"`
}

// BuscarFatura lê dashboard/fatura.
func BuscarFatura(ctx context.Context, g Getter) (*Fatura, error) {
	return get[Fatura](ctx, g, PathFatura)
}

// Mensalidade é a resposta de dashboard/v1/mensalidade. Na conta verificada só
// competenciaAnteriorAtrasada vinha preenchido; os demais campos têm formato não verificado.
type Mensalidade struct {
	Status                      any      `json:"status"`
	Valor                       *float64 `json:"valor"`
	DataVencimento              any      `json:"dataVencimento"`
	CompetenciaAnteriorAtrasada bool     `json:"competenciaAnteriorAtrasada"`
}

// BuscarMensalidade lê dashboard/v1/mensalidade.
func BuscarMensalidade(ctx context.Context, g Getter) (*Mensalidade, error) {
	return get[Mensalidade](ctx, g, PathMensalidade)
}
