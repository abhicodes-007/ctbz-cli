package output

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "regrava os arquivos golden em testdata/")

func init() { time.Local = time.FixedZone("BRT", -3*3600) }

func sampleList() *List {
	l := &List{Columns: []Column{
		{Key: "imposto", Header: "Imposto"},
		{Key: "cnpj", Header: "CNPJ"},
		{Key: "vencimento", Header: "Vencimento"},
		{Key: "valor", Header: "Valor"},
		{Key: "pago", Header: "Pago"},
		{Key: "obs", Header: "Observação"},
	}}
	l.Append("DAS", NewCNPJ("22.222.222/0001-22"), mustDate("20/10/2026"), Money(1234.5), false, "inclui juros, multa")
	l.Append("INSS – pró-labore", NewCNPJ("22222222000122"), mustDate("2026-11-20"), Money(-0.004), true, nil)
	l.Append("IRRF", CNPJ(""), Date{}, Money(1234567.891), nil, `aspas "duplas"`)
	return l
}

func sampleRecord() *Record {
	r := &Record{}
	r.Add("razao_social", "Razão social", "FULANO TECNOLOGIA LTDA").
		Add("cnpj", "CNPJ", NewCNPJ("22222222000122")).
		Add("abertura", "Abertura", DateFromMillis(1784581200000)).
		Add("login_em", "Login em", DateTime{time.Date(2026, 10, 3, 0, 11, 55, 0, time.UTC)}).
		Add("ramos_atividade", "Ramos de atividade", []string{"SERVICO", "COMERCIO"}).
		Add("mensalidade", "Mensalidade", Money(139)).
		Add("inscricao_estadual", "Inscrição estadual", nil).
		Add("outras_empresas", "Outras empresas", []Record{
			{Fields: []Field{{"cnpj", "CNPJ", NewCNPJ("11111111000111")}, {"razao_social", "Razão social", "FULANO MEI"}, {"situacao", "Situação", "INATIVO"}}},
		})
	return r
}

func mustJSON(s string) Data {
	d, err := FromJSON([]byte(s))
	if err != nil {
		panic(err)
	}
	return d
}

func mustDate(s string) Date {
	d, ok := ParseDate(s)
	if !ok {
		panic(s)
	}
	return d
}

func golden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (rode `go test ./internal/output -update`)", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s diferente do golden:\n--- obtido ---\n%s\n--- esperado ---\n%s", name, got, want)
	}
}

func TestGolden(t *testing.T) {
	jsonData, err := FromJSON([]byte(`[{"id":4801207998644224,"nome":{"label":"DAS"},"valor":12.3,"ativo":true},{"id":2,"extra":null,"nome":{"label":"<IRRF> & cia"}}]`))
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]Data{
		"list":     sampleList(),
		"record":   sampleRecord(),
		"fromjson": jsonData,
		"nested":   mustJSON(`{"emAtraso":[],"esteMes":[{"id":1,"valor":{"label":10.5}},{"id":2,"vencimento":"20/10/2026"}],"tags":["a","b"]}`),
	}
	for name, d := range cases {
		for _, f := range Formats {
			var buf bytes.Buffer
			if err := Write(&buf, f, d); err != nil {
				t.Fatalf("%s/%s: %v", name, f, err)
			}
			golden(t, name+"."+string(f), buf.Bytes())
		}
	}
}

func TestFormatBRL(t *testing.T) {
	for v, want := range map[float64]string{
		0:          "R$ 0,00",
		0.005:      "R$ 0,01",
		-0.004:     "R$ 0,00",
		12.3:       "R$ 12,30",
		1234.56:    "R$ 1.234,56",
		-1234.56:   "-R$ 1.234,56",
		1234567.89: "R$ 1.234.567,89",
	} {
		if got := FormatBRL(v); got != want {
			t.Errorf("FormatBRL(%v) = %q, quero %q", v, got, want)
		}
	}
}

func TestParseBRL(t *testing.T) {
	for in, want := range map[string]Money{"R$ 1.621,00": 1621, "R$ 932,31": 932.31, "-R$ 0,50": -0.5, "1.234.567,89": 1234567.89, " R$ 0,00 ": 0} {
		if got, ok := ParseBRL(in); !ok || got != want {
			t.Errorf("ParseBRL(%q) = %v, %v; quero %v", in, got, ok, want)
		}
	}
	for _, in := range []string{"", "R$", "abc", "R$ x,00"} {
		if _, ok := ParseBRL(in); ok {
			t.Errorf("ParseBRL(%q) deveria falhar", in)
		}
	}
	// ida e volta com FormatBRL
	if got, _ := ParseBRL(FormatBRL(-1234.5)); got != -1234.5 {
		t.Errorf("ida e volta: %v", got)
	}
}

func TestCPF(t *testing.T) {
	c := NewCPF("123.456.789-01")
	if c != "12345678901" || FormatCPF(c) != "123.456.789-01" || cellText(c, FormatCSV) != "12345678901" {
		t.Errorf("CPF: %q %q", c, FormatCPF(c))
	}
	if raw, _ := marshal(CPF("")); string(raw) != "null" {
		t.Errorf("CPF vazio em JSON = %s", raw)
	}
}

func TestText(t *testing.T) {
	longo := Text(strings.Repeat("á", 70) + " <fim>")
	if got := cellText(longo, FormatTable); len([]rune(got)) != textMax || !strings.HasSuffix(got, "…") {
		t.Errorf("tabela: %q", got)
	}
	if got := cellText(longo, FormatCSV); got != string(longo) {
		t.Errorf("CSV deveria trazer o texto completo: %q", got)
	}
	if raw, _ := marshal(longo); !strings.HasSuffix(string(raw), ` <fim>"`) {
		t.Errorf("JSON deveria trazer o texto completo, sem escapar <>: %s", raw)
	}
	if raw, _ := marshal(Text("")); string(raw) != "null" {
		t.Errorf("Text vazio em JSON = %s", raw)
	}
}

func TestParseFormat(t *testing.T) {
	for _, s := range []string{"table", "JSON", " csv "} {
		if _, err := ParseFormat(s); err != nil {
			t.Errorf("ParseFormat(%q): %v", s, err)
		}
	}
	if _, err := ParseFormat("yaml"); err == nil || !strings.Contains(err.Error(), "table, json ou csv") {
		t.Errorf("esperava erro para yaml, veio %v", err)
	}
}

func TestEmptyList(t *testing.T) {
	l := &List{Columns: []Column{{Key: "a", Header: "A"}}}
	for f, want := range map[Format]string{FormatTable: "", FormatJSON: "[]\n", FormatCSV: "a\n"} {
		var buf bytes.Buffer
		if err := Write(&buf, f, l); err != nil {
			t.Fatal(err)
		}
		if buf.String() != want {
			t.Errorf("%s: %q, quero %q", f, buf.String(), want)
		}
	}
}

func TestFromJSONScalar(t *testing.T) {
	d, err := FromJSON([]byte(`"OK"`))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	Write(&buf, FormatTable, d)
	if buf.String() != "OK\n" {
		t.Errorf("tabela = %q", buf.String())
	}
	if _, err := FromJSON([]byte(`{} {}`)); err == nil {
		t.Error("esperava erro para JSON com conteúdo extra")
	}
}
