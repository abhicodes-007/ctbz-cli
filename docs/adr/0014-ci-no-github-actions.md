# ADR-0014: CI no GitHub Actions com as mesmas verificações locais

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

Até a v0.8 cada commit era verificado localmente (`gofmt`, `go vet`, `go test`) antes do
push. Com contribuições de outras pessoas e releases automáticas (ADR-0003), a verificação
precisa rodar no servidor, para todo PR, e a documentação (ADR-0005) precisa ser publicada
sem passo manual.

## Decisão

- `.github/workflows/ci.yml` roda em todo PR e em push na `main`:
  - Go, em **Linux e macOS** (onde a CLI é usada; Windows só é compilado na release):
    `gofmt -l` vazio, `go vet ./...`, `go build ./...` e `go test -race ./...`;
  - documentação: `mkdocs build --strict` (links quebrados e páginas fora da navegação
    falham o build).
- Os testes do CI são os **offline**: contratos contra as fixtures, sem credenciais. Os
  contratos ao vivo ficam no job agendado de monitoramento (#53), que tem segredos.
- `.github/workflows/docs.yml` publica o site no GitHub Pages a cada push na `main`.
- Nenhuma ferramenta de lint além de `gofmt` e `go vet`: o projeto é pequeno e essas duas
  cobrem o que a revisão de código não pega; um linter novo entra por ADR.

## Consequências

- O que passa localmente passa no CI: os comandos são os mesmos do `docs/contribuindo.md`.
- O `-race` cobre as chamadas em paralelo do `ctbz resumo` e o re-login concorrente.
- Publicar a documentação exige configurar o Pages com a fonte "GitHub Actions" uma vez.

## Alternativas consideradas

- **golangci-lint** — mais regras, mas mais configuração e falsos positivos; pode entrar
  depois.
- **Matriz com Windows** — os testes usam `sh` (script de OTP); o binário de Windows é
  compilado pelo GoReleaser, sem testes.
