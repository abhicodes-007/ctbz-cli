package cli

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

// newEmpresaCertificadoCmd é "ctbz empresa certificado"; o mesmo comando existe como
// "ctbz certificado" (newCertificadoCmd).
func newEmpresaCertificadoCmd() *cobra.Command {
	return newCertificadoCmdComo("ctbz empresa certificado")
}

func newCertificadoCmd() *cobra.Command { return newCertificadoCmdComo("ctbz certificado") }

func newCertificadoCmdComo(nome string) *cobra.Command {
	return &cobra.Command{
		Use:   "certificado",
		Short: "Mostra a situação do certificado digital da empresa e da renovação",
		Long: `Mostra a situação do certificado digital (e-CNPJ) usado pela Contabilizei, a data de
vencimento, os dias que faltam, se já é possível renovar, o alerta de vencimento do painel e
a etapa do processo de compra ou renovação feito pela Contabilizei (ex.: PRE_CHECKOUT antes
de começar, agendamentoVideoconferencia durante a validação).

A senha e o arquivo do certificado nunca são consultados.`,
		Example: "  " + nome + "\n  " + nome + " -o json | jq .dias_para_vencer",
		Args:    exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			g := sessionGetter{s}
			c, err := api.BuscarCertificadoStatus(cmd.Context(), g)
			if err != nil {
				return err
			}
			p, err := api.BuscarCertificadoProcesso(cmd.Context(), g)
			if err != nil {
				return err
			}
			card, err := api.BuscarCertificadoCard(cmd.Context(), g)
			if err != nil {
				return err
			}
			rec := certificadoRecord(c, now())
			rec.Add("vencido_no_painel", "Vencido (painel)", card.CertificadoVencido).
				Add("prazo_finalizado", "Prazo finalizado", card.PrazoFinalizado).
				Add("renovacao_etapa", "Etapa da renovação", nilIfEmpty(p.EtapaAtual)).
				Add("renovacao_fluxo", "Fluxo da renovação", nilIfEmpty(p.FluxoDestino))
			return output.Write(s.out, f, rec)
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
