# ADR-0008: Cobra para a árvore de comandos

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

A CLI começou com o pacote `flag` da biblioteca padrão e um `switch` de comandos. O roadmap
prevê dezenas de comandos com subcomandos (`impostos guia`, `notas baixar`, `empresa usar`),
uma flag global (`-o`) aceita antes ou depois do comando e uma referência de comandos gerada
para o site. Com `flag`, cada um desses pontos exigiria código próprio (já havia um parser
manual só para o `-o` global).

## Decisão

- Usar [Cobra](https://github.com/spf13/cobra) (com pflag) para a árvore de comandos.
- Os comandos ficam em `internal/cli`, um arquivo por comando ou grupo; `cmd/ctbz/main.go` só
  chama `cli.Execute`. Assim `tools/gendocs` monta a mesma árvore para gerar a documentação.
- Ajuda e mensagens em português: modelo de uso traduzido, flags `-h` e `-v` redefinidas.
- Entrada e saída vêm do comando (`cmd.InOrStdin()`, `cmd.OutOrStdout()`, `cmd.ErrOrStderr()`),
  e os testes executam a árvore real com buffers.
- Variáveis de ambiente que servem de padrão para flags são lidas na execução, não no valor
  padrão da flag, para não vazarem para a ajuda e a documentação gerada.
- A referência em `docs/referencia/` é gerada por um gerador próprio de ~70 linhas, não pelo
  `cobra/doc`: queremos títulos em português e sem as dependências de man page e YAML.
- O comando `completion` do Cobra fica oculto da lista, mas disponível (`ctbz completion bash`).

## Consequências

- Subcomandos, flags herdadas, completion de shell e validação de argumentos sem código próprio.
- Dependência nova (cobra, pflag, mousetrap), todas amplamente usadas (gh, kubectl, hugo).
- Algumas mensagens internas do pflag continuam em inglês (ex.: "(default 3m0s)").

## Alternativas consideradas

- **Continuar com `flag`** — exigiria um dispatcher de subcomandos, parser de flags globais e
  gerador de ajuda próprios.
- **urfave/cli** — também adequado; Cobra tem mais adoção e geração de completion mais madura.
