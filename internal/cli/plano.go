package cli

import (
	"errors"
	"io"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newPlanoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plano",
		Short: "Mostra o plano contratado com a Contabilizei",
		Long: `Mostra o plano contratado com a Contabilizei: descrição, categoria, valor de tabela e ramos
de atividade cobertos. Os dados vêm do login (não fazem chamada à API); depois de uma
mudança de plano, rode "ctbz login" de novo.

Para o texto do contrato de serviço e da proposta do plano, use "ctbz plano contrato" e
"ctbz plano proposta".`,
		Example: `  ctbz plano
  ctbz plano contrato --texto > contrato.txt
  ctbz plano proposta > proposta.html`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			e, err := empresaDaSessao()
			if err != nil {
				return err
			}
			if e.PlanoPagamentoEmpresa == nil {
				return errors.New("a empresa não tem plano ativo nos dados do login")
			}
			p := e.PlanoPagamentoEmpresa.PagtoPlano
			rec := &output.Record{}
			rec.Add("plano", "Plano", nilIfEmpty(p.Descricao))
			rec.Add("categoria", "Categoria", nilIfEmpty(p.Categoria))
			rec.Add("valor", "Valor de tabela", moneyOrNil(p.Valor))
			rec.Add("ramos_atividade", "Ramos de atividade", nonNil(p.RamoAtividades))
			return output.Write(streamsOf(cmd).out, f, rec)
		},
	}
	cmd.AddCommand(newPlanoContratoCmd(), newPlanoPropostaCmd())
	return cmd
}

func newPlanoContratoCmd() *cobra.Command {
	var texto bool
	cmd := &cobra.Command{
		Use:   "contrato",
		Short: "Exporta o contrato de prestação de serviços (HTML ou texto)",
		Long: `Imprime no stdout o contrato de prestação de serviços aceito pela empresa, em HTML (como o
painel mostra) ou, com --texto, em texto simples. Redirecione para um arquivo para guardar.`,
		Example: `  ctbz plano contrato > contrato.html
  ctbz plano contrato --texto | less`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s := streamsOf(cmd)
			c, err := api.BuscarContratoServico(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return escreverDocumento(s.out, c.HTML, texto)
		},
	}
	cmd.Flags().BoolVar(&texto, "texto", false, "converte o HTML em texto simples")
	return cmd
}

func newPlanoPropostaCmd() *cobra.Command {
	var texto bool
	cmd := &cobra.Command{
		Use:   "proposta",
		Short: "Exporta a proposta do plano contratado, com a tabela de preços (HTML ou texto)",
		Long: `Imprime no stdout a proposta de prestação de serviços do plano contratado (tabela de
mensalidades por faixa de faturamento e o plano selecionado), em HTML ou, com --texto, em
texto simples.`,
		Example: `  ctbz plano proposta --texto`,
		Args:    exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			s := streamsOf(cmd)
			doc, err := api.BuscarPropostaPlano(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return escreverDocumento(s.out, doc, texto)
		},
	}
	cmd.Flags().BoolVar(&texto, "texto", false, "converte o HTML em texto simples")
	return cmd
}

// escreverDocumento imprime um documento HTML como veio ou convertido em texto.
func escreverDocumento(w io.Writer, doc string, texto bool) error {
	if doc == "" {
		return errors.New("a Contabilizei devolveu um documento vazio")
	}
	if texto {
		doc = htmlParaTexto(doc)
	}
	_, err := io.WriteString(w, doc)
	return err
}
