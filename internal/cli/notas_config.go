package cli

import (
	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newNotasConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Mostra a configuração do emissor de notas",
		Long: `Mostra se o emissor de notas está habilitado, a versão em uso, se a Contabilizei reporta
instabilidade com a prefeitura, a situação do certificado digital usado na emissão, se é
permitido emitir para o exterior e o município da empresa.

As alíquotas por atividade estão em "ctbz notas aliquotas".`,
		Example: `  ctbz notas config`,
		Args:    exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			g := sessionGetter{s}
			ini, err := api.BuscarEmissorInit(cmd.Context(), g)
			if err != nil {
				return err
			}
			v, err := api.BuscarVersaoEmissor(cmd.Context(), g)
			if err != nil {
				return err
			}
			return output.Write(s.out, f, emissorRecord(ini, v))
		},
	}
}

func emissorRecord(ini *api.EmissorInit, v *api.VersaoEmissor) *output.Record {
	rec := &output.Record{}
	rec.Add("emissor_habilitado", "Emissor habilitado", ini.EmissorEnabled)
	rec.Add("versao", "Versão", nilIfEmpty(v.VersaoNovoEmissor))
	rec.Add("instabilidade", "Instabilidade reportada", ini.HasInstability)
	rec.Add("permite_exterior", "Emite para o exterior", ini.PermiteEmissaoExterior)
	var validade, vencido, renovacao any
	if c := ini.CertificadoDigital; c != nil {
		validade, vencido, renovacao = dataOuTexto(c.DataVencimento), c.Vencido, c.EmRenovacao
	}
	rec.Add("certificado_validade", "Validade do certificado", validade)
	rec.Add("certificado_vencido", "Certificado vencido", vencido)
	rec.Add("certificado_em_renovacao", "Certificado em renovação", renovacao)
	var municipio, uf any
	if e := ini.Endereco; e != nil && e.Municipio != nil {
		municipio, uf = nilIfEmpty(e.Municipio.Nome), nilIfEmpty(e.Municipio.UF.ID)
	}
	rec.Add("municipio", "Município", municipio)
	rec.Add("uf", "UF", uf)
	return rec
}

func newNotasAliquotasCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "aliquotas",
		Short: "Lista as alíquotas e os códigos de serviço por atividade",
		Long: `Lista, para cada atividade (CNAE) da empresa, o item da lista de serviços usado na nota, a
alíquota do Simples (%), a parte que é ISS (%), o Fator R considerado (%) e se o anexo é
fixo. As alíquotas são separadas por mercado: tomadores no Brasil (interno) e no exterior
(externo, sem ISS).`,
		Example: `  ctbz notas aliquotas
  ctbz notas aliquotas -o csv`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			a, err := api.BuscarAliquotasEmissor(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, aliquotasList(a))
		},
	}
}

func aliquotasList(a *api.AliquotasEmissor) *output.List {
	l := &output.List{Columns: []output.Column{
		{Key: "mercado", Header: "Mercado"},
		{Key: "tipo", Header: "Tipo"},
		{Key: "cnae", Header: "CNAE"},
		{Key: "atividade", Header: "Atividade"},
		{Key: "item_servico", Header: "Item"},
		{Key: "descricao_item", Header: "Descrição do item"},
		{Key: "aliquota", Header: "Alíquota (%)"},
		{Key: "aliquota_iss", Header: "ISS (%)"},
		{Key: "fator_r", Header: "Fator R (%)"},
		{Key: "anexo_fixo", Header: "Anexo fixo"},
	}}
	for _, m := range []struct {
		nome string
		a    api.AliquotasMercado
	}{{"interno", a.Interno}, {"externo", a.Externo}} {
		for _, t := range []struct {
			nome  string
			itens []api.AliquotaAtividade
		}{{"servico", m.a.Servico}, {"comercio", m.a.Comercio}} {
			for _, it := range t.itens {
				l.Append(m.nome, t.nome, formatCNAE(it.CodigoCnae), output.Text(it.DescricaoCnae), nilIfEmpty(it.CodigoItemServico),
					output.Text(it.DescricaoItemServico), floatOrNil(it.AliquotaBase), floatOrNil(it.AliquotaISS),
					floatOrNil(it.FatorR), it.AnexoFixo)
			}
		}
	}
	return l
}
