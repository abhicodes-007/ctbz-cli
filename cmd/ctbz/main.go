// Comando ctbz: CLI para a plataforma web da Contabilizei.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/edusouza/ctbz-cli/internal/ctbz"
	"golang.org/x/term"
)

const usage = `ctbz — CLI para a Contabilizei

Uso:
  ctbz login [--otp N] [--otp-cmd CMD] [--cnpj CNPJ]   autentica (usuário/senha + OTP)
  ctbz status                                          mostra a sessão atual e testa se ainda vale
  ctbz empresa [--json]                                dados da empresa selecionada
  ctbz api [-X MÉTODO] [-d CORPO] CAMINHO              chama uma URL da plataforma com a sessão
  ctbz logout                                          apaga a sessão local

Variáveis de ambiente:
  CTBZ_USER, CTBZ_PASSWORD   credenciais (e-mail ou CPF, e senha)
  CTBZ_OTP_CMD               comando que imprime o OTP (habilita login e re-login automáticos)
  CTBZ_CNPJ                  empresa a selecionar quando houver mais de uma
  CTBZ_HOME                  diretório da sessão (padrão: ~/.config/ctbz)
  CTBZ_VERBOSE=1             mostra progresso do --otp-cmd
`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	var err error
	args := os.Args[2:]
	switch os.Args[1] {
	case "login":
		err = cmdLogin(ctx, args)
	case "logout":
		err = cmdLogout()
	case "status":
		err = cmdStatus(ctx)
	case "empresa":
		err = cmdEmpresa(ctx, args)
	case "api":
		err = cmdAPI(ctx, args)
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "comando desconhecido: %s\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
	switch {
	case err == nil:
	case errors.Is(err, errPending):
		os.Exit(3)
	case errors.Is(err, flag.ErrHelp):
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

func cmdLogout() error {
	store, err := ctbz.DefaultStore()
	if err != nil {
		return err
	}
	if err := store.Clear(); err != nil {
		return err
	}
	fmt.Println("Sessão local removida.")
	return nil
}

func cmdStatus(ctx context.Context) error {
	store, err := ctbz.DefaultStore()
	if err != nil {
		return err
	}
	sess, err := store.LoadSession()
	if err != nil {
		return err
	}
	fmt.Println(describeSession(sess))
	fmt.Println("Login em:", sess.CreatedAt.Local().Format(time.DateTime))
	if exp := sess.ExpiresAt(); !exp.IsZero() {
		fmt.Println("Cookies expiram em:", exp.Local().Format(time.DateTime))
	}
	c := sess.Client()
	if _, err := c.API(ctx, "GET", "appbar/get", nil); err != nil {
		if errors.Is(err, ctbz.ErrUnauthorized) {
			fmt.Println("Situação: EXPIRADA — rode `ctbz login`")
			return nil
		}
		return err
	}
	saveCookies(store, sess, c)
	fmt.Println("Situação: válida")
	return nil
}

func cmdEmpresa(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("empresa", flag.ContinueOnError)
	asJSON := fs.Bool("json", false, "imprime a resposta JSON completa")
	if err := fs.Parse(args); err != nil {
		return err
	}
	resp, err := authedAPI(ctx, "GET", "dadosempresa/get", nil)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return &ctbz.HTTPError{Step: "dadosempresa/get", Status: resp.Status, Body: resp.Body}
	}
	if *asJSON {
		return printJSON(os.Stdout, resp.Body)
	}
	var data struct {
		EmpresaAtual struct {
			CNPJ               string   `json:"cnpj"`
			RazaoSocial        string   `json:"razaoSocial"`
			InscricaoMunicipal string   `json:"inscricaoMunicipal"`
			RegimeTributario   string   `json:"regimeTributario"`
			StatusEmpresa      string   `json:"statusEmpresa"`
			RamosAtividade     []string `json:"ramosAtividade"`
			Plano              string   `json:"plano"`
			Certificado        *struct {
				Status struct {
					Descricao string `json:"descricao"`
				} `json:"status"`
				DataValidade string `json:"dataValidade"`
			} `json:"certificado"`
		} `json:"empresaAtual"`
		Empresas []struct {
			CNPJ          string `json:"cnpj"`
			RazaoSocial   string `json:"razaoSocial"`
			StatusEmpresa string `json:"statusEmpresa"`
		} `json:"empresas"`
	}
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return fmt.Errorf("resposta inesperada de dadosempresa/get: %w", err)
	}
	e := data.EmpresaAtual
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Razão social:\t%s\n", e.RazaoSocial)
	fmt.Fprintf(w, "CNPJ:\t%s\n", e.CNPJ)
	fmt.Fprintf(w, "Situação:\t%s\n", e.StatusEmpresa)
	fmt.Fprintf(w, "Regime tributário:\t%s\n", e.RegimeTributario)
	if e.InscricaoMunicipal != "" {
		fmt.Fprintf(w, "Inscrição municipal:\t%s\n", e.InscricaoMunicipal)
	}
	if len(e.RamosAtividade) > 0 {
		fmt.Fprintf(w, "Ramos de atividade:\t%s\n", strings.Join(e.RamosAtividade, ", "))
	}
	if e.Plano != "" {
		fmt.Fprintf(w, "Plano:\t%s\n", e.Plano)
	}
	if e.Certificado != nil {
		fmt.Fprintf(w, "Certificado digital:\t%s (validade %s)\n", e.Certificado.Status.Descricao, e.Certificado.DataValidade)
	}
	if len(data.Empresas) > 1 {
		fmt.Fprintf(w, "Outras empresas:\t\n")
		for _, o := range data.Empresas {
			if ctbz.OnlyDigits(o.CNPJ) != ctbz.OnlyDigits(e.CNPJ) {
				fmt.Fprintf(w, "  %s\t%s (%s)\n", o.CNPJ, o.RazaoSocial, o.StatusEmpresa)
			}
		}
	}
	return w.Flush()
}

func cmdAPI(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("api", flag.ContinueOnError)
	method := fs.String("X", "GET", "método HTTP")
	data := fs.String("d", "", "corpo JSON da requisição (use @arquivo para ler de um arquivo ou @- para stdin)")
	raw := fs.Bool("raw", false, "não formata a resposta JSON")
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), `Uso: ctbz api [-X MÉTODO] [-d CORPO] CAMINHO

CAMINHO relativo é resolvido contra o BFF da plataforma (/api/plataforma/):
  ctbz api dadosempresa/get
  ctbz api menu/get
Caminhos absolutos alcançam as demais APIs do app:
  ctbz api /api/legado/...

Flags:
`)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return errors.New("informe exatamente um CAMINHO")
	}
	body, err := readBody(*data)
	if err != nil {
		return err
	}
	resp, err := authedAPI(ctx, strings.ToUpper(*method), fs.Arg(0), body)
	if err != nil {
		return err
	}
	if !*raw && json.Valid(resp.Body) {
		if err := printJSON(os.Stdout, resp.Body); err != nil {
			return err
		}
	} else {
		os.Stdout.Write(resp.Body)
	}
	if resp.Status >= 400 {
		return fmt.Errorf("HTTP %d", resp.Status)
	}
	return nil
}

// authedAPI faz uma chamada com a sessão salva. Se a sessão tiver expirado e
// CTBZ_OTP_CMD estiver configurado, refaz o login sozinho e repete a chamada.
func authedAPI(ctx context.Context, method, path string, body []byte) (*ctbz.Response, error) {
	store, err := ctbz.DefaultStore()
	if err != nil {
		return nil, err
	}
	sess, err := store.LoadSession()
	if err != nil {
		return nil, err
	}
	c := sess.Client()
	resp, err := c.API(ctx, method, path, bodyReader(body))
	if errors.Is(err, ctbz.ErrUnauthorized) && os.Getenv("CTBZ_OTP_CMD") != "" {
		fmt.Fprintln(os.Stderr, "Sessão expirada; refazendo login com CTBZ_OTP_CMD…")
		sess, err = login(ctx, store, loginOpts{
			otpCmd:     os.Getenv("CTBZ_OTP_CMD"),
			otpTimeout: envDuration("CTBZ_OTP_TIMEOUT", 3*time.Minute),
			cnpj:       sess.CNPJ,
			restart:    true,
		})
		if err != nil {
			return nil, err
		}
		c = sess.Client()
		resp, err = c.API(ctx, method, path, bodyReader(body))
	}
	if err != nil {
		return nil, err
	}
	saveCookies(store, sess, c)
	return resp, nil
}

// saveCookies persiste cookies renovados pelo servidor durante as chamadas.
func saveCookies(store *ctbz.Store, sess *ctbz.Session, c *ctbz.Client) {
	sess.Cookies = c.Cookies
	if err := store.SaveSession(sess); err != nil {
		fmt.Fprintln(os.Stderr, "aviso: não foi possível atualizar a sessão:", err)
	}
}

func bodyReader(b []byte) io.Reader {
	if b == nil {
		return nil
	}
	return strings.NewReader(string(b))
}

func readBody(arg string) ([]byte, error) {
	switch {
	case arg == "":
		return nil, nil
	case arg == "@-":
		return io.ReadAll(os.Stdin)
	case strings.HasPrefix(arg, "@"):
		return os.ReadFile(arg[1:])
	default:
		return []byte(arg), nil
	}
}

func printJSON(w io.Writer, data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		_, err = w.Write(data)
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

// describeSession resume a sessão a partir dos dados do localStorage.
func describeSession(sess *ctbz.Session) string {
	var login struct {
		Email   string `json:"email"`
		Empresa struct {
			CNPJ        string `json:"cnpj"`
			RazaoSocial string `json:"razaoSocial"`
		} `json:"empresa"`
	}
	if raw, ok := sess.Storage["l"]; ok {
		_ = json.Unmarshal(raw, &login)
	}
	parts := []string{}
	if login.Empresa.RazaoSocial != "" {
		parts = append(parts, fmt.Sprintf("%s (CNPJ %s)", login.Empresa.RazaoSocial, login.Empresa.CNPJ))
	} else if sess.CNPJ != "" {
		parts = append(parts, "CNPJ "+sess.CNPJ)
	}
	if login.Email != "" {
		parts = append(parts, "usuário "+login.Email)
	}
	if len(parts) == 0 {
		return "sessão ativa"
	}
	return strings.Join(parts, " — ")
}

func currentCNPJ(sess *ctbz.Session) string {
	for _, key := range []string{"e", "l"} {
		raw, ok := sess.Storage[key]
		if !ok {
			continue
		}
		var v struct {
			CNPJ    string `json:"cnpj"`
			Empresa struct {
				CNPJ string `json:"cnpj"`
			} `json:"empresa"`
		}
		if json.Unmarshal(raw, &v) == nil {
			if v.CNPJ != "" {
				return ctbz.OnlyDigits(v.CNPJ)
			}
			if v.Empresa.CNPJ != "" {
				return ctbz.OnlyDigits(v.Empresa.CNPJ)
			}
		}
	}
	return ""
}

func baseURL() string {
	if v := os.Getenv("CTBZ_BASE_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}
	return ctbz.DefaultBaseURL
}

func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func verboseWriter() io.Writer {
	if os.Getenv("CTBZ_VERBOSE") != "" {
		return os.Stderr
	}
	return nil
}
