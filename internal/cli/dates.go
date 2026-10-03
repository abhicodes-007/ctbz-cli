package cli

import (
	"fmt"
	"time"

	"github.com/edusouza/ctbz-cli/internal/output"
)

// now é trocado nos testes para fixar a data atual.
var now = time.Now

// dateFromMillis converte epoch em milissegundos (formato das APIs); zero vira data vazia.
func dateFromMillis(ms int64) output.Date {
	if ms == 0 {
		return output.Date{}
	}
	return output.DateFromMillis(ms)
}

// diasEntre conta dias de calendário de a até b (negativo se b já passou).
func diasEntre(a, b time.Time) int {
	da := time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, time.UTC)
	db := time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, time.UTC)
	return int(db.Sub(da).Hours() / 24)
}

// competencia formata mês e ano como "07/2026"; zero vira vazio.
func competencia(mes, ano int) any {
	if mes == 0 || ano == 0 {
		return nil
	}
	return fmt.Sprintf("%02d/%d", mes, ano)
}
