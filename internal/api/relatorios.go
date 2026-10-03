package api

import (
	"context"
	"fmt"
)

// Relatórios contábeis: os três recebem ano e mês (1–12), nessa ordem, e devolvem a lista
// de contas da árvore contábil.
func PathBalancete(ano, mes int) string {
	return fmt.Sprintf("relatorios-ms/gerar-balancete/%d/%d", ano, mes)
}

func PathBalanco(ano, mes int) string {
	return fmt.Sprintf("relatorios-ms/gerarbalanco/%d/%d", ano, mes)
}

func PathRazao(ano, mes int) string {
	return fmt.Sprintf("relatorios-ms/gerarrazaotipoa/%d/%d", ano, mes)
}

// ContaRelatorio é uma conta de um relatório contábil. IdContaPai e Nivel formam a árvore
// (nível 1 = ATIVO, PASSIVO…); Tipo "S" é sintética (soma das filhas) e "A", analítica.
type ContaRelatorio struct {
	ID                     string               `json:"id"` // ex.: "1.01.01.01.00"
	IDContaPai             *string              `json:"idContaPai"`
	Descricao              string               `json:"descricao"`
	Natureza               string               `json:"natureza"` // DEVEDORA ou CREDORA
	Tipo                   string               `json:"tipo"`
	ClassificacaoConta     string               `json:"classificacaoConta"` // ATIVO, PASSIVO, PATRIMONIO, RESULTADO…
	Nivel                  int                  `json:"nivel"`
	SaldoAnterior          *float64             `json:"saldoAnterior"`
	TotalDebito            *float64             `json:"totalDebito"`
	TotalCredito           *float64             `json:"totalCredito"`
	SaldoExercicio         *float64             `json:"saldoExercicio"`
	SaldoExercicioAnterior *float64             `json:"saldoExercicioAnterior"`
	ListaLancamento        []LancamentoContabil `json:"listaLancamento"` // só no razão
}

// LancamentoContabil é um lançamento do razão.
type LancamentoContabil struct {
	ID                 int64    `json:"id"`
	Data               int64    `json:"data"` // epoch em ms
	Descricao          string   `json:"descricao"`
	Debito             *float64 `json:"debito"`
	Credito            *float64 `json:"credito"`
	SaldoExercicio     *float64 `json:"saldoExercicio"`
	Tipo               string   `json:"tipo"` // ex.: Sistema
	ContaContrapartida *struct {
		ID        string `json:"id"`
		Descricao string `json:"descricao"`
	} `json:"contaContrapartida"`
}

// BuscarRelatorio lê um relatório contábil (path de PathBalancete, PathBalanco ou PathRazao).
func BuscarRelatorio(ctx context.Context, g Getter, path string) ([]ContaRelatorio, error) {
	r, err := get[[]ContaRelatorio](ctx, g, path)
	if err != nil {
		return nil, err
	}
	return *r, nil
}
