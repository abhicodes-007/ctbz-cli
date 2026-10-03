# ADR-0010: Testes de contrato derivados dos tipos, com fixtures anonimizadas

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

A API é privada e pode mudar sem aviso. Precisamos detectar cedo quando um campo usado pela
CLI some ou muda de tipo, sem manter um esquema paralelo e sem guardar dados reais no repositório.

## Decisão

- **O contrato é o tipo Go** da resposta ([ADR-0009](0009-camada-api-tipada.md)):
  `internal/contract.Check` percorre o struct por reflexão e compara com o JSON.
  - Campo com tag `json` é obrigatório; `contract:"optional"` torna opcional.
  - `null` é aceito em qualquer campo (a decodificação tolera).
  - Divergências: *campo removido* e *tipo mudou* falham o teste; *campo novo* é só informativo.
  - Listas: verificados os 5 primeiros itens.
- **Fixtures** em `internal/api/testdata/<nome>.json`, geradas por `go run ./tools/capture`
  a partir da API real (ou de um arquivo com `-from`) e passadas por `contract.Anonymize`:
  nomes, documentos, e-mails, endereços, contas, valores monetários, textos livres e IDs
  internos viram valores fictícios; listas são cortadas em 3 itens; textos longos, omitidos.
  **Toda fixture é revisada antes do commit.**
- **Modo ao vivo:** `CTBZ_CONTRACT_LIVE=1 go test ./internal/api -run Live` usa a sessão do
  `ctbz login` e só faz `GET` dos endpoints com `LivePath`.

## Consequências

- Mudar um tipo atualiza o contrato automaticamente; não há esquema duplicado.
- O teste de fixtures protege contra regressões na decodificação; o modo ao vivo, contra
  mudanças da API.
- A anonimização é por regras: um campo pessoal com nome inesperado pode passar. Por isso a
  revisão manual é obrigatória e as regras crescem quando um caso novo aparece.

## Alternativas consideradas

- **JSON Schema por endpoint** — duplicaria os tipos e envelheceria separado deles.
- **Fixtures escritas à mão** — não refletem a API real; captura + anonimização é mais fiel.
