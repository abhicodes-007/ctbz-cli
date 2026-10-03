package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newEmpresaSociosCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "socios",
		Short: "Lista os sócios da empresa",
		Long: `Lista os sócios com CPF, se são administradores ou responsáveis perante a Receita,
categoria (com ou sem pró-labore), salário-base e data de entrada.`,
		Example: `  ctbz empresa socios
  ctbz empresa socios -o json`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			socios, err := api.BuscarSocios(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, sociosList(socios))
		},
	}
}

func sociosList(socios []api.Socio) *output.List {
	l := &output.List{Columns: []output.Column{
		{Key: "nome", Header: "Nome"},
		{Key: "cpf", Header: "CPF"},
		{Key: "administrador", Header: "Administrador"},
		{Key: "responsavel_receita", Header: "Responsável na Receita"},
		{Key: "categoria", Header: "Categoria"},
		{Key: "salario_base", Header: "Salário-base"},
		{Key: "entrada", Header: "Entrada"},
		{Key: "situacao", Header: "Situação"},
	}}
	for _, s := range socios {
		l.Append(s.Nome, output.NewCPF(s.CPF), s.Administrador, s.ResponsavelReceita, s.Categoria.Descricao,
			output.Money(s.SalarioBase), dateFromMillis(s.DataAdmissao), s.SituacaoColaborador.Descricao)
	}
	return l
}

func newEmpresaAtividadesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "atividades",
		Short: "Lista os CNAEs da empresa e os anexos do Simples Nacional",
		Long: `Lista as atividades (CNAEs) cadastradas na Contabilizei, marcando a principal, com o
ramo (serviço ou comércio) e os anexos do Simples Nacional em que cada uma pode ser tributada.`,
		Example: `  ctbz empresa atividades`,
		Args:    exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			cnaes, err := api.BuscarCNAEs(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, atividadesList(cnaes))
		},
	}
}

func atividadesList(cnaes []api.CNAEEmpresa) *output.List {
	l := &output.List{Columns: []output.Column{
		{Key: "cnae", Header: "CNAE"},
		{Key: "descricao", Header: "Descrição"},
		{Key: "principal", Header: "Principal"},
		{Key: "ramo", Header: "Ramo"},
		{Key: "anexos_simples", Header: "Anexos do Simples"},
	}}
	for _, c := range cnaes {
		principal := false
		anexos := []string{}
		for _, a := range c.Anexos {
			if !a.Ativo {
				continue
			}
			principal = principal || a.Principal
			anexos = append(anexos, anexoRomano(a.CodTabelaSimples))
		}
		l.Append(formatCNAE(c.CNAE.Codigo), c.CNAE.Descricao, principal, strings.ToLower(c.CNAE.TipoRamoAtividade), anexos)
	}
	return l
}

// formatCNAE aplica a máscara oficial 0000-0/00 a um código de 7 dígitos.
func formatCNAE(codigo string) string {
	if len(codigo) != 7 {
		return codigo
	}
	return codigo[:4] + "-" + codigo[4:5] + "/" + codigo[5:]
}

// anexoRomano converte o número do anexo do Simples Nacional (1 a 5) em algarismo romano.
func anexoRomano(n int) string {
	if r := []string{"", "I", "II", "III", "IV", "V"}; n >= 1 && n <= 5 {
		return r[n]
	}
	return "?"
}
