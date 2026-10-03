# ADR-0016: Monitorar a API com um job agendado que abre issue

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

A CLI usa APIs internas e não documentadas do painel. Uma mudança da Contabilizei (campo
renomeado, endpoint removido) só seria notada quando alguém usasse o comando afetado. Os
testes de contrato (ADR-0010) já sabem detectar isso contra a API real, e o catálogo de
endpoints pode ser regenerado a partir do front, mas os dois exigem login com OTP.

## Decisão

- `scripts/monitorar.sh` junta as duas verificações (contratos ao vivo e catálogo
  regenerado comparado ao versionado, ignorando a data de geração) e devolve um relatório em
  Markdown e código 1 quando algo mudou.
- `.github/workflows/monitor.yml` roda o script toda segunda-feira e sob demanda, depois de
  um `ctbz login` com credenciais e `CTBZ_OTP_CMD` em segredos do repositório.
- Mudança detectada abre a issue "Monitoramento: a Contabilizei mudou a API ou o front" ou
  comenta nela, se já estiver aberta (uma issue só, sem duplicatas).
- Falha de login não abre issue: aparece como job vermelho no Actions.

## Consequências

- Mudanças viram issues em até uma semana, com o relatório pronto para agir.
- Exige configurar um `CTBZ_OTP_CMD` que funcione sem interação (webhook próprio ou `gws`
  autenticado no job); sem os segredos, o job falha com instrução.
- Campos novos que a CLI não usa não abrem issue: só o que quebra o contrato.

## Alternativas consideradas

- **Rodar os contratos ao vivo no CI de todo PR** — exigiria segredos em PRs e um OTP por
  execução.
- **Comparar o catálogo inteiro, com a data** — abriria issue toda semana.
