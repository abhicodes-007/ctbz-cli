package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/ctbz"
)

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Apaga a sessão local",
		Long: `Apaga a sessão e qualquer login pendente salvos em CTBZ_HOME.

A Contabilizei não tem logout no servidor: a sessão continua válida lá até expirar.`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			store, err := ctbz.DefaultStore()
			if err != nil {
				return err
			}
			if err := store.Clear(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.ErrOrStderr(), "Sessão local removida.")
			return nil
		},
	}
}
