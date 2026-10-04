package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/ctbz"
)

// sessionSender implementa api.Sender com a sessão salva: confere a sessão antes e envia
// cada escrita uma única vez, sem retentativa nem re-login depois do envio (ADR-0018).
type sessionSender struct{ s streams }

func (w sessionSender) Send(ctx context.Context, method, path string, body any, v any) error {
	data, err := api.EncodeBody(body)
	if err != nil {
		return err
	}
	contentType := ""
	if data != nil {
		contentType = "application/json"
	}
	return w.send(ctx, method, path, data, contentType, v)
}

func (w sessionSender) SendMultipart(ctx context.Context, method, path string, campos []api.Campo, v any) error {
	data, contentType, err := api.EncodeMultipart(campos)
	if err != nil {
		return err
	}
	return w.send(ctx, method, path, data, contentType, v)
}

func (w sessionSender) send(ctx context.Context, method, path string, body []byte, contentType string, v any) error {
	resp, err := sendOnce(ctx, w.s, method, path, body, contentType)
	if err != nil {
		return err
	}
	if resp.Status < 200 || resp.Status > 299 {
		return ctbz.NewWriteError(method, path, resp)
	}
	return api.DecodeResponse(resp.Body, v)
}

// errExpiredOnSend indica que o servidor recusou a escrita por sessão expirada: nada foi aplicado.
var errExpiredOnSend = errors.New("a sessão expirou antes do envio e a escrita não foi aplicada; rode `ctbz login` e repita")

// sendOnce confere a sessão com um GET (que pode refazer o login, porque nada foi enviado
// ainda) e então envia a escrita uma única vez. Respostas HTTP de erro voltam sem erro, para
// quem chama decidir; erros de rede avisam que o resultado da escrita é incerto.
func sendOnce(ctx context.Context, s streams, method, path string, body []byte, contentType string) (*ctbz.Response, error) {
	if _, err := authedAPI(ctx, s, "GET", api.PathAppBar, nil); err != nil {
		return nil, fmt.Errorf("conferindo a sessão antes de enviar: %w", err)
	}
	store, err := ctbz.DefaultStore()
	if err != nil {
		return nil, err
	}
	sess, err := store.LoadSession()
	if err != nil {
		return nil, err
	}
	c := sess.Client()
	resp, err := c.Send(ctx, method, path, bytes.NewReader(body), contentType)
	switch {
	case errors.Is(err, ctbz.ErrUnauthorized):
		return nil, errExpiredOnSend
	case err != nil:
		return nil, fmt.Errorf("a conexão falhou durante o envio e a escrita pode ou não ter sido aplicada; confira antes de repetir: %w", err)
	}
	saveCookies(store, sess, c, s.err)
	return resp, nil
}
