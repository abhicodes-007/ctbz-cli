# Pró-labore e lucros

## Pró-labore vigente

```sh
ctbz prolabore
```

```text
Gerenciamento:  INTELIGENTE
Total:          R$ 1.621,00
Calculando:     não
Indisponível:   não
Sócios:
  Nome           CPF                   Valor  Recebe pró-labore  Responsável na Receita  Gestão       Atualizado em  Dependentes
  FULANO DE TAL  000.000.000-00  R$ 1.621,00  sim                sim                     INTELIGENTE  01/09/2026               2
```

- **Gerenciamento** `INTELIGENTE` significa que a Contabilizei calcula todo mês o pró-labore
  que leva ao menor imposto (inclusive o Fator R).
- O valor e as competências do card do painel (`valor_card`, `competencia_atual`,
  `competencia_anterior`) aparecem quando a Contabilizei já calculou o mês.
- Em JSON, os sócios vêm em `socios`; em CSV, a lista de sócios sai como JSON na célula.

## Histórico

```sh
ctbz prolabore historico               # todos os sócios, todo o histórico
ctbz prolabore historico --ano 2026
ctbz prolabore historico --socio ID    # ID de "ctbz prolabore -o json"
```

```text
Competência  Sócio          Pró-labore  Descontos
07/2026      FULANO DE TAL  R$ 1.621,00  R$ 178,31
```

A API devolve o histórico inteiro de um sócio de uma vez (não aceita ano nem página); a CLI
consulta cada sócio e filtra por `--ano`. Os descontos são o INSS e o IRRF retidos.
