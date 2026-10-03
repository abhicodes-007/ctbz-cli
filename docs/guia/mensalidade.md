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

## Plano e contrato

```sh
ctbz plano                              # plano contratado
ctbz plano contrato > contrato.html     # contrato de prestação de serviços (HTML)
ctbz plano contrato --texto | less      # o mesmo, em texto simples
ctbz plano proposta --texto             # proposta do plano, com a tabela de preços
```

```text
Plano:               2025 - Simples - Serviço - Básico [139]
Categoria:           BASICO
Valor de tabela:     R$ 139,00
Ramos de atividade:  SERVICO
```

- `ctbz plano` usa os dados gravados no login (não chama a API). Depois de mudar de plano,
  faça o login de novo.
- O valor de tabela é o do plano; a mensalidade cobrada (com serviços adicionais) está em
  `ctbz mensalidade`.
- `contrato` e `proposta` imprimem o documento no stdout. Com `--texto`, a CLI descarta
  estilos e tags, quebra linha por parágrafo e separa as células das tabelas.
