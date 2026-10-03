package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newImpostosHistoricoCmd() *cobra.Command {
	var f api.FiltroHistorico
	cmd := &cobra.Command{
		Use:   "historico",
		Short: "Lista o histórico de guias de impostos",
		Long: `Lista as guias de meses anteriores, com valor, valor pago, vencimento e situação,
lendo todas as páginas. Filtros por ano, mês (1 a 12) e situação (ex.: PAGO, PENDENTE).
Na tabela, o resumo (em dia / guias vencidas) vai para stderr.

Para baixar o PDF de uma guia do histórico, use "ctbz impostos baixar ID".`,
		Example: `  ctbz impostos historico --ano 2026
  ctbz impostos historico --ano 2026 --mes 7 -o csv`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if f.Mes < 0 || f.Mes > 12 {
				return usageError{fmt.Errorf("--mes deve estar entre 1 e 12")}
			}
			format, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			g := sessionGetter{s}
			f.Status = strings.ToUpper(f.Status)
			guias, err := api.BuscarHistoricoGuias(cmd.Context(), g, f)
			if err != nil {
				return err
			}
			if err := output.Write(s.out, format, historicoList(guias)); err != nil {
				return err
			}
			if format == output.FormatTable {
				if r, err := api.BuscarHistoricoResumo(cmd.Context(), g); err == nil {
					fmt.Fprintf(s.err, "Em dia: %s · Guias vencidas: %d\n", simNao(r.EmDia), r.QuantidadeGuiasVencidas)
				}
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&f.Ano, "ano", 0, "ano da competência")
	cmd.Flags().IntVar(&f.Mes, "mes", 0, "mês da competência (1 a 12)")
	cmd.Flags().StringVar(&f.Status, "status", "", "situação da guia (ex.: PAGO, PENDENTE)")
	return cmd
}

func historicoList(guias []api.GuiaHistorico) *output.List {
	l := &output.List{Columns: []output.Column{
		{Key: "competencia", Header: "Competência"},
		{Key: "id", Header: "ID"},
		{Key: "imposto", Header: "Imposto"},
		{Key: "vencimento", Header: "Vencimento"},
		{Key: "valor", Header: "Valor"},
		{Key: "valor_pago", Header: "Valor pago"},
		{Key: "situacao", Header: "Situação"},
		{Key: "tipo", Header: "Tipo"},
	}}
	for _, g := range guias {
		venc, _ := output.ParseDate(g.DataVencimento)
		l.Append(competencia(g.Competencia.Mes, g.Competencia.Ano), g.ID, firstNonEmpty(g.ImpostoDescricao, g.Imposto),
			venc, moneyOrNil(g.ValorPrincipal), moneyOrNil(g.ValorPago), g.Status, strings.ToLower(g.Tipo))
	}
	return l
}

func simNao(b bool) string {
	if b {
		return "sim"
	}
	return "não"
}
