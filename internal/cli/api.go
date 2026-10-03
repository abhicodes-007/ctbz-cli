package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/output"
)

func newAPICmd() *cobra.Command {
	var method, data string
	var raw bool
	cmd := &cobra.Command{
		Use:   "api CAMINHO",
		Short: "Chama uma URL da plataforma com a sessão atual",
		Long: `Chama qualquer endpoint da plataforma com os cookies da sessão.

CAMINHO relativo é resolvido contra o BFF da plataforma (/api/plataforma/);
caminhos absolutos ("/api/legado/...") alcançam as demais APIs do app.

A resposta sai em JSON formatado (CTBZ_OUTPUT é ignorado). Com -o table ou -o csv,
listas de objetos viram tabelas. --raw imprime o corpo exatamente como veio.
Respostas HTTP 4xx/5xx terminam com código de saída 1.`,
		Example: `  ctbz api dadosempresa/get
  ctbz api -o table menu/get
  ctbz api /api/legado/empresa/cnpj?cnpj=00000000000100
  ctbz api -X POST -d @corpo.json caminho/qualquer`,
		Args: exactArgs(1, "exatamente um CAMINHO"),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := outputFormat(cmd, string(output.FormatJSON))
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			body, err := readBody(data, s.in)
			if err != nil {
				return err
			}
			resp, err := authedAPI(cmd.Context(), s, strings.ToUpper(method), args[0], body)
			if err != nil {
				return err
			}
			if parsed, perr := output.FromJSON(resp.Body); !raw && perr == nil {
				if err := output.Write(s.out, f, parsed); err != nil {
					return err
				}
			} else {
				s.out.Write(resp.Body)
				if !raw && len(resp.Body) > 0 && resp.Body[len(resp.Body)-1] != '\n' {
					io.WriteString(s.out, "\n")
				}
			}
			if resp.Status >= 400 {
				return fmt.Errorf("HTTP %d", resp.Status)
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&method, "method", "X", "GET", "método HTTP")
	f.StringVarP(&data, "data", "d", "", "corpo JSON da requisição (@arquivo lê de um arquivo, @- da entrada padrão)")
	f.BoolVar(&raw, "raw", false, "imprime o corpo da resposta como veio, sem formatar")
	return cmd
}

func readBody(arg string, stdin io.Reader) ([]byte, error) {
	switch {
	case arg == "":
		return nil, nil
	case arg == "@-":
		return io.ReadAll(stdin)
	case strings.HasPrefix(arg, "@"):
		return os.ReadFile(arg[1:])
	default:
		return []byte(arg), nil
	}
}
