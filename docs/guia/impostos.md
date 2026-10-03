# Impostos

## Guias a pagar

```sh
ctbz impostos
```

```text
Grupo     ID                Imposto         Competência  Vencimento        Valor  Situação     Tipo
este_mes  1000000000000001  DARF Unificado  07/2026      06/10/2026  R$ 1.000,00  RECALCULADA  guia
Totais — Em atraso: R$ 0,00 (0) · Este mês: R$ 1.000,00 (1) · Próximo mês: R$ 0,00 (0)
```

- Grupos: `em_atraso`, `este_mes` e `proximo_mes`, como na tela de impostos do painel.
- `--atrasadas` mostra só as guias em atraso.
- Os totais por grupo saem no stderr (só na tabela), para não misturar com os dados.
- O `ID` é o que os outros comandos de impostos recebem.

### Alerta de atraso em scripts

`--fail-on-atraso` faz o comando terminar com **código 4** quando há guias em atraso:

```sh
# cron diário
ctbz impostos --fail-on-atraso -o csv > /dev/null || notify-send "Há impostos em atraso"
```

## Detalhe de uma guia

```sh
ctbz impostos guia 1000000000000001
```

Mostra valor total, valor original, juros e multa, situação, as ações disponíveis no painel
(ex.: `PAGAR`, `INFORMAR_PAGAMENTO`) e a explicação do imposto (o que é, frequência e o
impacto do atraso).

## Como o imposto foi calculado

```sh
ctbz impostos calculo
ctbz impostos tabela-irrf
```

`calculo` mostra a memória de cálculo do mês mais recente (a competência é escolhida pela
Contabilizei): faturamento, DAS do Simples Nacional, DARF de INSS e IRRF sobre o pró-labore,
faturamento e pró-labore dos últimos 12 meses e o percentual do Fator R. Valores ainda não
calculados aparecem vazios (`null` em JSON). `tabela-irrf` mostra as faixas do IRRF usadas no
cálculo do pró-labore.
