# Roadmap

Funcionalidades planejadas para o `ctbz`, organizadas em versões. Cada versão cobre **um contexto completo** da
Contabilizei. Cada versão tem um épico no GitHub e uma sub-issue por funcionalidade.

Princípios:

- **Somente leitura até a 1.0.** Nenhuma versão inclui ações que alterem dados na Contabilizei (emitir nota, confirmar
  pagamento, alterar pró-labore, enviar documento…). Elas estão listadas em [Depois da 1.0](#depois-da-10).
- **Saída padronizada e testes de contrato desde a v0.1**: todo comando novo aceita `--output table|json|csv` e
  ganha um contrato que detecta mudanças no formato das respostas.
- Endpoints de cada funcionalidade vêm de [docs/api](docs/api/README.md). Os itens marcados como *investigação*
  dependem de descobrir parâmetros ainda não testados.

Labels usadas: `épico`, `versão: vX.Y`, `contexto: …`, `tipo: funcionalidade|investigação|infra|documentação`.

## Visão geral

| Versão | Contexto | Épico | Funcionalidades | Situação |
|---|---|---|---|---|
| **v0.1** | Login e dados da empresa | [#1](https://github.com/edusouza/ctbz-cli/issues/1) | 8 | em andamento (7/8) |
| **v0.2** | Impostos | [#10](https://github.com/edusouza/ctbz-cli/issues/10) | 6 | concluída |
| **v0.3** | Pendências, rotinas e atendimento | [#17](https://github.com/edusouza/ctbz-cli/issues/17) | 5 | concluída |
| **v0.4** | Mensalidade e pagamentos da Contabilizei | [#23](https://github.com/edusouza/ctbz-cli/issues/23) | 4 | planejada |
| **v0.5** | Notas fiscais | [#28](https://github.com/edusouza/ctbz-cli/issues/28) | 5 | planejada |
| **v0.6** | Pró-labore e distribuição de lucros | [#34](https://github.com/edusouza/ctbz-cli/issues/34) | 5 | planejada |
| **v0.7** | Contabilidade: relatórios, caixa e extratos | [#40](https://github.com/edusouza/ctbz-cli/issues/40) | 6 | planejada |
| **v0.8** | Documentos e certificado digital | [#47](https://github.com/edusouza/ctbz-cli/issues/47) | 2 | planejada |
| **v1.0** | Estabilidade e distribuição (futuro) | [#50](https://github.com/edusouza/ctbz-cli/issues/50) | 3 | planejada |

## v0.1 — Login e dados da empresa

Épico: [#1](https://github.com/edusouza/ctbz-cli/issues/1) · contexto `autenticação`

Base da CLI: autenticar na Contabilizei (usuário/senha + OTP por e-mail + seleção de empresa), manter a sessão e consultar os dados cadastrais da empresa. Inclui a fundação usada por todas as versões seguintes: formato de saída padronizado e testes de contrato contra a API.

| | Issue | Funcionalidade | Tipo |
|---|---|---|---|
| ✅ | [#2](https://github.com/edusouza/ctbz-cli/issues/2) | Login com usuário, senha e OTP por e-mail | funcionalidade |
| ✅ | [#3](https://github.com/edusouza/ctbz-cli/issues/3) | Comandos básicos: status, empresa, api e logout | funcionalidade |
| ✅ | [#4](https://github.com/edusouza/ctbz-cli/issues/4) | Documentação da engenharia reversa em docs/ | documentação |
| ✅ | [#5](https://github.com/edusouza/ctbz-cli/issues/5) | Saída padronizada: --output table\|json\|csv | infra |
| ✅ | [#6](https://github.com/edusouza/ctbz-cli/issues/6) | Testes de contrato da API | infra |
| ✅ | [#7](https://github.com/edusouza/ctbz-cli/issues/7) | Listar empresas do usuário e trocar de empresa | funcionalidade |
| ✅ | [#8](https://github.com/edusouza/ctbz-cli/issues/8) | Dados completos da empresa: sócios, endereço, atividades e certificado | funcionalidade |
| ⬜ | [#9](https://github.com/edusouza/ctbz-cli/issues/9) | Validar OTP automático com Gmail real (gws) | investigação |

Fora de escopo: CI e releases (v1.0); Ações que alteram dados.

## v0.2 — Impostos

Épico: [#10](https://github.com/edusouza/ctbz-cli/issues/10) · contexto `impostos`

Tudo sobre os impostos da empresa: guias a pagar (atrasadas, do mês e do próximo mês), detalhes e memória de cálculo, download das guias em PDF, histórico de pagamentos, faturamento usado na apuração e parcelamentos.

| | Issue | Funcionalidade | Tipo |
|---|---|---|---|
| ✅ | [#11](https://github.com/edusouza/ctbz-cli/issues/11) | Listar guias a pagar | funcionalidade |
| ✅ | [#12](https://github.com/edusouza/ctbz-cli/issues/12) | Detalhes de uma guia e memória de cálculo | funcionalidade |
| ✅ | [#13](https://github.com/edusouza/ctbz-cli/issues/13) | Baixar guias em PDF | funcionalidade |
| ✅ | [#14](https://github.com/edusouza/ctbz-cli/issues/14) | Histórico de impostos | funcionalidade |
| ✅ | [#15](https://github.com/edusouza/ctbz-cli/issues/15) | Faturamento mensal usado na apuração | funcionalidade |
| ✅ | [#16](https://github.com/edusouza/ctbz-cli/issues/16) | Parcelamentos e débitos federais (leitura) | investigação |

Fora de escopo: Confirmar pagamento de guia, recálculo e contratação de parcelamento (escrita).

## v0.3 — Pendências, rotinas e atendimento

Épico: [#17](https://github.com/edusouza/ctbz-cli/issues/17) · contexto `pendências`

O que a empresa precisa resolver: pendências abertas, rotinas e obrigações do mês, conciliações pendentes e chamados com o atendimento. Culmina num resumo de uma tela com tudo que exige atenção.

| | Issue | Funcionalidade | Tipo |
|---|---|---|---|
| ✅ | [#18](https://github.com/edusouza/ctbz-cli/issues/18) | Listar pendências da empresa | funcionalidade |
| ✅ | [#19](https://github.com/edusouza/ctbz-cli/issues/19) | Central de rotinas e obrigações do mês | funcionalidade |
| ✅ | [#20](https://github.com/edusouza/ctbz-cli/issues/20) | Pendências de conciliação fiscal | funcionalidade |
| ✅ | [#21](https://github.com/edusouza/ctbz-cli/issues/21) | Chamados de atendimento | funcionalidade |
| ✅ | [#22](https://github.com/edusouza/ctbz-cli/issues/22) | Resumo geral: o que precisa de atenção | funcionalidade |

Fora de escopo: Resolver pendência, aceitar termos, abrir chamado (escrita).

## v0.4 — Mensalidade e pagamentos da Contabilizei

Épico: [#23](https://github.com/edusouza/ctbz-cli/issues/23) · contexto `mensalidade`

A relação financeira com a própria Contabilizei: mensalidade atual, faturas, histórico de pagamentos, débito automático, plano e contrato.

| | Issue | Funcionalidade | Tipo |
|---|---|---|---|
| ⬜ | [#24](https://github.com/edusouza/ctbz-cli/issues/24) | Mensalidade e fatura atual | funcionalidade |
| ⬜ | [#25](https://github.com/edusouza/ctbz-cli/issues/25) | Histórico de pagamentos e débito automático | funcionalidade |
| ⬜ | [#26](https://github.com/edusouza/ctbz-cli/issues/26) | Situação de inadimplência | funcionalidade |
| ⬜ | [#27](https://github.com/edusouza/ctbz-cli/issues/27) | Plano e contrato de serviço | funcionalidade |

Fora de escopo: Pagar fatura, cadastrar cartão, trocar plano (escrita).

## v0.5 — Notas fiscais

Épico: [#28](https://github.com/edusouza/ctbz-cli/issues/28) · contexto `notas fiscais`

Consulta das notas fiscais de serviço emitidas, dos tomadores e da configuração do emissor, além das notas tomadas/importadas.

| | Issue | Funcionalidade | Tipo |
|---|---|---|---|
| ⬜ | [#29](https://github.com/edusouza/ctbz-cli/issues/29) | Listar notas fiscais emitidas | funcionalidade |
| ⬜ | [#30](https://github.com/edusouza/ctbz-cli/issues/30) | Baixar PDF e XML das notas | investigação |
| ⬜ | [#31](https://github.com/edusouza/ctbz-cli/issues/31) | Tomadores (clientes) | funcionalidade |
| ⬜ | [#32](https://github.com/edusouza/ctbz-cli/issues/32) | Configuração do emissor e alíquotas | funcionalidade |
| ⬜ | [#33](https://github.com/edusouza/ctbz-cli/issues/33) | Notas tomadas e notas de entrada | investigação |

Fora de escopo: Emitir, cancelar, replicar ou agendar notas; cadastrar tomador (escrita).

## v0.6 — Pró-labore e distribuição de lucros

Épico: [#34](https://github.com/edusouza/ctbz-cli/issues/34) · contexto `pró-labore`

Remuneração dos sócios: pró-labore atual e histórico, parâmetros (INSS, IRRF, teto), distribuição de lucros e informes de rendimentos.

| | Issue | Funcionalidade | Tipo |
|---|---|---|---|
| ⬜ | [#35](https://github.com/edusouza/ctbz-cli/issues/35) | Pró-labore atual e histórico | funcionalidade |
| ⬜ | [#36](https://github.com/edusouza/ctbz-cli/issues/36) | Parâmetros de cálculo do pró-labore | funcionalidade |
| ⬜ | [#37](https://github.com/edusouza/ctbz-cli/issues/37) | Distribuição de lucros | funcionalidade |
| ⬜ | [#38](https://github.com/edusouza/ctbz-cli/issues/38) | Informe e comprovante de rendimentos dos sócios | investigação |
| ⬜ | [#39](https://github.com/edusouza/ctbz-cli/issues/39) | Fator R e simulador de impostos (leitura) | funcionalidade |

Fora de escopo: Alterar ou zerar pró-labore, ativar gestão inteligente, registrar distribuição (escrita).

## v0.7 — Contabilidade: relatórios, caixa e extratos

Épico: [#40](https://github.com/edusouza/ctbz-cli/issues/40) · contexto `contabilidade`

Relatórios contábeis (balancete, balanço, razão), lançamentos do caixa, extratos bancários e plano de contas.

| | Issue | Funcionalidade | Tipo |
|---|---|---|---|
| ⬜ | [#41](https://github.com/edusouza/ctbz-cli/issues/41) | Balancete mensal | funcionalidade |
| ⬜ | [#42](https://github.com/edusouza/ctbz-cli/issues/42) | Balanço patrimonial | funcionalidade |
| ⬜ | [#43](https://github.com/edusouza/ctbz-cli/issues/43) | Razão contábil | funcionalidade |
| ⬜ | [#44](https://github.com/edusouza/ctbz-cli/issues/44) | Caixa: lançamentos do mês | funcionalidade |
| ⬜ | [#45](https://github.com/edusouza/ctbz-cli/issues/45) | Extratos e contas bancárias | funcionalidade |
| ⬜ | [#46](https://github.com/edusouza/ctbz-cli/issues/46) | Plano de contas e classificações | funcionalidade |

Fora de escopo: Importar extrato, classificar/desmembrar lançamentos, reabrir balanço (escrita).

## v0.8 — Documentos e certificado digital

Épico: [#47](https://github.com/edusouza/ctbz-cli/issues/47) · contexto `documentos`

Documentos enviados à Contabilizei e situação do certificado digital (candidata: pode ser reordenada ou removida).

| | Issue | Funcionalidade | Tipo |
|---|---|---|---|
| ⬜ | [#48](https://github.com/edusouza/ctbz-cli/issues/48) | Documentos enviados | funcionalidade |
| ⬜ | [#49](https://github.com/edusouza/ctbz-cli/issues/49) | Certificado digital: situação e renovação (leitura) | funcionalidade |

Fora de escopo: Enviar documentos, emitir/renovar/remover certificado (escrita).

## v1.0 — Estabilidade e distribuição (futuro)

Épico: [#50](https://github.com/edusouza/ctbz-cli/issues/50) · contexto `infra`

Itens transversais adiados: integração contínua, releases com binários e monitoramento de mudanças no site da Contabilizei.

| | Issue | Funcionalidade | Tipo |
|---|---|---|---|
| ⬜ | [#51](https://github.com/edusouza/ctbz-cli/issues/51) | CI: testes, vet e lint no GitHub Actions | infra |
| ⬜ | [#52](https://github.com/edusouza/ctbz-cli/issues/52) | Releases com binários | infra |
| ⬜ | [#53](https://github.com/edusouza/ctbz-cli/issues/53) | Monitorar mudanças no front e na API | infra |

Fora de escopo: Ações de escrita (avaliadas depois da 1.0).

## Depois da 1.0

Ações de escrita, a serem avaliadas depois que a leitura estiver estável. Todas exigiriam confirmação explícita (`--yes`)
e `--dry-run`:

- **Impostos:** confirmar pagamento de guia, solicitar recálculo, contratar parcelamento
- **Pendências:** resolver pendências, aceitar termos e cartas de responsabilidade, abrir chamado
- **Mensalidade:** pagar fatura, cadastrar cartão
- **Notas fiscais:** emitir, cancelar, replicar e agendar NFS-e; cadastrar tomadores
- **Pró-labore:** alterar ou zerar pró-labore, gestão inteligente, registrar distribuição de lucros
- **Contabilidade:** importar extratos, classificar e desmembrar lançamentos
- **Documentos e certificado:** enviar documentos, emitir/renovar certificado
