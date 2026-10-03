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

func newBalancoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "balanco AAAA[-MM]",
		Short: "Mostra o balanço patrimonial (ativo, passivo e patrimônio líquido)",
		Long: `Mostra o balanço patrimonial: as contas de ativo, passivo e patrimônio líquido com o saldo
do exercício e o do exercício anterior, recuadas por nível na tabela. As contas de resultado
ficam de fora, como no painel.

Só com o ano, o balanço é o de dezembro (fechamento do exercício).`,
		Example: `  ctbz balanco 2025
  ctbz balanco 2026-09 -o csv`,
		Args: exactArgs(1, "o ano (AAAA) ou o mês (AAAA-MM)"),
		RunE: func(cmd *cobra.Command, args []string) error {
			mes, err := parseMes(args[0])
			if err != nil {
				if mes, err = time.Parse("2006", args[0]); err != nil {
					return usageError{fmt.Errorf("período inválido %q: use AAAA ou AAAA-MM", args[0])}
				}
				mes = time.Date(mes.Year(), time.December, 1, 0, 0, 0, 0, time.UTC)
			}
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			contas, err := api.BuscarRelatorio(cmd.Context(), sessionGetter{s}, api.PathBalanco(mes.Year(), int(mes.Month())))
			if err != nil {
				return err
			}
			return output.Write(s.out, f, balancoList(contas))
		},
	}
}

// classificacaoResultado são as contas de resultado, que o balanço não mostra.
const classificacaoResultado = "RESULTADO"

func balancoList(contas []api.ContaRelatorio) *output.List {
	l := &output.List{Columns: []output.Column{
		{Key: "conta", Header: "Conta"},
		{Key: "descricao", Header: "Descrição"},
		{Key: "nivel", Header: "Nível"},
		{Key: "grupo", Header: "Grupo"},
		{Key: "saldo", Header: "Saldo"},
		{Key: "saldo_exercicio_anterior", Header: "Exercício anterior"},
	}}
	for _, c := range contas {
		if c.ClassificacaoConta == classificacaoResultado {
			continue
		}
		l.Append(c.ID, output.Indent{Level: c.Nivel, Text: c.Descricao}, c.Nivel, nilIfEmpty(c.ClassificacaoConta),
			moneyOrNil(c.SaldoExercicio), moneyOrNil(c.SaldoExercicioAnterior))
	}
	return l
}
