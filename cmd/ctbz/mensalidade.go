package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/edusouza/ctbz-cli/internal/ctbz"
	"github.com/edusouza/ctbz-cli/internal/output"
)

const mensalidadeUsage = `Uso: ctbz mensalidade SUBCOMANDO [flags]

Subcomandos:
  historico   pagamentos anteriores e situação do débito automático
`

func cmdMensalidade(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, mensalidadeUsage)
		return errors.New("informe o subcomando")
	}
	switch args[0] {
	case "historico":
		return cmdMensalidadeHistorico(ctx, args[1:])
	case "help", "-h", "--help":
		fmt.Print(mensalidadeUsage)
		return nil
	default:
		fmt.Fprint(os.Stderr, mensalidadeUsage)
		return fmt.Errorf("subcomando desconhecido: %s", args[0])
	}
}

func cmdMensalidadeHistorico(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("mensalidade historico", flag.ContinueOnError)
	format := addOutputFlag(fs, defaultFormat())
	if err := fs.Parse(args); err != nil {
		return err
	}
	f, err := output.ParseFormat(*format)
	if err != nil {
		return err
	}
	recInit, err := authedGET(ctx, "payments/recorrencia/init")
	if err != nil {
		return err
	}
	hist, err := authedGET(ctx, "payments/recorrencia/historico")
	if err != nil {
		return err
	}
	rec, err := buildHistorico(recInit, hist)
	if err != nil {
		return err
	}
	return output.Write(os.Stdout, f, rec)
}

// authedGET chama um caminho do BFF e exige HTTP 200.
func authedGET(ctx context.Context, path string) ([]byte, error) {
	resp, err := authedAPI(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, &ctbz.HTTPError{Step: path, Status: resp.Status, Body: resp.Body}
	}
	return resp.Body, nil
}

// buildHistorico monta o registro de saída a partir das respostas de
// payments/recorrencia/init e payments/recorrencia/historico. Só campos
// conhecidos são copiados: a resposta pode trazer dados de cartão ou chaves
// do gateway de pagamento (adyenClientKey), que nunca devem chegar à saída.
func buildHistorico(initBody, histBody []byte) (*output.Record, error) {
	var ini map[string]any
	if err := decodeJSON(initBody, &ini); err != nil {
		return nil, fmt.Errorf("resposta inesperada de payments/recorrencia/init: %w", err)
	}
	var hist struct {
		Pagamentos []map[string]any `json:"pagamentos"`
	}
	if err := decodeJSON(histBody, &hist); err != nil {
		return nil, fmt.Errorf("resposta inesperada de payments/recorrencia/historico: %w", err)
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
		fmt.Fprintln(os.Stderr, "aviso: formato de pagamentos não reconhecido; veja o original com `ctbz api payments/recorrencia/historico`")
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
