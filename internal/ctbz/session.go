package ctbz

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Session é a sessão autenticada persistida em disco.
type Session struct {
	BaseURL   string                  `json:"base_url"`
	CNPJ      string                  `json:"cnpj"`
	CreatedAt time.Time               `json:"created_at"`
	Cookies   map[string]*http.Cookie `json:"cookies"`
	// Conteúdo que o navegador gravaria no localStorage (l, r, e).
	Storage map[string]json.RawMessage `json:"storage,omitempty"`
}

// Pending é um login interrompido aguardando OTP ou escolha de empresa.
type Pending struct {
	BaseURL    string      `json:"base_url"`
	WantedCNPJ string      `json:"wanted_cnpj,omitempty"`
	State      *LoginState `json:"state"`
}

// Store localiza os arquivos de estado da CLI.
type Store struct{ Dir string }

// DefaultStore usa $CTBZ_HOME ou ~/.config/ctbz.
func DefaultStore() (*Store, error) {
	if d := os.Getenv("CTBZ_HOME"); d != "" {
		return &Store{Dir: d}, nil
	}
	cfg, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return &Store{Dir: filepath.Join(cfg, "ctbz")}, nil
}

func (s *Store) sessionPath() string { return filepath.Join(s.Dir, "session.json") }
func (s *Store) pendingPath() string { return filepath.Join(s.Dir, "pending.json") }

func (s *Store) LoadSession() (*Session, error) {
	var sess Session
	if err := readJSON(s.sessionPath(), &sess); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("nenhuma sessão salva; rode `ctbz login`")
		}
		return nil, err
	}
	return &sess, nil
}

func (s *Store) SaveSession(sess *Session) error { return writeJSON(s.Dir, s.sessionPath(), sess) }

func (s *Store) LoadPending() (*Pending, error) {
	var p Pending
	if err := readJSON(s.pendingPath(), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) SavePending(p *Pending) error { return writeJSON(s.Dir, s.pendingPath(), p) }

func (s *Store) ClearPending() error { return removeIfExists(s.pendingPath()) }

func (s *Store) Clear() error {
	return errors.Join(removeIfExists(s.sessionPath()), removeIfExists(s.pendingPath()))
}

// ClientFor cria um Client já autenticado com os cookies da sessão.
func (sess *Session) Client() *Client {
	c := NewClient(sess.BaseURL)
	for k, v := range sess.Cookies {
		c.Cookies[k] = v
	}
	return c
}

// ExpiresAt retorna a menor expiração entre os cookies da sessão.
func (sess *Session) ExpiresAt() time.Time {
	var min time.Time
	for _, ck := range sess.Cookies {
		if ck.Expires.IsZero() {
			continue
		}
		if min.IsZero() || ck.Expires.Before(min) {
			min = ck.Expires
		}
	}
	return min
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// writeJSON grava com permissão 0600 (o arquivo contém cookies de sessão), de forma
// atômica: cada escrita usa um temporário próprio, então gravações em paralelo (ex.: as
// chamadas simultâneas do ctbz resumo) nunca misturam conteúdo; a última vence.
func writeJSON(dir, path string, v any) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp") // criado com 0600
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(f.Name(), path)
	}
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}

func removeIfExists(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
