# Testes de contrato

Como a CLI detecta mudanças na API da Contabilizei. Decisões em
[ADR-0009](adr/0009-camada-api-tipada.md) e [ADR-0010](adr/0010-testes-de-contrato.md).

## Rodar

```sh
go test ./...                                          # contratos contra as fixtures
CTBZ_CONTRACT_LIVE=1 go test ./internal/api -run Live -v   # contra a API real (precisa de ctbz login)
```

Saída de uma quebra:

```text
--- FAIL: TestContractsLive/dadosempresa
    empresaAtual.cnpj: campo removido (esperado texto)
    empresaAtual.plano: tipo mudou (esperado texto, veio objeto)
```

## Adicionar um endpoint

1. Em `internal/api/<contexto>.go`: constante `PathXxx`, tipo da resposta só com os campos
   usados e função `BuscarXxx(ctx, g)`.
2. Registrar em `Endpoints()` (`internal/api/api.go`). `LivePath` vazio quando o caminho
   depende de um ID de outra resposta.
3. Gerar a fixture: `go run ./tools/capture NOME` (ou `-from resposta.json`, ou `-path CAMINHO`).
   A captura poda a resposta aos campos do tipo e depois anonimiza
   ([ADR-0011](adr/0011-fixtures-podadas-ao-contrato.md)); `-full` mantém tudo.
4. **Revisar a fixture** (checklist abaixo) e rodar `go test ./...`.

## Checklist de revisão de fixture

- [ ] Nenhum nome de pessoa ou empresa real, CPF, CNPJ, e-mail, telefone, endereço, conta bancária
- [ ] Nenhum valor monetário real (aparecem como `1000` ou `1234.56`)
- [ ] Nenhum texto livre real (assuntos, mensagens, observações viram `TEXTO EXEMPLO`)
- [ ] Nenhuma URL assinada ou token
- [ ] Se algo passou, acrescentar uma regra em `internal/contract/anon.go` (com teste) e capturar de novo

## Campos opcionais

Quando a API omite um campo em parte das respostas (ex.: itens de menu sem `children`), marque
`contract:"optional"` no campo. Campos que podem vir `null` não precisam de marcação.

## Campos não verificados

Quando a conta usada para capturar não tem exemplos de um campo (ex.: uma lista sempre vazia),
declare-o como `json.RawMessage`. O contrato aceita qualquer valor nele, a poda o mantém
inteiro (a anonimização continua valendo) e o comando o mostra como a API o devolve
(`output.FromJSON`). Quando houver dados reais, troque por um tipo.

## Monitoramento

Os contratos ao vivo rodam toda semana no job de monitoramento, que abre uma issue quando
algum quebra. Ver [Monitoramento da API](monitoramento.md).
