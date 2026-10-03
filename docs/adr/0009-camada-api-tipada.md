# ADR-0009: Camada `internal/api` com tipos de resposta e funções de leitura

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

Os comandos começaram decodificando JSON direto em structs anônimos dentro de `internal/cli`.
Com dezenas de endpoints, isso mistura três coisas: o formato da API, a sessão (re-login,
cookies) e a apresentação. Também impede verificar os formatos da API de forma centralizada.

## Decisão

- `internal/api` guarda, por contexto (`empresa.go`, `impostos.go`…): as constantes de
  caminho (`PathXxx`), os tipos de resposta e funções `BuscarXxx(ctx, getter, …)`.
- `api.Getter` (`GetJSON(ctx, path, v)`) abstrai o transporte; `internal/cli` o implementa
  com a sessão salva e o re-login automático (`sessionGetter`).
- Os tipos declaram **só os campos que a CLI usa**. Campo novo na API não quebra nada.
- `api.Endpoints()` registra cada leitura (nome da fixture, caminho ao vivo, tipo) para os
  testes de contrato ([ADR-0010](0010-testes-de-contrato.md)).
- Nomes: tipos com o termo do domínio (`DadosEmpresa`, `Guia`), funções `Buscar…`.
- Conversões para a saída (`output.Record`/`List`) ficam em `internal/cli`, não em `internal/api`.

## Consequências

- `internal/api` é testável sem HTTP (um `Getter` falso que lê fixtures).
- A dependência é de cima para baixo: `cli → api → (nada)`; `cli → ctbz` para sessão.

## Alternativas consideradas

- **Métodos no `ctbz.Client`** — o cliente não conhece re-login (que depende de OTP e terminal),
  e os testes precisariam de um servidor HTTP para cada endpoint.
- **Gerar tipos a partir de respostas (quicktype)** — tipos com todos os campos viram contrato
  excessivo e quebrariam a cada campo removido que a CLI nem usa.
