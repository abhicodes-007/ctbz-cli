package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newImpostosGuiaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "guia ID",
		Short: "Mostra o detalhe de uma guia de imposto",
		Long: `Mostra o detalhe de uma guia: imposto, competência, vencimento, valor total, valor
original, juros e multa, situação e ações disponíveis no painel. O ID vem de "ctbz impostos".

Para ver como o imposto do mês foi calculado, use "ctbz impostos calculo".`,
		Example: `  ctbz impostos guia 1000000000000001`,
		Args:    exactArgs(1, "o ID da guia"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			g, err := api.BuscarGuia(cmd.Context(), sessionGetter{s}, id)
			if err != nil {
				return err
			}
			return output.Write(s.out, f, guiaRecord(g))
		},
	}
}

func guiaRecord(g *api.GuiaDetalhe) *output.Record {
	venc, _ := output.ParseDate(g.Vencimento)
	rec := &output.Record{}
	rec.Add("id", "ID", g.ID).
		Add("imposto", "Imposto", g.Oraculo.Nome).
		Add("identificador", "Identificador", g.IdentificadorImposto).
		Add("tipo", "Tipo", strings.ToLower(g.Tipo)).
		Add("competencia", "Competência", competenciaDeTexto(g.Competencia)).
		Add("vencimento", "Vencimento", venc).
		Add("valor_total", "Valor total", montante(g.ValorTotal)).
		Add("valor_original", "Valor original", montante(g.ValorOriginal)).
		Add("juros_e_multa", "Juros e multa", montante(g.ValorJurosEMulta)).
		Add("valor_estimado", "Valor estimado", montante(g.ValorEstimado)).
		Add("valor_em_atraso", "Valor em atraso", montante(g.ValorEmAtraso)).
		Add("situacao", "Situação", nonNil(g.Status)).
		Add("acoes", "Ações no painel", nonNil(g.AcoesBotoes)).
		Add("descricao", "Descrição", nilIfEmpty(g.Oraculo.Descricao)).
		Add("frequencia", "Frequência", nilIfEmpty(g.Oraculo.Frequencia)).
		Add("impacto_do_atraso", "Impacto do atraso", nilIfEmpty(g.Oraculo.Impacto))
	return rec
}

func montante(m *api.Montante) any {
	if m == nil {
		return nil
	}
	return moneyOrNil(m.Valor)
}

func parseID(s string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil || id <= 0 {
		return 0, usageError{fmt.Errorf("ID inválido: %q", s)}
	}
	return id, nil
}

func newImpostosCalculoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "calculo",
		Short: "Mostra como o imposto do mês foi calculado",
		Long: `Mostra a memória de cálculo do mês mais recente, como a tela "Como meu imposto foi
calculado": faturamento, DAS do Simples Nacional, DARF de INSS e IRRF sobre o pró-labore,
faturamento e pró-labore dos últimos 12 meses e o Fator R.

A competência é escolhida pela Contabilizei (a API não recebe mês).`,
		Example: `  ctbz impostos calculo
  ctbz impostos calculo -o json | jq .fator_r`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			c, err := api.BuscarCalculoImposto(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, calculoRecord(c))
		},
	}
}

func calculoRecord(c *api.CalculoImposto) *output.Record {
	das, inss, irrf := c.DasSimples, c.Darf.INSS, c.Darf.IRRF
	rec := &output.Record{}
	rec.Add("mes", "Mês", c.NomeMesCompetencia).
		Add("faturamento", "Faturamento", moneyOrNil(c.FaturamentoTotal)).
		Add("das_imposto_bruto", "DAS: imposto bruto", moneyOrNil(das.ImpostoBruto)).
		Add("das_deducao_retencao", "DAS: dedução de retenções", moneyOrNil(das.DeducaoRetencao)).
		Add("das_total", "DAS: total", moneyOrNil(das.ImpostoTotal)).
		Add("das_inconsistente", "DAS inconsistente", das.Inconsistente).
		Add("inss_prolabore", "INSS: pró-labore", moneyOrNil(inss.Prolabore)).
		Add("inss_aliquota", "INSS: alíquota (%)", floatOrNil(inss.Aliquota)).
		Add("inss_total", "INSS: total", moneyOrNil(inss.TotalImposto)).
		Add("irrf_base_calculo", "IRRF: base de cálculo", moneyOrNil(irrf.BaseCalculoIRRF)).
		Add("irrf_aliquota", "IRRF: alíquota (%)", floatOrNil(irrf.Aliquota)).
		Add("irrf_deducao", "IRRF: dedução", moneyOrNil(irrf.DeducaoIRRF)).
		Add("irrf_total", "IRRF: total", moneyOrNil(irrf.TotalImposto)).
		Add("darf_total", "DARF: total", moneyOrNil(c.Darf.Total)).
		Add("faturamento_12_meses", "Faturamento (12 meses)", moneyOrNil(c.ValorFaturamentoUltimos12Meses)).
		Add("prolabore_12_meses", "Pró-labore (12 meses)", moneyOrNil(c.ValorProlaboreUltimos12Meses)).
		Add("fator_r", "Fator R (%)", floatOrNil(c.PercentualFatorR))
	return rec
}

func floatOrNil(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}

func newImpostosTabelaIRRFCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "tabela-irrf",
		Short:   "Mostra a tabela progressiva do IRRF usada no cálculo do pró-labore",
		Example: `  ctbz impostos tabela-irrf`,
		Args:    exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			faixas, err := api.BuscarTabelaIRRF(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			l := &output.List{Columns: []output.Column{
				{Key: "base_calculo", Header: "Base de cálculo"},
				{Key: "aliquota", Header: "Alíquota"},
				{Key: "deducao", Header: "Dedução"},
			}}
			for _, fx := range faixas {
				l.Append(fx.BaseCalculo, nilIfEmpty(fx.Aliquota), nilIfEmpty(fx.Deducao))
			}
			return output.Write(s.out, f, l)
		},
	}
}
