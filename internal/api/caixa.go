package api

import (
	"context"
	"fmt"
)

// PathCaixa é a lista de lançamentos do caixa de um mês (1–12). O painel pede 1000 registros
// e manda a página literalmente como "null"; a CLI faz o mesmo.
func PathCaixa(ano, mes int) string {
	return fmt.Sprintf("caixa/listpaginada/%d/%d/%d/null", ano, mes, TamanhoPaginaCaixa)
}

// TamanhoPaginaCaixa é o número de registros pedido por mês, como no painel.
const TamanhoPaginaCaixa = 1000

const PathContasUsuario = "movimentacao-financeira/contasUsuario"

// LancamentoCaixa é um lançamento do caixa (entrada ou saída classificada).
type LancamentoCaixa struct {
	ID                   int64    `json:"id"`
	Data                 int64    `json:"data"` // epoch em ms
	Descricao            string   `json:"descricao"`
	Valor                *float64 `json:"valor"` // negativo nas saídas
	Situacao             string   `json:"situacao"`
	Tipo                 string   `json:"tipo"`
	IDContaUsuario       *int64   `json:"idContaUsuario"` // classificação (ContaUsuario)
	ConfirmadoViaSistema bool     `json:"confirmadoViaSistema"`
}

// Caixa é a resposta de caixa/listpaginada.
type Caixa struct {
	List  []LancamentoCaixa `json:"list"`
	Total int               `json:"total"`
}

// BuscarCaixa lê os lançamentos do caixa de um mês.
func BuscarCaixa(ctx context.Context, g Getter, ano, mes int) (*Caixa, error) {
	return get[Caixa](ctx, g, PathCaixa(ano, mes))
}

// ContaUsuario é uma conta do plano de contas simplificado usado para classificar os
// lançamentos (ex.: "Pagamento de Fornecedores"), ligada a uma conta contábil.
type ContaUsuario struct {
	ID                     int64  `json:"id"`
	Descricao              string `json:"descricao"`
	DescricaoContaContabil string `json:"descricaoContaContabil"`
	Classificacao          string `json:"classificacao"` // ex.: RECEITA, DESPESA
	Situacao               string `json:"situacao"`      // ATIVO ou INATIVO
}

// BuscarContasUsuario lê o plano de contas usado nas classificações.
func BuscarContasUsuario(ctx context.Context, g Getter) ([]ContaUsuario, error) {
	c, err := get[[]ContaUsuario](ctx, g, PathContasUsuario)
	if err != nil {
		return nil, err
	}
	return *c, nil
}
