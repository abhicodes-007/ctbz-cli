// Package output renderiza os dados dos comandos em três formatos:
//
//   - table: para pessoas (colunas alinhadas, R$ 1.234,56, dd/mm/aaaa);
//   - json:  para scripts (campos na ordem declarada, valores crus, datas ISO);
//   - csv:   para planilhas (cabeçalho com as chaves, valores crus).
//
// Os comandos descrevem os dados com List (várias linhas) ou Record (um
// registro) e usam os tipos Money, Date, DateTime e CNPJ para que a
// formatação seja a mesma em toda a CLI.
package output

import (
	"fmt"
	"io"
	"strings"
)

// Format é um formato de saída.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
)

// Formats lista os formatos aceitos, na ordem exibida na ajuda.
var Formats = []Format{FormatTable, FormatJSON, FormatCSV}

// ParseFormat valida o nome de um formato (sem diferenciar maiúsculas).
func ParseFormat(s string) (Format, error) {
	f := Format(strings.ToLower(strings.TrimSpace(s)))
	for _, ok := range Formats {
		if f == ok {
			return f, nil
		}
	}
	return "", fmt.Errorf("formato de saída inválido %q: use table, json ou csv", s)
}

// Data é qualquer coisa que os comandos podem imprimir.
type Data interface{ isData() }

// Column descreve uma coluna de uma List. Key é o nome do campo em JSON e no
// cabeçalho CSV; Header é o título mostrado na tabela.
type Column struct {
	Key    string
	Header string
}

// List é uma sequência de linhas com as mesmas colunas.
type List struct {
	Columns []Column
	Rows    [][]any
}

// Field é um campo de um Record.
type Field struct {
	Key   string
	Label string
	Value any
}

// Record é um único registro com campos ordenados.
type Record struct {
	Fields []Field
}

// Scalar é um valor solto (ex.: resposta de texto de uma API).
type Scalar struct{ Value any }

func (List) isData()   {}
func (Record) isData() {}
func (Scalar) isData() {}

// Add acrescenta um campo ao registro e o devolve, para encadear.
func (r *Record) Add(key, label string, value any) *Record {
	r.Fields = append(r.Fields, Field{Key: key, Label: label, Value: value})
	return r
}

// Append acrescenta uma linha; o número de valores deve bater com as colunas.
func (l *List) Append(values ...any) {
	if len(values) != len(l.Columns) {
		panic(fmt.Sprintf("output: linha com %d valores para %d colunas", len(values), len(l.Columns)))
	}
	l.Rows = append(l.Rows, values)
}

// Write renderiza os dados no formato pedido.
func Write(w io.Writer, f Format, d Data) error {
	switch f {
	case FormatJSON:
		return writeJSON(w, d)
	case FormatCSV:
		return writeCSV(w, d)
	case FormatTable, "":
		return writeTable(w, d)
	default:
		return fmt.Errorf("formato de saída desconhecido: %q", f)
	}
}

// RecordsToList converte registros com os mesmos campos numa List, usando os
// campos do primeiro registro como colunas.
func RecordsToList(recs []Record) *List {
	l := &List{}
	if len(recs) == 0 {
		return l
	}
	for _, f := range recs[0].Fields {
		l.Columns = append(l.Columns, Column{Key: f.Key, Header: f.Label})
	}
	for _, r := range recs {
		byKey := map[string]any{}
		for _, f := range r.Fields {
			byKey[f.Key] = f.Value
		}
		row := make([]any, len(l.Columns))
		for i, c := range l.Columns {
			row[i] = byKey[c.Key]
		}
		l.Rows = append(l.Rows, row)
	}
	return l
}
