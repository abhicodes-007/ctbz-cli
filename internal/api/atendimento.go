package api

import "context"

const (
	PathChamadosEmAndamento = "atendimento/chamados?em-andamento=true"
	PathChamadosFinalizados = "atendimento/chamados?finalizados=true"
)

// LinkChamado é a página do chamado na central de ajuda da Contabilizei.
func LinkChamado(id string) string {
	return "https://suporte.contabilizei.com.br/hc/pt-br/requests/" + id
}

// Chamado é um chamado de atendimento (tela "Meus chamados"). Os campos são os lidos pelo
// front; nenhuma resposta real com itens foi vista, por isso todos são opcionais.
type Chamado struct {
	ID              any    `json:"id" contract:"optional"` // número ou texto
	Assunto         string `json:"assunto" contract:"optional"`
	Canal           string `json:"canal" contract:"optional"`
	Status          string `json:"status" contract:"optional"`
	Atualizado      string `json:"atualizado" contract:"optional"`
	PrevisaoRetorno string `json:"previsaoRetorno" contract:"optional"`
}

// BuscarChamados lê os chamados em andamento ou, com finalizados, os finalizados.
// Para contas com muitos chamados, os finalizados falham no servidor (HTTP 400); veja
// ADR-0013 e DadosEmpresa.Chamados.
func BuscarChamados(ctx context.Context, g Getter, finalizados bool) ([]Chamado, error) {
	path := PathChamadosEmAndamento
	if finalizados {
		path = PathChamadosFinalizados
	}
	c, err := get[[]Chamado](ctx, g, path)
	if err != nil {
		return nil, err
	}
	return *c, nil
}
