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

## Histórico de pagamentos e débito automático

```sh
ctbz mensalidade historico
```

```text
Débito automático ativo:       sim
Débito automático habilitado:  sim
Status:                        ATIVO
Competência:                   10/2026
Próxima cobrança:              10/11/2026
Pagamentos:
  Data        Competência        Valor  Status
  10/09/2026  09/2026        R$ 189,90  PAGO
  10/09/2026  08/2026      R$ 1.189,90  PAGO
```

- Lê `payments/recorrencia/init` e `payments/recorrencia/historico`.
- O formato de cada pagamento não está documentado e não foi verificado contra uma conta com
  histórico: os nomes dos campos (data, competência, valor, status) são inferidos, com
  alternativas aceitas. Se nenhum for reconhecido, o comando avisa no stderr e indica
  `ctbz api payments/recorrencia/historico`.
- Só campos conhecidos são copiados para a saída: chaves do gateway de pagamento e dados de
  cartão nunca aparecem.
- Para o pagamento recorrente de **impostos**, veja `ctbz impostos recorrente`.

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
