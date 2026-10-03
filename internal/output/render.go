package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

// ---------- JSON

func writeJSON(w io.Writer, d Data) error {
	raw, err := marshal(d)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return err
	}
	buf.WriteByte('\n')
	_, err = w.Write(buf.Bytes())
	return err
}

// marshal gera JSON compacto preservando a ordem dos campos.
func marshal(v any) ([]byte, error) {
	switch x := v.(type) {
	case nil:
		return []byte("null"), nil
	case *Record:
		return marshalRecord(*x)
	case Record:
		return marshalRecord(x)
	case []Record:
		parts := make([][]byte, len(x))
		for i, r := range x {
			b, err := marshalRecord(r)
			if err != nil {
				return nil, err
			}
			parts[i] = b
		}
		return joinArray(parts), nil
	case *List:
		return marshalList(*x)
	case List:
		return marshalList(x)
	case Scalar:
		return marshal(x.Value)
	case []any:
		parts := make([][]byte, len(x))
		for i, e := range x {
			b, err := marshal(e)
			if err != nil {
				return nil, err
			}
			parts[i] = b
		}
		return joinArray(parts), nil
	case Money:
		return []byte(moneyRaw(x)), nil
	case Date:
		if x.IsZero() {
			return []byte("null"), nil
		}
		return json.Marshal(x.Format("2006-01-02"))
	case DateTime:
		if x.IsZero() {
			return []byte("null"), nil
		}
		return json.Marshal(x.Format(time.RFC3339))
	case CNPJ:
		if x == "" {
			return []byte("null"), nil
		}
		return json.Marshal(string(x))
	case CPF:
		if x == "" {
			return []byte("null"), nil
		}
		return json.Marshal(string(x))
	case string:
		return marshalString(x)
	default:
		return json.Marshal(v)
	}
}

// marshalString não escapa <, > e & (o padrão de encoding/json), para que
// textos fiquem legíveis.
func marshalString(s string) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(s); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func marshalRecord(r Record) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, f := range r.Fields {
		if i > 0 {
			buf.WriteByte(',')
		}
		k, err := marshalString(f.Key)
		if err != nil {
			return nil, err
		}
		v, err := marshal(f.Value)
		if err != nil {
			return nil, fmt.Errorf("campo %s: %w", f.Key, err)
		}
		buf.Write(k)
		buf.WriteByte(':')
		buf.Write(v)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func marshalList(l List) ([]byte, error) {
	parts := make([][]byte, len(l.Rows))
	for i, row := range l.Rows {
		r := Record{}
		for j, c := range l.Columns {
			r.Fields = append(r.Fields, Field{Key: c.Key, Value: row[j]})
		}
		b, err := marshalRecord(r)
		if err != nil {
			return nil, err
		}
		parts[i] = b
	}
	return joinArray(parts), nil
}

func joinArray(parts [][]byte) []byte {
	return append(append([]byte{'['}, bytes.Join(parts, []byte{','})...), ']')
}

// ---------- CSV

func writeCSV(w io.Writer, d Data) error {
	cw := csv.NewWriter(w)
	switch x := d.(type) {
	case *List:
		writeCSVList(cw, *x)
	case List:
		writeCSVList(cw, x)
	case *Record:
		writeCSVRecord(cw, *x)
	case Record:
		writeCSVRecord(cw, x)
	case Scalar:
		cw.Write([]string{cellText(x.Value, FormatCSV)})
	default:
		return fmt.Errorf("output: tipo não suportado em CSV: %T", d)
	}
	cw.Flush()
	return cw.Error()
}

func writeCSVList(cw *csv.Writer, l List) {
	header := make([]string, len(l.Columns))
	for i, c := range l.Columns {
		header[i] = c.Key
	}
	cw.Write(header)
	for _, row := range l.Rows {
		cells := make([]string, len(row))
		for i, v := range row {
			cells[i] = cellText(v, FormatCSV)
		}
		cw.Write(cells)
	}
}

// Um registro vira uma linha com cabeçalho, como uma lista de um item.
func writeCSVRecord(cw *csv.Writer, r Record) {
	header := make([]string, len(r.Fields))
	cells := make([]string, len(r.Fields))
	for i, f := range r.Fields {
		header[i] = f.Key
		cells[i] = cellText(f.Value, FormatCSV)
	}
	cw.Write(header)
	cw.Write(cells)
}

// ---------- Tabela

func writeTable(w io.Writer, d Data) error {
	var buf bytes.Buffer
	switch x := d.(type) {
	case *List:
		tableList(&buf, *x, "")
	case List:
		tableList(&buf, x, "")
	case *Record:
		tableRecord(&buf, *x)
	case Record:
		tableRecord(&buf, x)
	case Scalar:
		buf.WriteString(cellText(x.Value, FormatTable) + "\n")
	default:
		return fmt.Errorf("output: tipo não suportado em tabela: %T", d)
	}
	_, err := w.Write(buf.Bytes())
	return err
}

func tableList(buf *bytes.Buffer, l List, indent string) {
	if len(l.Rows) == 0 {
		return
	}
	n := len(l.Columns)
	cells := make([][]string, 0, len(l.Rows)+1)
	header := make([]string, n)
	for i, c := range l.Columns {
		header[i] = c.Header
		if header[i] == "" {
			header[i] = c.Key
		}
	}
	cells = append(cells, header)
	right := make([]bool, n)
	for i := range right {
		right[i] = true
	}
	for _, row := range l.Rows {
		line := make([]string, n)
		for i, v := range row {
			line[i] = strings.ReplaceAll(cellText(v, FormatTable), "\n", " ")
			if v != nil && !numeric(v) {
				right[i] = false
			}
		}
		cells = append(cells, line)
	}
	writeAligned(buf, cells, right, indent)
}

func tableRecord(buf *bytes.Buffer, r Record) {
	var simple [][]string
	flush := func() {
		writeAligned(buf, simple, []bool{false, false}, "")
		simple = nil
	}
	for _, f := range r.Fields {
		var nested *List
		switch x := f.Value.(type) {
		case []Record:
			nested = RecordsToList(x)
		case *List:
			nested = x
		case List:
			nested = &x
		case []any: // listas vindas de FromJSON
			if l := listFromArray(x); len(x) == 0 || l.Columns[0].Key != "valor" {
				nested = l
			}
		}
		if nested != nil {
			if len(nested.Rows) == 0 {
				continue
			}
			flush()
			buf.WriteString(f.Label + ":\n")
			tableList(buf, *nested, "  ")
			continue
		}
		text := cellText(f.Value, FormatTable)
		if text == "" {
			continue // campos vazios poluem a visão em tabela
		}
		simple = append(simple, []string{f.Label + ":", text})
	}
	flush()
}

// writeAligned alinha colunas contando runas (acentos ocupam uma posição).
func writeAligned(buf *bytes.Buffer, rows [][]string, right []bool, indent string) {
	if len(rows) == 0 {
		return
	}
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, c := range row {
			if n := utf8.RuneCountInString(c); n > widths[i] {
				widths[i] = n
			}
		}
	}
	for _, row := range rows {
		var line strings.Builder
		line.WriteString(indent)
		for i, c := range row {
			if i > 0 {
				line.WriteString("  ")
			}
			pad := strings.Repeat(" ", widths[i]-utf8.RuneCountInString(c))
			if right[i] {
				line.WriteString(pad + c)
			} else {
				line.WriteString(c + pad)
			}
		}
		buf.WriteString(strings.TrimRight(line.String(), " ") + "\n")
	}
}
