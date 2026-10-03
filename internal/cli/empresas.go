package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/ctbz"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newEmpresasCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "empresas",
		Short: "Lista as empresas do usuário, marcando a atual",
		Long: `Lista todas as empresas vinculadas ao usuário, inclusive as inativas, com a situação
de cada uma. A coluna "atual" marca a empresa da sessão; para trocar, use "ctbz empresa usar".`,
		Example: `  ctbz empresas
  ctbz empresas -o csv`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			d, err := api.BuscarDadosEmpresa(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, empresasList(d))
		},
	}
}

func empresasList(d *api.DadosEmpresa) *output.List {
	atual := ctbz.OnlyDigits(d.EmpresaAtual.CNPJ)
	l := &output.List{Columns: []output.Column{
		{Key: "atual", Header: "Atual"},
		{Key: "cnpj", Header: "CNPJ"},
		{Key: "razao_social", Header: "Razão social"},
		{Key: "situacao", Header: "Situação"},
	}}
	for _, e := range d.Empresas {
		l.Append(ctbz.OnlyDigits(e.CNPJ) == atual, output.NewCNPJ(e.CNPJ), e.RazaoSocial, e.StatusEmpresa)
	}
	return l
}

func newEmpresaUsarCmd() *cobra.Command {
	var o loginOpts
	cmd := &cobra.Command{
		Use:   "usar CNPJ",
		Short: "Troca a empresa da sessão",
		Long: `Troca a empresa da sessão refazendo o login com --cnpj: a Contabilizei escolhe a
empresa só no login, então um novo código OTP é enviado por e-mail. Com CTBZ_OTP_CMD
definido, a troca é automática; sem ele, o código é pedido no terminal ou a troca fica
pendente para "ctbz login --otp NNNNNN".`,
		Example: `  ctbz empresa usar 00.000.000/0001-00`,
		Args:    exactArgs(1, "o CNPJ da empresa"),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := streamsOf(cmd)
			o.cnpj = ctbz.OnlyDigits(args[0])
			if len(o.cnpj) != 14 {
				return usageError{fmt.Errorf("CNPJ inválido: %q", args[0])}
			}
			store, err := ctbz.DefaultStore()
			if err != nil {
				return err
			}
			if sess, err := store.LoadSession(); err == nil && sess.CNPJ == o.cnpj {
				fmt.Fprintf(s.err, "A sessão já está na empresa %s.\n", output.FormatCNPJ(output.CNPJ(o.cnpj)))
				return nil
			}
			o.restart = true
			o.applyEnv(cmd)
			sess, err := login(cmd.Context(), store, s, o)
			if err != nil {
				return err
			}
			fmt.Fprintf(s.err, "Empresa atual: %s\n", describeSession(sess))
			return nil
		},
	}
	return cmd
}
