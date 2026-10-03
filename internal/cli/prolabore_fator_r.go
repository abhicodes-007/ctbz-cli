package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newProlaboreFatorRCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fator-r",
		Short: "Mostra a situação do Fator R e os anexos possíveis de cada atividade",
		Long: `Mostra se a Contabilizei ajusta o pró-labore para o Fator R (motor do Fator R), o Fator R
atual (pró-labore ÷ faturamento dos últimos 12 meses), os valores acumulados usados no
cálculo e, para cada atividade, os anexos do Simples Nacional em que pode ser tributada.

Com Fator R a partir de 28%, as atividades sujeitas a ele saem do Anexo V (alíquota inicial
de 15,5%) para o Anexo III (6%). Os dados vêm do simulador de impostos avançado do painel
(só leitura) e da memória de cálculo do mês.`,
		Example: `  ctbz prolabore fator-r
  ctbz prolabore fator-r -o json | jq .fator_r`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			g := sessionGetter{s}
			p, err := api.BuscarProlaboreParametros(cmd.Context(), g)
			if err != nil {
				return err
			}
			sim, err := api.BuscarSimuladorImpostos(cmd.Context(), g)
			if err != nil {
				return err
			}
			c, err := api.BuscarCalculoImposto(cmd.Context(), g)
			if err != nil {
				return err
			}
			return output.Write(s.out, f, fatorRRecord(p, sim, c))
		},
	}
}

func fatorRRecord(p *api.ProlaboreParametros, sim *api.SimuladorImpostos, c *api.CalculoImposto) *output.Record {
	atividades := make([]output.Record, 0, len(sim.Atividades))
	for _, a := range sim.Atividades {
		var anexos, aliquotas []string
		fixo := false
		for _, an := range a.Anexos {
			anexos = append(anexos, anexoRomano(an.Anexo))
			aliquotas = append(aliquotas, an.DescricaoAliquota)
			fixo = fixo || an.AnexoFixo
		}
		r := output.Record{}
		r.Add("cnae", "CNAE", a.Codigo)
		r.Add("atividade", "Atividade", output.Text(a.Descricao))
		r.Add("anexos", "Anexos", strings.Join(anexos, " ou "))
		r.Add("aliquotas", "Alíquotas", output.Text(strings.Join(aliquotas, "; ")))
		r.Add("anexo_fixo", "Anexo fixo", fixo)
		atividades = append(atividades, r)
	}
	rec := &output.Record{}
	rec.Add("motor_fator_r", "Motor do Fator R", p.IsMotorFatorR || sim.MotorFatorR)
	rec.Add("fator_r", "Fator R (%)", floatOrNil(c.PercentualFatorR))
	rec.Add("prolabore_12_meses", "Pró-labore 12 meses", moneyOrNil(c.ValorProlaboreUltimos12Meses))
	rec.Add("faturamento_12_meses", "Faturamento 12 meses", moneyOrNil(c.ValorFaturamentoUltimos12Meses))
	rec.Add("simulador", "Simulador", nilIfEmpty(sim.Disponibilidade))
	rec.Add("atividades", "Atividades", atividades)
	return rec
}
