# ADR-0013: Não repetir chamadas que falham sempre; usar uma fonte alternativa

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

`GET atendimento/chamados?finalizados=true` responde `400` com
`Memcache put: Item may not be more than 1048503 bytes in length` para contas com muitos
chamados: o servidor tenta guardar a lista inteira num cache limitado a 1 MB. O erro é
**determinístico** (a lista só cresce), então repetir a chamada, com ou sem espera, só
atrasa o mesmo erro e gera carga desnecessária. O próprio painel não mostra os finalizados
nessa situação.

`dadosempresa/get` devolve os 100 chamados mais recentes da empresa, inclusive os
finalizados, com id, assunto, status, data de criação e link.

## Decisão

- A CLI **não** faz retentativa automática de chamadas à API. Erros transitórios (rede,
  `5xx`) aparecem para o usuário, que pode rodar o comando de novo.
- Quando um endpoint falha de forma conhecida e existe outra fonte para os mesmos dados, o
  comando usa a fonte alternativa e avisa no stderr o que mudou (ex.: "mostrando os 100
  chamados mais recentes"). O código de saída continua 0.
- O fallback só é acionado por resposta HTTP de erro do servidor (`*ctbz.HTTPError`), nunca
  por sessão expirada (tratada pelo re-login) ou erro de rede.

## Consequências

- Comandos não ficam presos em esperas longas, e o comportamento é previsível em scripts.
- Cada fallback é específico do comando e precisa de teste com o erro simulado.
- Se a Contabilizei corrigir o endpoint, a fonte principal volta a ser usada sem mudança.

## Alternativas consideradas

- **Retentativa com backoff** (sugerida na issue #21) — inútil para um erro determinístico;
  só multiplicaria o tempo de resposta.
- **Falhar com mensagem clara** — correto, mas deixa o usuário sem os dados que existem em
  outro endpoint.
