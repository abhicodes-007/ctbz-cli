package output

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Money é um valor em reais. Tabela: "R$ 1.234,56"; JSON: 1234.56; CSV: "1234.56".
type Money float64

// Date é uma data sem horário. Tabela: "02/01/2006"; JSON/CSV: "2006-01-02".
// A data zero é tratada como ausente.
type Date struct{ time.Time }

// DateTime é um instante. Tabela: "02/01/2006 15:04" (horário local);
// JSON/CSV: RFC 3339.
type DateTime struct{ time.Time }

// CNPJ guarda só os dígitos. Tabela: "00.000.000/0000-00"; JSON/CSV: dígitos.
type CNPJ string

// CPF guarda só os dígitos. Tabela: "000.000.000-00"; JSON/CSV: dígitos.
type CPF string

// Text é um texto longo (ex.: descrições). Tabela: cortado em 60 caracteres com "…";
// JSON/CSV: completo. Texto vazio é tratado como ausente.
type Text string

// Indent é um texto em uma árvore (ex.: plano de contas). Tabela: recuado em dois espaços
// por nível abaixo do primeiro; JSON/CSV: só o texto (o nível vai em outra coluna).
type Indent struct {
	Level int
	Text  string
}

// NewCPF remove a pontuação de um CPF.
func NewCPF(s string) CPF { return CPF(NewCNPJ(s)) }

// FormatCPF aplica a máscara 000.000.000-00 quando há 11 dígitos.
func FormatCPF(c CPF) string {
	s := string(c)
	if len(s) != 11 {
		return s
	}
	return s[0:3] + "." + s[3:6] + "." + s[6:9] + "-" + s[9:11]
}

// NewCNPJ remove a pontuação de um CNPJ.
func NewCNPJ(s string) CNPJ {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return CNPJ(b.String())
}

// ParseDate aceita "dd/mm/aaaa" e "aaaa-mm-dd" (também com horário ISO).
func ParseDate(s string) (Date, bool) {
	s = strings.TrimSpace(s)
	for _, layout := range []string{"02/01/2006", "2006-01-02", time.RFC3339, "2006-01-02T15:04:05"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return Date{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)}, true
		}
	}
	return Date{}, false
}

// DateFromMillis converte epoch em milissegundos (formato comum nas APIs).
func DateFromMillis(ms int64) Date {
	t := time.UnixMilli(ms).In(time.Local)
	return Date{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)}
}

// FormatBRL formata um valor como "R$ 1.234,56" (negativos: "-R$ 1.234,56").
func FormatBRL(v float64) string {
	cents := int64(math.Round(math.Abs(v) * 100))
	inteiro := strconv.FormatInt(cents/100, 10)
	var b strings.Builder
	for i, r := range inteiro {
		if i > 0 && (len(inteiro)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(r)
	}
	sinal := ""
	if v < 0 && cents != 0 {
		sinal = "-"
	}
	return fmt.Sprintf("%sR$ %s,%02d", sinal, b.String(), cents%100)
}

// ParseBRL lê um valor formatado em reais ("R$ 1.234,56", "-R$ 0,50", "1.234,56").
// Texto que não é um valor devolve false.
func ParseBRL(s string) (Money, bool) {
	s = strings.TrimSpace(s)
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(s, "-"), "R$"))
	s = strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", ".")
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil || s == "" {
		return 0, false
	}
	if neg {
		v = -v
	}
	return Money(v), true
}

// FormatCNPJ aplica a máscara 00.000.000/0000-00 quando há 14 dígitos.
func FormatCNPJ(c CNPJ) string {
	s := string(c)
	if len(s) != 14 {
		return s
	}
	return s[0:2] + "." + s[2:5] + "." + s[5:8] + "/" + s[8:12] + "-" + s[12:14]
}

// cellText converte um valor em texto para tabela ou CSV.
func cellText(v any, f Format) string {
	table := f == FormatTable
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case bool:
		if table {
			if x {
				return "sim"
			}
			return "não"
		}
		return strconv.FormatBool(x)
	case Money:
		if table {
			return FormatBRL(float64(x))
		}
		return moneyRaw(x)
	case Date:
		if x.IsZero() {
			return ""
		}
		if table {
			return x.Format("02/01/2006")
		}
		return x.Format("2006-01-02")
	case DateTime:
		if x.IsZero() {
			return ""
		}
		if table {
			return x.Local().Format("02/01/2006 15:04")
		}
		return x.Format(time.RFC3339)
	case Text:
		if table {
			return truncate(string(x), textMax)
		}
		return string(x)
	case Indent:
		if table && x.Level > 1 {
			return strings.Repeat("  ", x.Level-1) + x.Text
		}
		return x.Text
	case CNPJ:
		if table {
			return FormatCNPJ(x)
		}
		return string(x)
	case CPF:
		if table {
			return FormatCPF(x)
		}
		return string(x)
	case []string:
		if table {
			return strings.Join(x, ", ")
		}
		return strings.Join(x, "; ")
	case json.Number:
		return x.String()
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprint(x)
	case fmt.Stringer:
		return x.String()
	default:
		// Estruturas aninhadas viram JSON compacto.
		raw, err := marshal(v)
		if err != nil {
			return fmt.Sprint(v)
		}
		s := string(raw)
		if table {
			s = truncate(s, textMax)
		}
		return s
	}
}

// moneyRaw arredonda para centavos sem produzir "-0.00".
func moneyRaw(m Money) string {
	v := math.Round(float64(m)*100) / 100
	if v == 0 {
		v = 0 // descarta o zero negativo
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}

// numeric indica valores que ficam alinhados à direita na tabela.
func numeric(v any) bool {
	switch v.(type) {
	case Money, json.Number, float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

// textMax é o tamanho máximo de Text e de JSON aninhado numa célula da tabela.
const textMax = 60

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}
