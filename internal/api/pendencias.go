package api

import "context"

const PathPendenciasEmpresa = "home/pendencia/pendenciasEmpresa"

// PendenciaEmpresa é uma pendência cadastrada para a empresa (tela inicial do painel).
type PendenciaEmpresa struct {
	ID            int64  `json:"id"`
	DataCriacao   int64  `json:"dataCriacao"` // epoch em ms
	DataLimite    string `json:"dataLimite"`  // dd/mm/aaaa
	Detalhe       string `json:"detalhe"`
	TipoPendencia struct {
		Codigo int    `json:"codigo"`
		Titulo string `json:"titulo"`
	} `json:"tipoPendencia"`
	SituacaoPendencia struct {
		ID        string `json:"id"` // ex.: PENDENTE, FINALIZADO
		Descricao string `json:"descricao"`
	} `json:"situacaoPendencia"`
}

// BuscarPendenciasEmpresa lê as pendências cadastradas.
func BuscarPendenciasEmpresa(ctx context.Context, g Getter) ([]PendenciaEmpresa, error) {
	p, err := get[[]PendenciaEmpresa](ctx, g, PathPendenciasEmpresa)
	if err != nil {
		return nil, err
	}
	return *p, nil
}
