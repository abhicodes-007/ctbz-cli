package contract

import (
	"bytes"
	"encoding/json"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Anonymize troca dados pessoais e financeiros de uma resposta JSON por valores
// fictícios, mantendo a estrutura e os tipos. É usada para gerar as fixtures de
// teste a partir de respostas reais. O resultado deve ser revisado antes do commit:
// as regras cobrem os campos conhecidos, não qualquer texto livre.
func Anonymize(data []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var doc any
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}
	a := &anonymizer{ids: map[string]string{}}
	doc = a.value("", doc, false)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type anonymizer struct {
	ids  map[string]string // ID real → ID fictício (mesmo ID, mesmo substituto)
	next int64
}

var (
	reEmail      = regexp.MustCompile(`[\w.+-]+@[\w-]+(\.[\w-]+)+`)
	reCNPJMasked = regexp.MustCompile(`\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}`)
	reCPFMasked  = regexp.MustCompile(`\d{3}\.\d{3}\.\d{3}-\d{2}`)
	reLongDigits = regexp.MustCompile(`\d{8,}`)
	reURL        = regexp.MustCompile(`https?://\S+`)

	// Chaves cujo valor identifica pessoas, empresas ou endereços.
	keyName    = regexp.MustCompile(`(?i)(nome|razao|fantasia|titular|responsavel|socio|cliente|tomador)`)
	keyAddress = regexp.MustCompile(`(?i)(logradouro|endereco|bairro|complemento|rua|cidade|municipio$)`)
	// Chaves que casam com as regras acima mas não identificam ninguém.
	keyNotPersonal = regexp.MustCompile(`(?i)^nomeMes`)
	keyFree        = regexp.MustCompile(`(?i)^(subject|assunto|mensagem|observacao|detalhe|comentario|motivo)$`)
	keySecret      = regexp.MustCompile(`(?i)^(ref|token|hash|senha|chave|key|clientKey|adyenClientKey|jiraIssueId)$`)
	keyDigits      = regexp.MustCompile(`(?i)(cpf|cnpj|cep|telefone|celular|conta$|numeroConta|agencia|inscricao|pis|nit|documento|identificador|^numero$|rg$|cnh|eleitor|rne|passaporte)`)
	keyBirth       = regexp.MustCompile(`(?i)nascimento`)
	keyMoney       = regexp.MustCompile(`(?i)(valor|saldo|total|faturamento|receita|credito|debito|preco|montante|juros|multa|prolabore|lucro|distribu|adiantamento|base|imposto|economia|cenario|pago|custo)`)
	keyID          = regexp.MustCompile(`(?i)^id|id$|^(mes|ano|periodo|dia|quantidade\w*|nr\w*|numero\w*)$`)
)

const (
	longText = 300 // textos maiores (contratos, HTML) são omitidos
	maxItems = 3   // listas são cortadas neste tamanho
)

// value anonimiza v. money indica que algum objeto acima é monetário
// (ex.: "valor": {"label": 1234.5}): números dentro dele também são trocados.
func (a *anonymizer) value(key string, v any, money bool) any {
	money = money || keyMoney.MatchString(key)
	switch x := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys) // ordem determinística para os IDs fictícios
		for _, k := range keys {
			x[k] = a.value(k, x[k], money)
		}
		return x
	case []any:
		if len(x) > maxItems {
			x = x[:maxItems] // fixtures pequenas e com menos dados reais para revisar
		}
		for i := range x {
			x[i] = a.value(key, x[i], money)
		}
		return x
	case string:
		return a.text(key, x)
	case json.Number:
		return a.number(key, x, money)
	}
	return v
}

func (a *anonymizer) text(key, s string) string {
	switch {
	case s == "":
		return s
	case len(s) > longText:
		return "<texto omitido>"
	case keyNotPersonal.MatchString(key):
		return s
	case keySecret.MatchString(key):
		return "XXXX"
	case strings.Contains(strings.ToLower(key), "email"):
		return "fulano@exemplo.com"
	case keyFree.MatchString(key):
		return "TEXTO EXEMPLO"
	case keyDigits.MatchString(key):
		return zeroDigits(s)
	case keyName.MatchString(key):
		return "FULANO DE TAL"
	case keyAddress.MatchString(key):
		return "LOGRADOURO EXEMPLO"
	}
	// Texto livre: aplica todos os padrões.
	s = reURL.ReplaceAllString(s, "https://exemplo.invalid/arquivo")
	s = reEmail.ReplaceAllString(s, "fulano@exemplo.com")
	s = reCNPJMasked.ReplaceAllString(s, "22.222.222/0001-22")
	s = reCPFMasked.ReplaceAllString(s, "000.000.000-00")
	return reLongDigits.ReplaceAllStringFunc(s, zeroDigits)
}

// zeroDigits troca todos os dígitos por zero, preservando máscara e tamanho.
func zeroDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return '0'
		}
		return r
	}, s)
}

func (a *anonymizer) number(key string, n json.Number, money bool) json.Number {
	s := n.String()
	isInt := !strings.ContainsAny(s, ".eE")
	if keyBirth.MatchString(key) {
		return "0" // data de nascimento identifica a pessoa
	}
	if money && !keyID.MatchString(key) {
		if isInt {
			return "1000"
		}
		return "1234.56"
	}
	if !isInt {
		return n
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil || i < 1e14 {
		return n // contagens, anos, epoch em milissegundos (~1.7e12): não identificam ninguém
	}
	// IDs internos (16 dígitos) viram IDs fictícios estáveis dentro da resposta.
	if id, ok := a.ids[s]; ok {
		return json.Number(id)
	}
	a.next++
	id := strconv.FormatInt(1_000_000_000_000_000+a.next, 10)
	a.ids[s] = id
	return json.Number(id)
}
