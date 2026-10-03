package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/edusouza/ctbz-cli/internal/api"
)

func TestCertificadoRecord(t *testing.T) {
	hoje := time.Date(2027, 7, 12, 15, 0, 0, 0, time.Local)
	venc := time.Date(2027, 7, 22, 9, 0, 0, 0, time.Local).UnixMilli()
	r := certificadoRecord(&api.CertificadoStatus{Situacao: "VALIDO", Valido: true, DataVencimento: &venc}, hoje)
	got := map[string]any{}
	for _, f := range r.Fields {
		got[f.Key] = f.Value
	}
	if got["dias_para_vencer"] != 10 {
		t.Errorf("dias_para_vencer = %v", got["dias_para_vencer"])
	}
	semCert := certificadoRecord(&api.CertificadoStatus{Situacao: "SEM_CERTIFICADO"}, hoje)
	for _, f := range semCert.Fields {
		if f.Key == "dias_para_vencer" && f.Value != nil {
			t.Errorf("sem certificado, dias deveria ser nulo: %v", f.Value)
		}
	}
}

func TestEmpresaCertificadoCmd(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/certificado/status": `{"situacao":"VALIDO","dataVencimento":1816282860000,"mensagemErro":null,"valido":true,"aptoRenovacao":false}`,
	}))
	out, stderr, code := execCLI(t, "", "empresa", "certificado", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	if !strings.HasPrefix(out, "situacao,valido,vencimento,dias_para_vencer,apto_renovacao,mensagem\nVALIDO,true,2027-07-") {
		t.Errorf("CSV:\n%s", out)
	}
}

func TestDiasEntre(t *testing.T) {
	a := time.Date(2026, 12, 31, 23, 0, 0, 0, time.Local)
	for b, want := range map[time.Time]int{
		time.Date(2027, 1, 1, 1, 0, 0, 0, time.Local):   1,
		time.Date(2026, 12, 31, 0, 0, 0, 0, time.Local): 0,
		time.Date(2026, 12, 1, 0, 0, 0, 0, time.Local):  -30,
	} {
		if got := diasEntre(a, b); got != want {
			t.Errorf("diasEntre(%v, %v) = %d, quero %d", a, b, got, want)
		}
	}
}
