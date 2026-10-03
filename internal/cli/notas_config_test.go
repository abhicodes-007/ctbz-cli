package cli

import (
	"strings"
	"testing"
)

func TestNotasConfig(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/novo-emissor/listagem/init":     fixture(t, "emissor_init"),
		"/api/plataforma/novo-emissor/v2/versao-emissor": fixture(t, "versao_emissor"),
		"/api/plataforma/notafiscal/listaliquotaatividade": `{"regimeTributario":"SIMPLES",
"interno":{"servico":[{"codigoCnae":"6204000","descricaoCnae":"Consultoria em TI","codigoItemServico":"01.06",
 "descricaoItemServico":"Assessoria e consultoria em informática.","aliquotaBase":6,"aliquotaISS":2.01,"fatorR":33.3333,"anexoFixo":false}],
 "comercio":[]},
"externo":{"servico":[{"codigoCnae":"6204000","descricaoCnae":"Consultoria em TI","codigoItemServico":"01.06",
 "descricaoItemServico":"Assessoria e consultoria em informática.","aliquotaBase":3.05,"aliquotaISS":0,"fatorR":null,"anexoFixo":true}],
 "comercio":[]}}`,
	}))

	out, stderr, code := execCLI(t, "", "notas", "config", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	want := "emissor_habilitado,versao,instabilidade,permite_exterior,certificado_validade,certificado_vencido," +
		"certificado_em_renovacao,municipio,uf\n" +
		"true,V2,false,true,2027-07-22,false,false,FULANO DE TAL,PR\n"
	if out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}

	out, _, code = execCLI(t, "", "notas", "aliquotas", "-o", "csv")
	want = "mercado,tipo,cnae,atividade,item_servico,descricao_item,aliquota,aliquota_iss,fator_r,anexo_fixo\n" +
		"interno,servico,6204-0/00,Consultoria em TI,01.06,Assessoria e consultoria em informática.,6,2.01,33.3333,false\n" +
		"externo,servico,6204-0/00,Consultoria em TI,01.06,Assessoria e consultoria em informática.,3.05,0,,true\n"
	if code != ExitOK || out != want {
		t.Errorf("alíquotas (código %d):\n%s\nesperado:\n%s", code, out, want)
	}
}

func TestNotasAliquotasFixture(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/notafiscal/listaliquotaatividade": fixture(t, "aliquotas_emissor")}))
	out, _, code := execCLI(t, "", "notas", "aliquotas", "-o", "json")
	if code != ExitOK || !strings.Contains(out, `"mercado": "externo"`) || !strings.Contains(out, `"cnae": "6204-0/00"`) {
		t.Errorf("código %d:\n%s", code, out)
	}
}
