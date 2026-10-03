package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/ctbz"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Mostra a sessão atual e testa se ainda é válida",
		Args:  exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			store, err := ctbz.DefaultStore()
			if err != nil {
				return err
			}
			sess, err := store.LoadSession()
			if err != nil {
				return err
			}
			situacao := "válida"
			c := sess.Client()
			if _, err := c.API(cmd.Context(), "GET", "appbar/get", nil); err != nil {
				if !errors.Is(err, ctbz.ErrUnauthorized) {
					return err
				}
				situacao = "expirada"
				fmt.Fprintln(s.err, "Sessão expirada: rode `ctbz login`.")
			} else {
				saveCookies(store, sess, c, s.err)
			}
			info := sessionInfo(sess)
			rec := &output.Record{}
			rec.Add("empresa", "Empresa", info.RazaoSocial).
				Add("cnpj", "CNPJ", output.NewCNPJ(sess.CNPJ)).
				Add("usuario", "Usuário", info.Email).
				Add("login_em", "Login em", output.DateTime{Time: sess.CreatedAt}).
				Add("expira_em", "Expira em", output.DateTime{Time: sess.ExpiresAt()}).
				Add("situacao", "Situação", situacao)
			return output.Write(s.out, f, rec)
		},
	}
}
