# Mensalidade e pagamentos

## Mensalidade atual

```sh
ctbz mensalidade
```

```text
Competência:                     10/2026
Valor:                           R$ 15,90
Situação:                        Em aberto
Observação:                      + 1 Serviço adicional
Competência anterior em atraso:  não
```

- Os dados vêm do card de mensalidade do painel (`dashboard/fatura`) e do indicador de
  atraso (`dashboard/v1/mensalidade`).
- O vencimento aparece quando a fatura já foi gerada.
- `--fail-on-atraso` termina com **código 4** quando há mensalidade de competência anterior
  em atraso.

## Situação com a Contabilizei

```sh
ctbz mensalidade situacao
```

```text
Em dia:    sim
Situação:  OK
```

- Usa a mesma consulta de inadimplência do painel, que responde texto puro: `OK` quando a
  empresa está em dia. Outras respostas não foram vistas; qualquer valor diferente de `OK` é
  mostrado como veio, com `em_dia` falso.
- `--fail-on-inadimplencia` termina com **código 4** quando a empresa não está em dia.
