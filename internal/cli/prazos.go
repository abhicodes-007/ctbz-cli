package cli

import "github.com/edusouza/ctbz-cli/internal/output"

// diasDeAviso é a janela em que um prazo aberto é marcado como "próxima" (ADR-0012).
const diasDeAviso = 7

// Valores da coluna alerta (ADR-0012).
const (
	alertaVencida = "vencida"
	alertaProxima = "próxima"
)

// alertaDePrazo classifica o prazo de um item aberto; itens concluídos ou sem prazo não
// têm alerta (nil).
func alertaDePrazo(prazo output.Date, aberto bool) any {
	if !aberto || prazo.IsZero() {
		return nil
	}
	switch dias := diasEntre(now(), prazo.Time); {
	case dias < 0:
		return alertaVencida
	case dias <= diasDeAviso:
		return alertaProxima
	}
	return nil
}
