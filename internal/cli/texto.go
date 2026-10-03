package cli

import (
	"html"
	"regexp"
	"strings"

	"github.com/edusouza/ctbz-cli/internal/output"
)

var (
	reTag     = regexp.MustCompile(`<[^>]*>`)
	reEspacos = regexp.MustCompile(`\s+`)
)

// textoSimples tira as tags de um trecho de HTML vindo da API e junta os espaços.
func textoSimples(s string) output.Text {
	s = html.UnescapeString(reTag.ReplaceAllString(s, " "))
	return output.Text(strings.TrimSpace(reEspacos.ReplaceAllString(s, " ")))
}
