package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/ctbz"
	"github.com/edusouza/ctbz-cli/internal/output"
)

// statusFinalizados são os status de chamado (Zendesk) considerados finalizados.
var statusFinalizados = map[string]bool{"solved": true, "closed": true}

func newChamadosCmd() *cobra.Command {
	var finalizados bool
	cmd := &cobra.Command{
		Use:   "chamados",
		Short: "Lista os chamados de atendimento (em andamento ou finalizados)",
		Long: `Lista os chamados abertos com o atendimento da Contabilizei, com assunto, status, canal,
datas e o link para a central de ajuda.

--finalizados lista os já resolvidos. Em contas com muitos chamados, a Contabilizei não
consegue montar essa lista (erro do servidor); nesse caso a CLI mostra os finalizados
entre os 100 chamados mais recentes da empresa e avisa no stderr.`,
		Example: `  ctbz chamados
  ctbz chamados --finalizados -o csv`,
		Args: exactArgs(0, "nenhum argumento"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			g := sessionGetter{s}
			cs, err := api.BuscarChamados(cmd.Context(), g, finalizados)
			var httpErr *ctbz.HTTPError
			if finalizados && errors.As(err, &httpErr) {
				// Fonte alternativa (ADR-0013): o servidor falha sempre para listas grandes.
				fmt.Fprintf(s.err, "aviso: a Contabilizei não listou os chamados finalizados (HTTP %d); "+
					"mostrando os finalizados entre os 100 mais recentes\n", httpErr.Status)
				d, err := api.BuscarDadosEmpresa(cmd.Context(), g)
				if err != nil {
					return err
				}
				return output.Write(s.out, f, chamadosRecentesList(d.Chamados))
			}
			if err != nil {
				return err
			}
			return output.Write(s.out, f, chamadosList(cs))
		},
	}
	cmd.Flags().BoolVar(&finalizados, "finalizados", false, "lista os chamados finalizados")
	return cmd
}

func novaListaDeChamados() *output.List {
	return &output.List{Columns: []output.Column{
		{Key: "id", Header: "ID"},
		{Key: "assunto", Header: "Assunto"},
		{Key: "status", Header: "Status"},
		{Key: "canal", Header: "Canal"},
		{Key: "criado", Header: "Criado em"},
		{Key: "atualizado", Header: "Atualizado"},
		{Key: "previsao_retorno", Header: "Previsão de retorno"},
		{Key: "link", Header: "Link"},
	}}
}

func chamadosList(cs []api.Chamado) *output.List {
	l := novaListaDeChamados()
	for _, c := range cs {
		id := idDeChamado(c.ID)
		l.Append(nilIfEmpty(id), output.Text(c.Assunto), nilIfEmpty(c.Status), nilIfEmpty(c.Canal), nil,
			dataOuTexto(c.Atualizado), dataOuTexto(c.PrevisaoRetorno), linkDeChamado(id))
	}
	return l
}

// chamadosRecentesList lista os finalizados entre os chamados de dadosempresa/get.
func chamadosRecentesList(cs []api.ChamadoResumo) *output.List {
	l := novaListaDeChamados()
	for _, c := range cs {
		if !statusFinalizados[c.Status] {
			continue
		}
		l.Append(nilIfEmpty(c.ID), output.Text(c.Subject), c.Status, nil, dataOuTexto(c.CreatedAt),
			nil, nil, linkDeChamado(c.ID))
	}
	return l
}

// idDeChamado normaliza o id, que pode vir como número ou texto.
func idDeChamado(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case nil:
		return ""
	}
	return fmt.Sprint(v)
}

func linkDeChamado(id string) any {
	if id == "" {
		return nil
	}
	return api.LinkChamado(id)
}

// dataOuTexto converte datas conhecidas; outro texto volta como veio.
func dataOuTexto(s string) any {
	if d, ok := output.ParseDate(s); ok {
		return d
	}
	return nilIfEmpty(s)
}
