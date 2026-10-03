# ctbz

CLI para a Contabilizei

ctbz opera a Contabilizei pela linha de comando: faz o mesmo login do site
(usuário, senha e código enviado por e-mail) e consulta os dados da empresa.

Dados vão para stdout no formato escolhido com -o; mensagens e progresso vão para stderr.

## Subcomandos

- [`ctbz api`](ctbz_api.md): Chama uma URL da plataforma com a sessão atual
- [`ctbz empresa`](ctbz_empresa.md): Mostra os dados da empresa selecionada
- [`ctbz login`](ctbz_login.md): Autentica na Contabilizei (usuário, senha e código por e-mail)
- [`ctbz logout`](ctbz_logout.md): Apaga a sessão local
- [`ctbz status`](ctbz_status.md): Mostra a sessão atual e testa se ainda é válida
- [`ctbz version`](ctbz_version.md): Mostra a versão do ctbz

## Flags

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
  -v, --version         mostra a versão do ctbz
```
