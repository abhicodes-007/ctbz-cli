# ADR-0004: Branches empilhadas, um PR por issue e commits atômicos

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

O roadmap tem dezenas de sub-issues. PRs grandes são difíceis de revisar; PRs independentes
sobre `main` conflitariam entre si, porque as funcionalidades se apoiam umas nas outras.

## Decisão

- **Um PR por issue**, com o corpo terminando em `Closes #N`.
- **Pilha:** cada branch nasce da anterior e o PR aponta para ela
  (`main ← v0.1/00-processo ← v0.1/01-cobra ← v0.1/06-testes-contrato ← …`).
  Fazer merge de baixo para cima; ao fundir um PR, o GitHub redireciona o seguinte para `main`.
- **Nome da branch:** `vX.Y/NN-slug`, com `NN` = número da issue (ou `00`/`0x` para trabalho de base sem issue).
- **Commits atômicos:** cada commit compila, passa em `go vet` e `go test ./...` e faz uma coisa só.
  Mensagens em português no padrão [Conventional Commits](https://www.conventionalcommits.org/pt-br/)
  (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`), com `#N` da issue quando houver.
- **Correções de revisão** num PR da base da pilha entram como commit novo nessa branch e são
  propagadas para cima com `git rebase --update-refs` (sem reescrever o que já foi revisado acima).

## Consequências

- Revisão rápida: cada PR mostra só o diff da sua issue.
- Custo: rebase em cascata quando um PR de baixo muda. Mitigado mantendo os PRs pequenos.

## Alternativas consideradas

- **Um PR por épico** — 9 PRs grandes, difíceis de revisar.
- **PRs independentes sobre `main`** — conflitos constantes entre funcionalidades dependentes.
