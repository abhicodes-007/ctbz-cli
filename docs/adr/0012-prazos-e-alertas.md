# ADR-0012: Calcular alertas de prazo na CLI e sinalizá-los com o código 4

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

Pendências, rotinas e obrigações têm prazo, e a pergunta mais comum é "o que está vencido
ou vence logo?". As APIs devolvem só a data (e às vezes uma situação como `PENDENTE` ou
`EM_ABERTO`), sem dizer se o prazo já passou. Cada comando decidindo isso do seu jeito daria
resultados diferentes para a mesma data.

## Decisão

- Comandos com prazo trazem a coluna `alerta`, calculada pela CLI só para itens **abertos**
  (qualquer situação que não seja de concluído):
  - `vencida`: o prazo é anterior a hoje;
  - `próxima`: o prazo é hoje ou nos próximos **7 dias**;
  - vazio (`null` no JSON): sem alerta, item concluído ou sem prazo.
- "Hoje" é a data local da máquina (`now()` em `internal/cli`, trocável nos testes).
- Por padrão, listas de pendências mostram só itens abertos; `--todas` inclui os concluídos.
- `--fail-on-vencidas` faz o comando terminar com código 4 (atenção) quando há item
  `vencida`, como `--fail-on-atraso` em `ctbz impostos`.

## Consequências

- Scripts filtram por `alerta` sem reimplementar a regra de datas.
- A janela de 7 dias é fixa; se for preciso outra, vira flag sem quebrar o contrato.
- Os valores de `alerta` fazem parte do contrato de saída (ADR-0006): renomeá-los é
  mudança incompatível.

## Alternativas consideradas

- **Mostrar só a data** — empurra a conta para cada usuário e cada script.
- **Colorir a tabela** — não funciona em JSON/CSV nem em terminais sem cor, e não ajuda scripts.
- **Dias até o prazo como número** — útil, mas menos direto que uma etiqueta; pode ser
  adicionado depois sem quebrar nada.
