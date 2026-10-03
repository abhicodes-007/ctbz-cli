# ADR-0015: Releases com GoReleaser e notas tiradas do CHANGELOG

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

Até a v0.8 a CLI só podia ser instalada com `go install`, o que exige Go. A v1.0 precisa de
binários prontos para Linux, macOS e Windows, e de notas de release que não dupliquem o
trabalho do `CHANGELOG.md` (ADR-0003).

## Decisão

- Uma tag `vX.Y.Z` dispara `.github/workflows/release.yml`, que roda os testes e o
  GoReleaser (`.goreleaser.yaml`): binários `CGO_ENABLED=0` para linux, darwin e windows em
  amd64 e arm64, arquivos `tar.gz` (`zip` no Windows) com README e CHANGELOG, e
  `checksums.txt`.
- A versão é injetada com `-X main.version={{ .Tag }}`; commit e data vêm das informações
  de VCS que o Go embute no build (`ctbz version`).
- As notas da release são a seção da versão no `CHANGELOG.md`, extraída por
  `go run ./tools/releasenotes vX.Y.Z`. O changelog automático do GoReleaser fica desligado.

## Consequências

- Lançar uma versão continua sendo: cortar o CHANGELOG no PR de release e criar a tag depois
  do merge (ADR-0003). O resto é automático.
- Um teste garante que toda versão do CHANGELOG tem notas; uma release sem notas falha.
- Os pacotes não incluem licença enquanto o repositório não tiver uma.

## Alternativas consideradas

- **Changelog gerado a partir das issues e PRs** — duplicaria o CHANGELOG, que é escrito
  para usuários e já é obrigatório em cada PR.
- **Homebrew tap** — exige outro repositório e um token com escrita nele; fica para quando
  houver demanda.
- **Script shell para as notas** — `awk`/`sed` sem teste; uma ferramenta Go segue o padrão
  de `tools/` e é testada no CI.
