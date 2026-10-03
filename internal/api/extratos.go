package api

import "context"

const (
	PathExtratos        = "movimentacao-financeira/v2/extratos"
	PathContasBancarias = "contabancaria/list"
)

// Extrato é a situação do extrato de uma conta bancária em um mês.
type Extrato struct {
	Ano              int    `json:"ano"`
	Mes              int    `json:"mes"`
	IDContaBancaria  int64  `json:"idContaBancaria"`
	Banco            string `json:"banco"`
	Agencia          string `json:"agencia"`
	NumeroConta      string `json:"numeroConta"`
	Situacao         string `json:"situacao"`         // ex.: ABERTO
	StatusIntegracao any    `json:"statusIntegracao"` // formato não verificado (vinha null)
}

// BuscarExtratos lista os extratos de todas as contas (a API não tem filtro).
func BuscarExtratos(ctx context.Context, g Getter) ([]Extrato, error) {
	e, err := get[[]Extrato](ctx, g, PathExtratos)
	if err != nil {
		return nil, err
	}
	return *e, nil
}

// ContaBancaria é uma conta bancária cadastrada da empresa.
type ContaBancaria struct {
	ID               int64    `json:"id"`
	NomeBanco        string   `json:"nomeBanco"`
	CodigoBanco      string   `json:"codigoBanco"`
	Agencia          string   `json:"agencia"`
	ContaCorrente    string   `json:"contaCorrente"`
	VlrSaldoInicial  *float64 `json:"vlrSaldoInicial"`
	DataSaldoInicial int64    `json:"dataSaldoInicial"` // epoch em ms
	StatusIntegracao string   `json:"statusIntegracao"` // ex.: INTEGRADA
	FluxoIntegracao  string   `json:"fluxoIntegracao"`  // ex.: CONTABILIZEI_BANK
}

// ContasBancarias é a resposta de contabancaria/list (que também traz a lista de bancos).
type ContasBancarias struct {
	ContasBancarias []ContaBancaria `json:"contasBancarias"`
}

// BuscarContasBancarias lista as contas bancárias da empresa.
func BuscarContasBancarias(ctx context.Context, g Getter) (*ContasBancarias, error) {
	return get[ContasBancarias](ctx, g, PathContasBancarias)
}
