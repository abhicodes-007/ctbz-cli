package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// FromJSON converte uma resposta JSON arbitrária em Data, preservando a ordem
// dos campos e a precisão dos números (IDs de 16 dígitos não cabem em float64):
//
//   - objeto            → Record
//   - lista de objetos  → List (colunas = união das chaves, na ordem em que aparecem)
//   - lista de escalares → List com a coluna "valor"
//   - escalar           → Scalar
func FromJSON(data []byte) (Data, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	v, err := decodeValue(dec)
	if err != nil {
		return nil, err
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, errors.New("JSON com conteúdo extra após o valor principal")
	}
	switch x := v.(type) {
	case *Record:
		return x, nil
	case []any:
		return listFromArray(x), nil
	default:
		return Scalar{Value: v}, nil
	}
}

func listFromArray(items []any) *List {
	l := &List{}
	allRecords := len(items) > 0
	for _, it := range items {
		if _, ok := it.(*Record); !ok {
			allRecords = false
			break
		}
	}
	if !allRecords {
		l.Columns = []Column{{Key: "valor", Header: "valor"}}
		for _, it := range items {
			l.Rows = append(l.Rows, []any{it})
		}
		return l
	}
	index := map[string]int{}
	for _, it := range items {
		for _, f := range it.(*Record).Fields {
			if _, ok := index[f.Key]; !ok {
				index[f.Key] = len(l.Columns)
				l.Columns = append(l.Columns, Column{Key: f.Key, Header: f.Key})
			}
		}
	}
	for _, it := range items {
		row := make([]any, len(l.Columns))
		for _, f := range it.(*Record).Fields {
			row[index[f.Key]] = f.Value
		}
		l.Rows = append(l.Rows, row)
	}
	return l
}

func decodeValue(dec *json.Decoder) (any, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			r := &Record{}
			for dec.More() {
				kt, err := dec.Token()
				if err != nil {
					return nil, err
				}
				key, ok := kt.(string)
				if !ok {
					return nil, fmt.Errorf("chave JSON inválida: %v", kt)
				}
				v, err := decodeValue(dec)
				if err != nil {
					return nil, err
				}
				r.Fields = append(r.Fields, Field{Key: key, Label: key, Value: v})
			}
			if _, err := dec.Token(); err != nil { // '}'
				return nil, err
			}
			return r, nil
		case '[':
			arr := []any{}
			for dec.More() {
				v, err := decodeValue(dec)
				if err != nil {
					return nil, err
				}
				arr = append(arr, v)
			}
			if _, err := dec.Token(); err != nil { // ']'
				return nil, err
			}
			return arr, nil
		}
		return nil, fmt.Errorf("delimitador JSON inesperado: %v", t)
	default:
		return tok, nil // string, json.Number, bool ou nil
	}
}
