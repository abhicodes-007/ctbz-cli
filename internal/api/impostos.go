package api

import (
	"context"
	"fmt"
)

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

const (
	PathCalculoImposto = "impostos/como-imposto-foi-calculado/init"
	PathTabelaIRRF     = "impostos/como-imposto-foi-calculado/tabela-irrf"
)

// PathGuia é o detalhe de uma guia (o ID vem de GuiasAPagar).
func PathGuia(id int64) string {
	return fmt.Sprintf("impostos/v5/impostos-a-pagar/guia/%d", id)
}

// Montante é um valor no formato {"valor": …, "descricao": …}; nulo quando não se aplica.
type Montante struct {
	Valor *float64 `json:"valor"`
}

// GuiaDetalhe é a resposta de impostos/v5/impostos-a-pagar/guia/{id}.
type GuiaDetalhe struct {
	ID                   int64  `json:"id"`
	Tipo                 string `json:"tipo"`
	IdentificadorImposto string `json:"identificadorImposto"`
	Competencia          string `json:"competencia"`
	Vencimento           string `json:"vencimento"` // dd/mm/aaaa
	Oraculo              struct {
		Nome       string `json:"nome"`
		Descricao  string `json:"descricao"`
		Frequencia string `json:"frequencia"`
		Impacto    string `json:"impacto"`
	} `json:"oraculo"`
	Status           []string  `json:"status"`
	ValorTotal       *Montante `json:"valorTotal"`
	ValorOriginal    *Montante `json:"valorOriginal"`
	ValorJurosEMulta *Montante `json:"valorJurosEMulta"`
	ValorEstimado    *Montante `json:"valorEstimado"`
	ValorEmAtraso    *Montante `json:"valorEmAtraso"`
	AcoesBotoes      []string  `json:"acoesBotoes"`
}

// BuscarGuia lê o detalhe de uma guia.
func BuscarGuia(ctx context.Context, g Getter, id int64) (*GuiaDetalhe, error) {
	return get[GuiaDetalhe](ctx, g, PathGuia(id))
}

// CalculoImposto é a resposta de "como meu imposto foi calculado": a memória de cálculo
// do mês mais recente (DAS do Simples e DARF de INSS/IRRF sobre o pró-labore).
type CalculoImposto struct {
	NomeMesCompetencia string   `json:"nomeMesCompetencia"`
	FaturamentoTotal   *float64 `json:"faturamentoTotal"`
	DasSimples         struct {
		ImpostoBruto    *float64 `json:"impostoBruto"`
		DeducaoRetencao *float64 `json:"deducaoRetencao"`
		ImpostoTotal    *float64 `json:"impostoTotal"`
		Inconsistente   bool     `json:"inconsistente"`
	} `json:"dasSimples"`
	Darf struct {
		INSS struct {
			Prolabore    *float64 `json:"prolabore"`
			Aliquota     *float64 `json:"aliquota"`
			Teto         *float64 `json:"teto"`
			TotalImposto *float64 `json:"totalImposto"`
		} `json:"inss"`
		IRRF struct {
			Prolabore       *float64 `json:"prolabore"`
			BaseCalculoIRRF *float64 `json:"baseCalculoIrrf"`
			Aliquota        *float64 `json:"aliquota"`
			DeducaoIRRF     *float64 `json:"deducaoIrrf"`
			TotalImposto    *float64 `json:"totalImposto"`
		} `json:"irrf"`
		Total *float64 `json:"total"`
	} `json:"darf"`
	ValorFaturamentoUltimos12Meses *float64 `json:"valorFaturamentoUltimos12Meses"`
	ValorProlaboreUltimos12Meses   *float64 `json:"valorProlaboreUltimos12Meses"`
	PercentualFatorR               *float64 `json:"percentualFatorR"`
	HistoricoFaturamento           []struct {
		Mes              string   `json:"mes"` // ex.: "set./26"
		ValorFaturamento *float64 `json:"valorFaturamento"`
		ValorProlabore   *float64 `json:"valorProlabore"`
	} `json:"historicoFaturamento"`
}

// BuscarCalculoImposto lê a memória de cálculo do mês mais recente.
func BuscarCalculoImposto(ctx context.Context, g Getter) (*CalculoImposto, error) {
	return get[CalculoImposto](ctx, g, PathCalculoImposto)
}

// FaixaIRRF é uma faixa da tabela progressiva do IRRF, em texto como a API devolve.
type FaixaIRRF struct {
	BaseCalculo string `json:"baseCalculo"`
	Aliquota    string `json:"aliquota"`
	Deducao     string `json:"deducao"`
}

// BuscarTabelaIRRF lê a tabela progressiva do IRRF vigente.
func BuscarTabelaIRRF(ctx context.Context, g Getter) ([]FaixaIRRF, error) {
	t, err := get[[]FaixaIRRF](ctx, g, PathTabelaIRRF)
	if err != nil {
		return nil, err
	}
	return *t, nil
}
