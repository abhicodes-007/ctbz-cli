package cli

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newEmpresaCertificadoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "certificado",
		Short: "Mostra a situação do certificado digital da empresa",
		Long: `Mostra a situação do certificado digital (e-CNPJ) usado pela Contabilizei, a data de
vencimento, os dias que faltam e se já é possível renovar.`,
		Example: `  ctbz empresa certificado
  ctbz empresa certificado -o json | jq .dias_para_vencer`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			c, err := api.BuscarCertificadoStatus(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, certificadoRecord(c, now()))
		},
	}
}

func certificadoRecord(c *api.CertificadoStatus, hoje time.Time) *output.Record {
	var vencimento output.Date
	var dias any
	if c.DataVencimento != nil {
		vencimento = dateFromMillis(*c.DataVencimento)
		dias = diasEntre(hoje, vencimento.Time)
	}
	rec := &output.Record{}
	rec.Add("situacao", "Situação", c.Situacao).
		Add("valido", "Válido", c.Valido).
		Add("vencimento", "Vencimento", vencimento).
		Add("dias_para_vencer", "Dias para vencer", dias).
		Add("apto_renovacao", "Pode renovar", c.AptoRenovacao).
		Add("mensagem", "Mensagem", nilIfEmpty(c.MensagemErro))
	return rec
}
