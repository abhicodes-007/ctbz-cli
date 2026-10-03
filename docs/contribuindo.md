# Contribuindo

Convenções do dia a dia. Decisões maiores estão nos [ADRs](adr/README.md).

## Fluxo

1. Uma issue por funcionalidade; um PR por issue, empilhado sobre o anterior
   ([ADR-0004](adr/0004-branches-empilhadas-e-prs-por-issue.md)).
2. Branch `vX.Y/NN-slug`; commits atômicos em [Conventional Commits](https://www.conventionalcommits.org/pt-br/),
   em português.
3. Cada commit compila e passa em `gofmt -l .`, `go vet ./...` e `go test ./...`.

## Definição de pronto de um PR

- [ ] Testes cobrindo o comportamento novo (incluindo erro e caso vazio)
- [ ] Documentação de usuário atualizada (`README.md` e/ou página em `docs/`, com entrada no `mkdocs.yml`)
- [ ] Referência de comandos regenerada, se a CLI mudou (`go run ./tools/gendocs`)
- [ ] Linha no `CHANGELOG.md` em `[Unreleased]` quando o usuário perceber a mudança ([ADR-0003](adr/0003-versionamento-e-changelog.md))
- [ ] ADR novo, se houve decisão com alternativas razoáveis
- [ ] Nenhum dado pessoal real (CPF, CNPJ, nomes, e-mails, endereços, valores) em código, testes ou docs

## Idioma e nomes

- Textos para o usuário (mensagens, ajuda, docs, CHANGELOG) e commits: **português do Brasil**.
- Código: termos do **domínio em português** (`Guia`, `Empresa`, `Competencia`), termos
  **técnicos em inglês** (`Client`, `Do`, `Write`). Sem acentos em identificadores.
- Chaves de saída JSON/CSV: `snake_case` em português sem acento ([ADR-0006](adr/0006-saida-padronizada.md)).
- Mensagens de erro: minúsculas, sem ponto final, com contexto (`fmt.Errorf("lendo guias: %w", err)`).

## Testes

- Tabelas de casos (`for _, tc := range ...`) e `t.Run` quando houver vários cenários.
- HTTP: `httptest.Server`; nunca chamar a Contabilizei real em `go test` comum.
- Respostas da API: fixtures **anonimizadas** em `internal/api/testdata/`, geradas e
  revisadas como descrito em [Testes de contrato](contratos.md).
- Formatos de saída: golden files (`go test ./internal/output -update` regrava).

## Segurança

- Só `GET` nos comandos de domínio ([ADR-0002](adr/0002-somente-leitura-ate-1-0.md)).
- Sessão e credenciais nunca vão para logs, testes ou commits.
