package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newImpostosRecorrenteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recorrente",
		Short: "Mostra a situação do pagamento recorrente (débito automático) de impostos",
		Long: `Mostra se o pagamento recorrente de impostos (débito automático no cartão ou na conta PJ)
está disponível e ativo, a competência atual, as datas do próximo pagamento e da próxima
tentativa e quantos pagamentos estão agendados, concluídos e recusados.

Dados de cartão e chaves de pagamento nunca são mostrados: só a quantidade de cartões salvos.
Para os meses anteriores, use "ctbz impostos recorrente historico".`,
		Example: `  ctbz impostos recorrente
  ctbz impostos recorrente historico -o csv`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			r, err := api.BuscarRecorrencia(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, recorrenciaRecord(r))
		},
	}
	cmd.AddCommand(newImpostosRecorrenteHistoricoCmd())
	return cmd
}

func recorrenciaRecord(r *api.Recorrencia) *output.Record {
	rec := &output.Record{}
	rec.Add("disponivel", "Disponível", r.Habilitado)
	rec.Add("ativo", "Ativo", r.Ativado)
	rec.Add("situacao", "Situação", nilIfEmpty(r.Status))
	rec.Add("competencia", "Competência", nilIfEmpty(strings.ReplaceAll(r.Competencia, "-", "/")))
	rec.Add("proximo_pagamento", "Próximo pagamento", dataOuTexto(r.DataProximoPagamento))
	rec.Add("proxima_tentativa", "Próxima tentativa", dataOuTexto(r.DataProximaTentativa))
	rec.Add("cartoes_salvos", "Cartões salvos", len(r.CartoesSalvos))
	rec.Add("pagamentos_agendados", "Pagamentos agendados", len(r.PagamentosAgendados))
	rec.Add("pagamentos_concluidos", "Pagamentos concluídos", len(r.PagamentosConcluido))
	rec.Add("pagamentos_recusados", "Pagamentos recusados", len(r.PagamentosRecusado))
	return rec
}

func newImpostosRecorrenteHistoricoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "historico",
		Short: "Lista os pagamentos recorrentes de impostos dos últimos meses",
		Long: `Lista, por competência, as guias pagas pelo pagamento recorrente, com situação e valor,
e o custo de operação de cada mês. A Contabilizei só devolve um período recente; para os
meses anteriores, use "ctbz impostos historico".`,
		Example: `  ctbz impostos recorrente historico`,
		Args:    exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			h, err := api.BuscarHistoricoRecorrencia(cmd.Context(), sessionGetter{s})
			if err != nil {
				return err
			}
			return output.Write(s.out, f, historicoRecorrenciaList(h))
		},
	}
}

func historicoRecorrenciaList(h *api.HistoricoRecorrencia) *output.List {
	l := &output.List{Columns: []output.Column{
		{Key: "competencia", Header: "Competência"},
		{Key: "item", Header: "Item"},
		{Key: "situacao", Header: "Situação"},
		{Key: "valor", Header: "Valor"},
	}}
	for _, p := range h.Pagamentos {
		comp := competenciaDeValores(p.Mes, p.Ano)
		for _, g := range p.Guias {
			l.Append(comp, g.Nome, nilIfEmpty(g.Status), moneyOrNil(g.Valor))
		}
		if p.CustoOperacao != nil {
			l.Append(comp, "Custo de operação", nil, output.Money(*p.CustoOperacao))
		}
	}
	return l
}

// competenciaDeValores monta "MM/AAAA" com mês e ano que podem vir como número ou texto.
func competenciaDeValores(mes, ano any) any {
	if mes == nil || ano == nil {
		return nil
	}
	m := fmt.Sprint(mes)
	if len(m) == 1 {
		m = "0" + m
	}
	return m + "/" + fmt.Sprint(ano)
}
