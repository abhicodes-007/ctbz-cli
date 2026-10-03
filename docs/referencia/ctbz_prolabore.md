# ctbz prolabore

Mostra o pró-labore vigente por sócio e o tipo de gerenciamento

Mostra o pró-labore da empresa como a central de pró-labore do painel: tipo de
gerenciamento (ex.: INTELIGENTE, quando a Contabilizei calcula o valor ideal), total, valor e
competências do card do painel e, para cada sócio, valor, se recebe pró-labore, se é o
responsável na Receita, data da última atualização e quantidade de dependentes.

O histórico mensal está em "ctbz prolabore historico".

## Uso

```
ctbz prolabore
```

## Exemplos

```sh
  ctbz prolabore
  ctbz prolabore -o json | jq '.socios[] | {nome, valor}'
```

## Subcomandos

- [`ctbz prolabore historico`](ctbz_prolabore_historico.md): Lista o histórico mensal de pró-labore dos sócios
- [`ctbz prolabore parametros`](ctbz_prolabore_parametros.md): Mostra os valores usados no cálculo do pró-labore (INSS e IRRF)

## Flags globais

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
```

Veja também: [`ctbz`](ctbz.md).
