package api

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

const (
	PathConciliacaoInit       = "conciliacao-fiscal/v2/init"
	PathConciliacaoPendencias = "conciliacao-fiscal/v2/pendencias"
)

// Tipos de pendência de conciliação (parâmetro tipoPendencia).
const (
	// PendenciaNotaSemRecebimento é uma nota fiscal emitida sem recebimento vinculado.
	PendenciaNotaSemRecebimento = "RECEITA_SEM_RECEBIMENTO"
	// PendenciaRecebimentoSemNota é um recebimento no extrato sem nota fiscal vinculada.
	PendenciaRecebimentoSemNota = "RECEBIMENTO_SEM_RECEITA"
)

// ConciliacaoResumo é a resposta de conciliacao-fiscal/v2/init.
type ConciliacaoResumo struct {
	QtdNotasFiscaisPendentes              int    `json:"qtdNotasFiscaisPendentes"`
	QtdRecebimentosPendentes              int    `json:"qtdRecebimentosPendentes"`
	QtdConciliacoesAutomaticasMesAnterior int    `json:"qtdConciliacoesAutomaticasMesAnterior"`
	CompetenciaMesAnterior                string `json:"competenciaMesAnterior"` // AAAA-MM
}

// BuscarConciliacaoResumo lê as quantidades de pendências de conciliação.
func BuscarConciliacaoResumo(ctx context.Context, g Getter) (*ConciliacaoResumo, error) {
	return get[ConciliacaoResumo](ctx, g, PathConciliacaoInit)
}

// PaginaConciliacao é uma página de conciliacao-fiscal/v2/pendencias. O formato dos itens
// ainda não foi verificado (a conta de captura não tinha pendências).
type PaginaConciliacao struct {
	Pagina         []json.RawMessage `json:"pagina"`
	TotalRegistros int               `json:"totalRegistros"`
	TotalPaginas   int               `json:"totalPaginas"`
}

// BuscarPendenciasConciliacao lê todas as páginas de pendências abertas de um tipo
// (PendenciaNotaSemRecebimento ou PendenciaRecebimentoSemNota) com data entre de e ate.
func BuscarPendenciasConciliacao(ctx context.Context, g Getter, tipo string, de, ate time.Time) ([]json.RawMessage, error) {
	var out []json.RawMessage
	for pagina := 1; ; pagina++ {
		q := url.Values{
			"pagina":         {strconv.Itoa(pagina)},
			"totalPagina":    {"10"}, // o mesmo tamanho de página do painel
			"status":         {"PENDENTE"},
			"tipoPendencia":  {tipo},
			"periodoInicial": {de.Format("2006-01-02")},
			"periodoFinal":   {ate.Format("2006-01-02")},
		}
		p, err := get[PaginaConciliacao](ctx, g, PathConciliacaoPendencias+"?"+q.Encode())
		if err != nil {
			return nil, err
		}
		out = append(out, p.Pagina...)
		if pagina >= p.TotalPaginas || pagina >= maxPaginas {
			return out, nil
		}
	}
}
