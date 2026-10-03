package ctbz

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestSaveSessionConcorrente(t *testing.T) {
	store := &Store{Dir: t.TempDir()}
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sess := &Session{CNPJ: "22222222000122", CreatedAt: time.Unix(int64(i), 0)}
			if err := store.SaveSession(sess); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	sess, err := store.LoadSession()
	if err != nil || sess.CNPJ != "22222222000122" {
		t.Fatalf("sessão corrompida: %+v %v", sess, err)
	}
	info, err := os.Stat(filepath.Join(store.Dir, "session.json"))
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("permissão %v (%v), quero 0600", info.Mode().Perm(), err)
	}
	if tmps, _ := filepath.Glob(filepath.Join(store.Dir, "*.tmp")); len(tmps) > 0 {
		t.Errorf("temporários esquecidos: %v", tmps)
	}
}
