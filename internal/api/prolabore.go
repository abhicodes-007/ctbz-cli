package api

import (
	"context"
	"strconv"
)

const (
	PathProlaboreCentral   = "prolabore/central/init"
	PathProlaboreDashboard = "dashboard/prolabore"
)

// PathProlaboreHistorico é o histórico mensal de pró-labore de um sócio (id de
// ProlaboreCentral.Socios). A API não tem filtro por ano nem paginação.
func PathProlaboreHistorico(idSocio int64) string {
	return "prolabore/central/historico/" + strconv.FormatInt(idSocio, 10)
}

// SocioProlabore é um sócio na central de pró-labore.
type SocioProlabore struct {
	ID                    int64    `json:"id"`
	Nome                  string   `json:"nome"`
	CPF                   string   `json:"cpf"`
	PossuiProlabore       bool     `json:"possuiProlabore"`
	ResponsavelReceita    bool     `json:"responsavelReceita"`
	ValorProlabore        *float64 `json:"valorProlabore"`
	DataUltimaAtualizacao string   `json:"dataUltimaAtualizacao"` // dd/mm/aaaa
	InsightGestao         string   `json:"insightGestao"`         // ex.: INTELIGENTE
	NomesDependentes      []string `json:"nomesDependentes"`
}

// ProlaboreCentral é a resposta de prolabore/central/init.
type ProlaboreCentral struct {
	TipoGerenciamento     string           `json:"tipoGerenciamento"` // ex.: INTELIGENTE
	ElegivelNoMotor       bool             `json:"elegivelNoMotor"`
	Socios                []SocioProlabore `json:"socios"`
	TotalProLabore        *float64         `json:"totalProLabore"`
	ProlaboreIndisponivel bool             `json:"prolaboreIndisponivel"`
	ZerarProlabore        bool             `json:"zerarProlabore"`
}

// BuscarProlaboreCentral lê a central de pró-labore.
func BuscarProlaboreCentral(ctx context.Context, g Getter) (*ProlaboreCentral, error) {
	return get[ProlaboreCentral](ctx, g, PathProlaboreCentral)
}

// ProlaboreDashboard é o card de pró-labore do painel. Na conta verificada valor e
// competências vinham null; as competências ficam como any até haver exemplo.
type ProlaboreDashboard struct {
	ValorProlabore      *float64 `json:"valorProlabore"`
	CompetenciaAtual    any      `json:"competenciaAtual"`
	CompetenciaAnterior any      `json:"competenciaAnterior"`
	PresenteMotor       bool     `json:"presenteMotor"`
	Calculando          bool     `json:"calculando"`
}

// BuscarProlaboreDashboard lê o card de pró-labore.
func BuscarProlaboreDashboard(ctx context.Context, g Getter) (*ProlaboreDashboard, error) {
	return get[ProlaboreDashboard](ctx, g, PathProlaboreDashboard)
}

// ProlaboreMes é um mês do histórico, com valores já formatados pela API ("R$ 100,00").
type ProlaboreMes struct {
	Competencia string `json:"competencia"` // ex.: "Julho/2026"
	Nome        string `json:"nome"`
	Prolabore   string `json:"prolabore"`
	Descontos   string `json:"descontos"`
}

// BuscarProlaboreHistorico lê o histórico de um sócio.
func BuscarProlaboreHistorico(ctx context.Context, g Getter, idSocio int64) ([]ProlaboreMes, error) {
	h, err := get[[]ProlaboreMes](ctx, g, PathProlaboreHistorico(idSocio))
	if err != nil {
		return nil, err
	}
	return *h, nil
}

const PathProlaboreParametros = "prolabore/init"

// ProlaboreParametros são os valores usados no cálculo do pró-labore, já formatados pela
// API ("R$ 1.621,00"), exceto a porcentagem do INSS.
type ProlaboreParametros struct {
	SalarioMinimo                          string  `json:"salarioMinimo"`
	ValorMaximoContribuicaoInss            string  `json:"valorMaximoContribuicaoInss"`
	ValorMaximoProlabore                   string  `json:"valorMaximoProlabore"`
	PorcentagemInss                        float64 `json:"porcentagemInss"`
	ValorMinimoProlaboreParaIncidenciaIrrf string  `json:"valorMinimoProlaboreParaIncidenciaIrrf"`
}

// BuscarProlaboreParametros lê prolabore/init.
func BuscarProlaboreParametros(ctx context.Context, g Getter) (*ProlaboreParametros, error) {
	return get[ProlaboreParametros](ctx, g, PathProlaboreParametros)
}
