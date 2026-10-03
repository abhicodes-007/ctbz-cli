package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newBalanceteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "balancete AAAA-MM",
		Short: "Mostra o balancete de verificação de um mês",
		Long: `Mostra o balancete do mês: cada conta da árvore contábil com saldo anterior, débitos,
créditos e saldo do exercício. Na tabela, as contas aparecem recuadas por nível; em CSV e
JSON, a descrição vem sem recuo e o nível fica na coluna "nivel" (bom para planilhas).

O padrão do painel é dezembro do ano anterior; aqui o mês é obrigatório.`,
		Example: `  ctbz balancete 2026-08
  ctbz balancete 2026-08 -o csv > balancete-2026-08.csv`,
		Args: exactArgs(1, "o mês (AAAA-MM)"),
		RunE: func(cmd *cobra.Command, args []string) error {
			mes, err := parseMes(args[0])
			if err != nil {
				return err
			}
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			contas, err := api.BuscarRelatorio(cmd.Context(), sessionGetter{s}, api.PathBalancete(mes.Year(), int(mes.Month())))
			if err != nil {
				return err
			}
			return output.Write(s.out, f, balanceteList(contas))
		},
	}
}

func balanceteList(contas []api.ContaRelatorio) *output.List {
	l := &output.List{Columns: []output.Column{
		{Key: "conta", Header: "Conta"},
		{Key: "descricao", Header: "Descrição"},
		{Key: "nivel", Header: "Nível"},
		{Key: "saldo_anterior", Header: "Saldo anterior"},
		{Key: "debitos", Header: "Débitos"},
		{Key: "creditos", Header: "Créditos"},
		{Key: "saldo", Header: "Saldo"},
	}}
	for _, c := range contas {
		l.Append(c.ID, output.Indent{Level: c.Nivel, Text: c.Descricao}, c.Nivel, moneyOrNil(c.SaldoAnterior),
			moneyOrNil(c.TotalDebito), moneyOrNil(c.TotalCredito), moneyOrNil(c.SaldoExercicio))
	}
	return l
}

// parseMes lê um mês AAAA-MM como argumento.
func parseMes(s string) (time.Time, error) {
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return t, usageError{fmt.Errorf("mês inválido %q: use AAAA-MM", s)}
	}
	return t, nil
}
