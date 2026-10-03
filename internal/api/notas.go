package api

import (
	"context"
	"fmt"
	"net/url"
)

// PathNotasEmitidas é a listagem v2 do emissor de NFS-e (a que o painel usa quando
// novo-emissor/v2/feature-flag é verdadeiro).
const PathNotasEmitidas = "novo-emissor/v2/listagem/notas/filtro"

// Campos de filtro aceitos pela listagem (um por vez).
const (
	FiltroNomeTomador = "nomeTomador"
	FiltroDocumento   = "documento"
	FiltroNumeroNota  = "numeroNota"
)

// NotaEmitida é uma NFS-e emitida. Os campos são os lidos pelo front; a conta verificada
// não tinha notas, por isso todos são opcionais e os de formato incerto ficam como any.
type NotaEmitida struct {
	ID                   any      `json:"id" contract:"optional"`
	Numero               any      `json:"numero" contract:"optional"`
	CPFCNPJTomador       string   `json:"cpfCnpjTomador" contract:"optional"`
	NomeRazaoTomador     string   `json:"nomeRazaoTomador" contract:"optional"`
	ValorServico         *float64 `json:"valorServico" contract:"optional"`
	DataEmissao          any      `json:"dataEmissao" contract:"optional"`
	DataEmissaoFormatada string   `json:"dataEmissaoFormatada" contract:"optional"`
	StatusNotaFiscal     string   `json:"statusNotaFiscal" contract:"optional"`
	SituacaoNota         string   `json:"situacaoNota" contract:"optional"`
	TomadorExterior      bool     `json:"tomadorExterior" contract:"optional"`
}

// ListaNotas é uma página da listagem de notas.
type ListaNotas struct {
	List  []NotaEmitida `json:"list"`
	Total int           `json:"total"`
}

// FiltroNotas seleciona as notas de um mês (1–12), opcionalmente filtradas por um campo.
type FiltroNotas struct {
	Ano, Mes     int
	Campo, Termo string // Campo é um dos Filtro*; vazio, sem filtro
}

// tamanhoPaginaNotas é o tamanho de página padrão do painel.
const tamanhoPaginaNotas = 10

// BuscarNotasEmitidas lê todas as páginas de notas emitidas de um mês.
func BuscarNotasEmitidas(ctx context.Context, g Getter, f FiltroNotas) ([]NotaEmitida, error) {
	var out []NotaEmitida
	for pagina := 1; pagina <= maxPaginas; pagina++ {
		path := fmt.Sprintf("%s?pagina=%d&limite=%d&ano=%d&mes=%d", PathNotasEmitidas, pagina, tamanhoPaginaNotas, f.Ano, f.Mes)
		if f.Campo != "" {
			path += "&" + f.Campo + "=" + url.QueryEscape(f.Termo)
		}
		l, err := get[ListaNotas](ctx, g, path)
		if err != nil {
			return nil, err
		}
		out = append(out, l.List...)
		if len(l.List) == 0 || len(out) >= l.Total {
			break
		}
	}
	return out, nil
}
