# ctbz

CLI para a Contabilizei

ctbz opera a Contabilizei pela linha de comando: faz o mesmo login do site
(usuário, senha e código enviado por e-mail) e consulta os dados da empresa.

Dados vão para stdout no formato escolhido com -o; mensagens e progresso vão para stderr.

## Subcomandos

- [`ctbz api`](ctbz_api.md): Chama uma URL da plataforma com a sessão atual
- [`ctbz chamados`](ctbz_chamados.md): Lista os chamados de atendimento (em andamento ou finalizados)
- [`ctbz empresa`](ctbz_empresa.md): Mostra os dados da empresa selecionada
- [`ctbz empresas`](ctbz_empresas.md): Lista as empresas do usuário, marcando a atual
- [`ctbz impostos`](ctbz_impostos.md): Lista as guias de impostos a pagar (em atraso, do mês e do próximo mês)
- [`ctbz login`](ctbz_login.md): Autentica na Contabilizei (usuário, senha e código por e-mail)
- [`ctbz logout`](ctbz_logout.md): Apaga a sessão local
- [`ctbz mensalidade`](ctbz_mensalidade.md): Mostra a mensalidade atual da Contabilizei (valor, vencimento e situação)
- [`ctbz pendencias`](ctbz_pendencias.md): Lista as pendências da empresa (tipo, detalhe, prazo e situação)
- [`ctbz resumo`](ctbz_resumo.md): Mostra numa lista só o que precisa de atenção
- [`ctbz rotinas`](ctbz_rotinas.md): Lista as rotinas e obrigações do mês (da empresa e da Contabilizei)
- [`ctbz status`](ctbz_status.md): Mostra a sessão atual e testa se ainda é válida
- [`ctbz version`](ctbz_version.md): Mostra a versão do ctbz

## Flags

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
  -v, --version         mostra a versão do ctbz
```
