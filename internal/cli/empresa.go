package cli

import (
	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/ctbz"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newEmpresaCmd() *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "empresa",
		Short: "Mostra os dados da empresa selecionada",
		Long: `Mostra razão social, CNPJ, situação, regime tributário, plano, certificado digital
e as outras empresas do usuário.

A resposta crua da API está em "ctbz api dadosempresa/get".`,
		Example: `  ctbz empresa
  ctbz empresa -o json | jq -r .certificado_validade`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			fallback := ""
			if asJSON {
				fallback = string(output.FormatJSON)
			}
			f, err := outputFormat(cmd, fallback)
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			data, err := api.BuscarDadosEmpresa(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, empresaRecord(data))
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "atalho para -o json")
	return cmd
}

// empresaRecord monta a saída de "ctbz empresa".
func empresaRecord(d *api.DadosEmpresa) *output.Record {
	e := d.EmpresaAtual
	var certSituacao any
	var certValidade output.Date
	if e.Certificado != nil {
		certSituacao = e.Certificado.Status.Descricao
		certValidade, _ = output.ParseDate(e.Certificado.DataValidade)
	}
	outras := []output.Record{}
	for _, o := range d.Empresas {
		if ctbz.OnlyDigits(o.CNPJ) == ctbz.OnlyDigits(e.CNPJ) {
			continue
		}
		r := output.Record{}
		r.Add("cnpj", "CNPJ", output.NewCNPJ(o.CNPJ)).
			Add("razao_social", "Razão social", o.RazaoSocial).
			Add("situacao", "Situação", o.StatusEmpresa)
		outras = append(outras, r)
	}
	rec := &output.Record{}
	rec.Add("razao_social", "Razão social", e.RazaoSocial).
		Add("cnpj", "CNPJ", output.NewCNPJ(e.CNPJ)).
		Add("situacao", "Situação", e.StatusEmpresa).
		Add("regime_tributario", "Regime tributário", e.RegimeTributario).
		Add("inscricao_municipal", "Inscrição municipal", nilIfEmpty(e.InscricaoMunicipal)).
		Add("ramos_atividade", "Ramos de atividade", nonNil(e.RamosAtividade)).
		Add("plano", "Plano", nilIfEmpty(e.Plano)).
		Add("certificado_situacao", "Certificado digital", certSituacao).
		Add("certificado_validade", "Validade do certificado", certValidade).
		Add("outras_empresas", "Outras empresas", outras)
	return rec
}
