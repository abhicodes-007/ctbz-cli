package api

import "context"

const (
	PathContratoServico = "contrato/buscarContratoServico/"
	// PathPropostaPlano responde HTML puro (não JSON) e fica na base do legado.
	PathPropostaPlano = "/api/legado/contrato/buscarContratoPlanoPagamento/"
)

// ContratoServico é o contrato de prestação de serviços aceito pela empresa.
type ContratoServico struct {
	HTML                string `json:"html"`
	Versao              string `json:"versao"`
	IDHistoricoContrato int64  `json:"idHistoricoContrato"`
}

// BuscarContratoServico lê o contrato de serviço (HTML dentro do JSON).
func BuscarContratoServico(ctx context.Context, g Getter) (*ContratoServico, error) {
	return get[ContratoServico](ctx, g, PathContratoServico)
}

// BuscarPropostaPlano lê a proposta do plano contratado, em HTML.
func BuscarPropostaPlano(ctx context.Context, g TextGetter) (string, error) {
	return g.GetText(ctx, PathPropostaPlano)
}
