// Package otp obtém o código de verificação enviado pela Contabilizei,
// seja executando um comando externo (ex.: um script que lê o Gmail) ou
// perguntando ao usuário no terminal.
package otp

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var reCode = regexp.MustCompile(`\b(\d{6})\b`)

// Extract retorna o primeiro código de 6 dígitos encontrado no texto.
func Extract(s string) (string, bool) {
	m := reCode.FindStringSubmatch(s)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// Command executa um comando de shell até ele imprimir um código de 6 dígitos.
//
// O comando recebe as variáveis de ambiente:
//
//	CTBZ_OTP_SINCE      instante (epoch, segundos) em que o código foi solicitado
//	CTBZ_OTP_SINCE_ISO  o mesmo instante em RFC 3339
//
// para que possa ignorar e-mails antigos. Uma saída sem código ou um status de
// saída diferente de zero significa "ainda não chegou": o comando é repetido a
// cada Interval até Timeout.
type Command struct {
	Shell    string // padrão: "sh"
	Cmd      string
	Since    time.Time
	Interval time.Duration
	Timeout  time.Duration
	Log      io.Writer // progresso (opcional)
}

func (c *Command) Fetch(ctx context.Context) (string, error) {
	shell := c.Shell
	if shell == "" {
		shell = "sh"
	}
	interval, timeout := c.Interval, c.Timeout
	if interval <= 0 {
		interval = 5 * time.Second
	}
	if timeout <= 0 {
		timeout = 2 * time.Minute
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var lastErr error
	for attempt := 1; ; attempt++ {
		code, err := c.runOnce(ctx, shell)
		if err == nil {
			return code, nil
		}
		// Um comando morto pelo prazo esconderia o erro realmente útil.
		if ctx.Err() == nil || lastErr == nil {
			lastErr = err
		}
		if c.Log != nil {
			fmt.Fprintf(c.Log, "aguardando OTP (tentativa %d): %v\n", attempt, err)
		}
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("OTP não obtido em %s; último erro: %w", timeout, lastErr)
		case <-time.After(interval):
		}
	}
}

func (c *Command) runOnce(ctx context.Context, shell string) (string, error) {
	cmd := exec.CommandContext(ctx, shell, "-c", c.Cmd)
	cmd.Env = append(os.Environ(),
		"CTBZ_OTP_SINCE="+strconv.FormatInt(c.Since.Unix(), 10),
		"CTBZ_OTP_SINCE_ISO="+c.Since.UTC().Format(time.RFC3339),
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			return "", err
		}
		return "", fmt.Errorf("%w: %s", err, firstLine(msg))
	}
	if code, ok := Extract(stdout.String()); ok {
		return code, nil
	}
	return "", errors.New("comando não imprimiu um código de 6 dígitos")
}

// Prompt pergunta o código no terminal.
func Prompt(in io.Reader, out io.Writer, label string) (string, error) {
	r := bufio.NewReader(in)
	for {
		fmt.Fprint(out, label)
		line, err := r.ReadString('\n')
		if code, ok := Extract(line); ok {
			return code, nil
		}
		if err != nil {
			return "", fmt.Errorf("lendo OTP: %w", err)
		}
		fmt.Fprintln(out, "código inválido: digite os 6 dígitos")
	}
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
