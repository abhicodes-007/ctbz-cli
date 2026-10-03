package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newMensalidadeHistoricoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "historico",
		Short: "Lista os pagamentos anteriores da mensalidade e a situação do débito automático",
		Long: `Lista os pagamentos anteriores da mensalidade (data, competência, valor e status) e mostra a
situação do débito automático: se está habilitado e ativo, a competência e a próxima cobrança.

Só campos conhecidos são copiados para a saída: chaves do gateway de pagamento e dados de cartão
nunca são mostrados. O formato de cada pagamento não é documentado; se nenhum campo for
reconhecido, o comando avisa e indica "ctbz api payments/recorrencia/historico".`,
		Example: `  ctbz mensalidade historico
  ctbz mensalidade historico -o json | jq '.pagamentos[].valor'`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			g := sessionGetter{s}
			var ini, hist json.RawMessage
			if err := g.GetJSON(cmd.Context(), api.PathRecorrenciaInit, &ini); err != nil {
				return err
			}
			if err := g.GetJSON(cmd.Context(), api.PathRecorrenciaHistorico, &hist); err != nil {
				return err
			}
			rec, err := buildHistorico(ini, hist, s.err)
			if err != nil {
				return err
			}
			return output.Write(s.out, f, rec)
		},
	}
}

// buildHistorico monta o registro de saída a partir das respostas de
// payments/recorrencia/init e payments/recorrencia/historico. Só campos
// conhecidos são copiados: a resposta pode trazer dados de cartão ou chaves
// do gateway de pagamento (adyenClientKey), que nunca devem chegar à saída.
// Avisos vão para errOut.
func buildHistorico(initBody, histBody []byte, errOut io.Writer) (*output.Record, error) {
	var ini map[string]any
	if err := decodeJSON(initBody, &ini); err != nil {
		return nil, fmt.Errorf("resposta inesperada de %s: %w", api.PathRecorrenciaInit, err)
	}
	var hist struct {
		Pagamentos []map[string]any `json:"pagamentos"`
	}
	if err := decodeJSON(histBody, &hist); err != nil {
		return nil, fmt.Errorf("resposta inesperada de %s: %w", api.PathRecorrenciaHistorico, err)
	}

	pagamentos := make([]output.Record, 0, len(hist.Pagamentos))
	reconhecidos := false
	for _, p := range hist.Pagamentos {
		r := output.Record{}
		data, okData := pick(p, "dataPagamento", "dataPagto", "data", "dataVencimento", "vencimento")
		competencia, okComp := pick(p, "competencia", "mesCompetencia")
		valor, okValor := pick(p, "valorPago", "valor", "total")
		status, okStatus := pick(p, "status", "situacao")
		reconhecidos = reconhecidos || okData || okComp || okValor || okStatus
		r.Add("data", "Data", toDate(data)).
			Add("competencia", "Competência", toText(competencia)).
			Add("valor", "Valor", toMoney(valor)).
			Add("status", "Status", toText(status))
		pagamentos = append(pagamentos, r)
	}
	if len(hist.Pagamentos) > 0 && !reconhecidos {
		fmt.Fprintln(errOut, "aviso: formato de pagamentos não reconhecido; veja o original com `ctbz api payments/recorrencia/historico`")
	}

	status, _ := pick(ini, "status")
	competencia, _ := pick(ini, "competencia")
	proximo, _ := pick(ini, "dataProximoPagamento")
	habilitado, _ := pick(ini, "habilitado")
	ativado, _ := pick(ini, "ativado")

	rec := &output.Record{}
	rec.Add("debito_automatico_ativo", "Débito automático ativo", toBool(ativado)).
		Add("debito_automatico_habilitado", "Débito automático habilitado", toBool(habilitado)).
		Add("status", "Status", toText(status)).
		Add("competencia", "Competência", toText(competencia)).
		Add("proximo_pagamento", "Próxima cobrança", toDate(proximo)).
		Add("pagamentos", "Pagamentos", pagamentos)
	return rec, nil
}

func decodeJSON(b []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v)
}

// pick devolve o primeiro valor não nulo entre as chaves candidatas.
func pick(m map[string]any, keys ...string) (any, bool) {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			return v, true
		}
	}
	return nil, false
}

func toText(v any) any {
	switch x := v.(type) {
	case nil:
		return nil
	case string:
		return nilIfEmpty(x)
	case json.Number:
		return x.String()
	case bool:
		return strconv.FormatBool(x)
	default:
		return nil
	}
}

func toBool(v any) any {
	b, ok := v.(bool)
	if !ok {
		return nil
	}
	return b
}

// toDate aceita texto (dd/mm/aaaa ou ISO) ou epoch em milissegundos.
func toDate(v any) output.Date {
	switch x := v.(type) {
	case string:
		d, _ := output.ParseDate(x)
		return d
	case json.Number:
		if ms, err := x.Int64(); err == nil {
			return output.DateFromMillis(ms)
		}
	}
	return output.Date{}
}

// toMoney aceita número ou texto ("1.234,56" ou "1234.56"). Sem valor: nil.
func toMoney(v any) any {
	switch x := v.(type) {
	case json.Number:
		if f, err := x.Float64(); err == nil {
			return output.Money(f)
		}
	case string:
		s := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(x), "R$"))
		if strings.Contains(s, ",") {
			s = strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", ".")
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return output.Money(f)
		}
	}
	return nil
}
