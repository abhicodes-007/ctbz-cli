# Monitoramento da API

A CLI depende de APIs internas do painel, que a Contabilizei pode mudar sem aviso. Um job
agendado (`.github/workflows/monitor.yml`, toda segunda-feira, ou manualmente em
**Actions → Monitoramento → Run workflow**) verifica:

1. **Contratos ao vivo**: os mesmos testes de contrato do `go test`, mas contra a API real
   (`CTBZ_CONTRACT_LIVE=1`, só `GET`). Falham quando um campo usado pela CLI some ou muda de
   tipo.
2. **Catálogo de endpoints**: `scripts/extrair-endpoints.py` baixa os bundles do painel e
   regenera `docs/api/catalogo.md`; qualquer diferença com o versionado indica que o front
   passou a chamar outros endpoints.

Quando algo muda, o job abre a issue **"Monitoramento: a Contabilizei mudou a API ou o
front"** com o relatório (ou comenta nela, se já estiver aberta).

## Configurar

Em **Settings → Secrets and variables → Actions**, crie:

| Segredo | Conteúdo |
|---|---|
| `CTBZ_USER` | e-mail ou CPF do login |
| `CTBZ_PASSWORD` | senha |
| `CTBZ_CNPJ` | CNPJ da empresa (se o usuário tiver mais de uma) |
| `CTBZ_OTP_CMD` | comando que imprime o código OTP (ver abaixo) |

O login sempre pede um código por e-mail, então o job precisa de um `CTBZ_OTP_CMD` que
funcione sem ninguém por perto ([OTP automático](otp-automatico/README.md)). Opções:

- um serviço próprio que recebe o e-mail (encaminhado por um filtro do Gmail) e devolve o
  último código: `CTBZ_OTP_CMD='curl -fsS -H "Authorization: Bearer …" https://meu-servidor/ultimo-otp?desde=$CTBZ_OTP_SINCE'`;
- o script do Gmail com o `gws`, instalando e autenticando o `gws` num passo extra do
  workflow com credenciais de uma conta de serviço guardadas em outro segredo.

Os segredos ficam só no GitHub; o job não grava sessão nem credenciais no repositório.

## Rodar localmente

Com uma sessão válida (`ctbz login`):

```sh
sh scripts/monitorar.sh            # relatório no stdout; código 1 se algo mudou
git checkout docs/api/catalogo.md  # descarta o catálogo regenerado, se não for usar
```

## Quando a issue abrir

- **Contrato quebrado**: ajuste o tipo em `internal/api`, recapture a fixture
  (`go run ./tools/capture NOME`, revise o arquivo) e trate a mudança no comando.
- **Catálogo diferente**: veja os endpoints novos ou removidos, atualize
  `docs/api/catalogo.md` e, se for o caso, `docs/api/endpoints-verificados.md`.
