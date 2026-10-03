package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newContasCmd() *cobra.Command {
	var busca, situacao string
	cmd := &cobra.Command{
		Use:   "contas",
		Short: "Lista o plano de contas usado para classificar os lançamentos",
		Long: `Lista as contas usadas para classificar entradas e saídas (as mesmas da tela de
classificação de extratos): descrição, conta contábil correspondente, classificação (receita,
despesa…) e situação.

--busca filtra por um trecho da descrição ou da conta contábil, sem diferenciar maiúsculas
(acentos contam: use um trecho como "alug" para achar "Aluguel" e "Aluguéis"). --situacao
filtra por ativo ou inativo; por padrão aparecem todas.`,
		Example: `  ctbz contas --busca alug
  ctbz contas --situacao ativo -o csv`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			situacao = strings.ToUpper(situacao)
			if situacao != "" && situacao != "ATIVO" && situacao != "INATIVO" {
				return usageError{fmt.Errorf("--situacao deve ser ativo ou inativo: %q", situacao)}
			}
			s := streamsOf(cmd)
			cs, err := api.BuscarContasUsuario(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, contasList(cs, busca, situacao))
		},
	}
	cmd.Flags().StringVar(&busca, "busca", "", "trecho da descrição ou da conta contábil")
	cmd.Flags().StringVar(&situacao, "situacao", "", "ativo ou inativo (padrão: todas)")
	return cmd
}

func contasList(cs []api.ContaUsuario, busca, situacao string) *output.List {
	l := &output.List{Columns: []output.Column{
		{Key: "descricao", Header: "Descrição"},
		{Key: "conta_contabil", Header: "Conta contábil"},
		{Key: "classificacao", Header: "Classificação"},
		{Key: "situacao", Header: "Situação"},
		{Key: "id", Header: "ID"},
	}}
	busca = strings.ToLower(busca)
	for _, c := range cs {
		if situacao != "" && c.Situacao != situacao {
			continue
		}
		if busca != "" && !strings.Contains(strings.ToLower(c.Descricao+" "+c.DescricaoContaContabil), busca) {
			continue
		}
		l.Append(c.Descricao, nilIfEmpty(c.DescricaoContaContabil), nilIfEmpty(c.Classificacao), nilIfEmpty(c.Situacao), c.ID)
	}
	return l
}
