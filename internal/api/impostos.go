package api

import "context"

const PathGuiasAPagar = "impostos/v5/impostos-a-pagar/guias"

// Rotulo é o formato de exibição usado pelas rotas v5: {"label": …, "descricao": …}.
type Rotulo struct {
	Label string `json:"label"`
}

// ValorRotulo é um valor monetário no formato de exibição v5; nulo enquanto calcula.
type ValorRotulo struct {
	Label *float64 `json:"label"`
}

// GuiaResumo é uma guia (ou parcela) na lista de impostos a pagar.
type GuiaResumo struct {
	ID                   int64       `json:"id"`
	Nome                 Rotulo      `json:"nome"`
	Tipo                 string      `json:"tipo"` // GUIA ou PARCELA
	IdentificadorImposto string      `json:"identificadorImposto"`
	VencimentoOriginal   string      `json:"vencimentoOriginal"` // AAAA-MM-DD
	Valor                ValorRotulo `json:"valor"`
	Competencia          string      `json:"competencia"` // ex.: "Jul de 2026"
	Status               Rotulo      `json:"status"`
}

// GuiasAPagar é a resposta de impostos/v5/impostos-a-pagar/guias.
type GuiasAPagar struct {
	EmAtraso   []GuiaResumo `json:"emAtraso"`
	EsteMes    []GuiaResumo `json:"esteMes"`
	ProximoMes []GuiaResumo `json:"proximoMes"`
}

// BuscarGuiasAPagar lê as guias em atraso, do mês e do próximo mês.
func BuscarGuiasAPagar(ctx context.Context, g Getter) (*GuiasAPagar, error) {
	return get[GuiasAPagar](ctx, g, PathGuiasAPagar)
}
