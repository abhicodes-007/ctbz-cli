# ADR-0005: Documentação em Markdown publicada com MkDocs Material

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

Toda a documentação do projeto (guia de uso, engenharia reversa, ADRs) vai virar um site no
GitHub Pages. Ela já existe em Markdown em `docs/` e precisa continuar legível no próprio GitHub.

## Decisão

- `docs/` é a fonte do site, publicado com [MkDocs Material](https://squidfunk.github.io/mkdocs-material/)
  a partir de `mkdocs.yml` na raiz.
- Links entre páginas são relativos a arquivos `.md` (funcionam no GitHub e no MkDocs).
  Links para arquivos fora de `docs/` (código, `CHANGELOG.md`, `ROADMAP.md`) usam URL absoluta do GitHub.
- A referência de comandos (`docs/referencia/`) é **gerada** a partir da própria CLI
  (`go run ./tools/gendocs`), e um teste falha se ela estiver desatualizada.
- Páginas novas entram na `nav` do `mkdocs.yml` no mesmo PR.
- O idioma da documentação é o português do Brasil.

## Consequências

- Nenhuma ferramenta extra é necessária para desenvolver; o MkDocs só roda na publicação.
- Versões fixadas em `docs/requirements.txt` (`mkdocs<2`, `mkdocs-material` 9.x): o time do
  Material alerta que o MkDocs 2.0 remove plugins e temas sem caminho de migração. Reavaliar
  (ex.: [Zensical](https://zensical.org/)) quando o MkDocs 1.x deixar de ser mantido.
- Validação local: `pip install -r docs/requirements.txt && mkdocs build --strict`.
- O workflow de publicação no GitHub Pages fica para a v1.0 (junto com CI).

## Alternativas consideradas

- **Hugo** — mais rápido, mas exige tema e front matter; MkDocs lê o Markdown atual sem mudanças.
- **Docusaurus** — traz Node.js para um projeto Go.
