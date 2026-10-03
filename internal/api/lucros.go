package api

import (
	"context"
	"fmt"
)

const PathDistribuicaoLucros = "informerendimento/recuperardadosdistribuicaocliente"

// PathRestricoesInforme são as restrições do informe de rendimentos de um ano (pendências
// documentais, débitos federais, reabertura do balanço).
func PathRestricoesInforme(ano int) string {
	return fmt.Sprintf("informerendimento/v2/%d/restricoes", ano)
}

// DistribuicaoLucros é a distribuição de lucros do exercício aberto. A API não recebe ano:
// o exercício vem em Ano (0 quando ainda não há exercício para distribuir).
type DistribuicaoLucros struct {
	Ano                               int      `json:"ano"`
	Saldo                             *float64 `json:"saldo"`
	TotalDistribuido                  *float64 `json:"totalDistribuido"`
	TotalAdiantamentos                *float64 `json:"totalAdiantamentos"`
	LimitePermitidoDistribuicaoLucros *float64 `json:"limitePermitidoDistribuicaoLucros"`
	ExercicioFechado                  bool     `json:"exercicioFechado"`
	PodeAlterar                       bool     `json:"podeAlterar"`
	MotivoNaoPodeAlterar              any      `json:"motivoNaoPodeAlterar"`
	DataLimite                        any      `json:"dataLimite"`
	// Itens lidos pelo front ({socio, valor}); a conta verificada não tinha distribuição.
	LucrosSocios []struct {
		Socio string   `json:"socio" contract:"optional"`
		Valor *float64 `json:"valor" contract:"optional"`
	} `json:"lucrosSocios"`
}

// BuscarDistribuicaoLucros lê a distribuição de lucros.
func BuscarDistribuicaoLucros(ctx context.Context, g Getter) (*DistribuicaoLucros, error) {
	return get[DistribuicaoLucros](ctx, g, PathDistribuicaoLucros)
}

// RestricoesInforme indica o que impede o informe de rendimentos de um ano.
type RestricoesInforme struct {
	Restricoes struct {
		PendenciaDocumental struct {
			PossuiPendencia    bool   `json:"possuiPendencia"`
			FluxoRegularizacao string `json:"fluxoRegularizacao"`
		} `json:"pendenciaDocumental"`
		DebitosFederais struct {
			PossuiPendencia           bool `json:"possuiPendencia"`
			DivergenciaContabilFiscal bool `json:"divergenciaContabilFiscal"`
		} `json:"debitosFederais"`
	} `json:"restricoes"`
	ProcessoReabertura struct {
		Status string `json:"status"` // ex.: NENHUM, EM_ANDAMENTO
	} `json:"processoReabertura"`
}

// BuscarRestricoesInforme lê as restrições do informe de um ano.
func BuscarRestricoesInforme(ctx context.Context, g Getter, ano int) (*RestricoesInforme, error) {
	return get[RestricoesInforme](ctx, g, PathRestricoesInforme(ano))
}
