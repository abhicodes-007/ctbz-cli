package api

import "context"

const PathCentralRotinas = "dashboard/v2/central-rotinas"

// IndicadorPendencia é um item de pendenciasCriticas/outrasPendencias: só o indicador é
// comum a todos; os demais campos variam por tipo de pendência.
type IndicadorPendencia struct {
	PossuiPendencia bool `json:"possuiPendencia"`
}

// RotinaEmpresa é uma rotina mensal de responsabilidade da empresa.
type RotinaEmpresa struct {
	Tipo         string `json:"tipo"`  // ex.: IMPORTACAO_EXTRATO, IMPOSTO, VENCIMENTO_MENSALIDADE
	Prazo        string `json:"prazo"` // AAAA-MM-DD
	Status       string `json:"status"`
	Automatica   bool   `json:"automatica"`
	Propriedades *struct {
		TituloModal    string   `json:"tituloModal" contract:"optional"`
		MesReferencia  string   `json:"mesReferencia" contract:"optional"`
		ValorPagamento *float64 `json:"valorPagamento" contract:"optional"`
	} `json:"propriedades"`
}

// RotinaContabilizei é uma obrigação entregue pela Contabilizei (eSocial, DCTFWeb…).
type RotinaContabilizei struct {
	Titulo   string `json:"titulo"`
	Prazo    string `json:"prazo"` // AAAA-MM-DD
	Status   string `json:"status"`
	Conteudo struct {
		Sigla string `json:"sigla"`
	} `json:"conteudo"`
}

// CentralRotinas é a resposta de dashboard/v2/central-rotinas.
type CentralRotinas struct {
	Pendencias struct {
		Criticas map[string]IndicadorPendencia `json:"pendenciasCriticas"`
		Outras   map[string]IndicadorPendencia `json:"outrasPendencias"`
	} `json:"pendencias"`
	Rotinas             []RotinaEmpresa      `json:"rotinas"`
	RotinasContabilizei []RotinaContabilizei `json:"rotinasContabilizei"`
}

// BuscarCentralRotinas lê pendências críticas e rotinas do painel.
func BuscarCentralRotinas(ctx context.Context, g Getter) (*CentralRotinas, error) {
	return get[CentralRotinas](ctx, g, PathCentralRotinas)
}
