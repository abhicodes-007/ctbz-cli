# Changelog

Todas as mudanças relevantes deste projeto são documentadas aqui.

O formato segue o [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/) e o projeto
usa [Versionamento Semântico](https://semver.org/lang/pt-BR/).

## [Unreleased]

### Added

- `ctbz login`: autenticação com usuário/senha e código OTP enviado por e-mail, obtido por
  comando externo (`--otp-cmd`/`CTBZ_OTP_CMD`), pelo terminal ou em duas etapas (`--otp`).
- Seleção de empresa no login (`--cnpj`/`CTBZ_CNPJ`) e re-login automático quando a sessão expira.
- `ctbz status`, `ctbz empresa`, `ctbz api` e `ctbz logout`.
- Formatos de saída `table`, `json` e `csv` (`-o`/`--output`, `CTBZ_OUTPUT`).
- Script `scripts/otp-gmail-gws.sh` para ler o OTP do Gmail com o Google Workspace CLI,
  com mensagem clara (código 127) quando `gws` ou `jq` não estão instalados.
- Documentação da engenharia reversa, ADRs e roadmap.
- `ctbz version` (e `--version`), com versão, commit e data do build.
- Ajuda em português e referência de comandos gerada em `docs/referencia/`.
- Completion de shell: `ctbz completion bash|zsh|fish|powershell`.
- `ctbz empresas`: lista as empresas do usuário (inclusive inativas), marcando a atual.
- `ctbz empresa usar CNPJ`: troca a empresa da sessão refazendo o login.
- `ctbz empresa` mostra também natureza jurídica, data de abertura, início na Contabilizei,
  inscrição estadual e endereço.
- `ctbz empresa certificado`: situação, vencimento e dias para vencer do certificado digital.
- `ctbz empresa socios`: sócios com CPF, papel (administrador, responsável na Receita),
  categoria, salário-base e data de entrada.
- `ctbz empresa atividades`: CNAEs da empresa, a principal e os anexos do Simples Nacional.

### Changed

- Erros de uso (comando ou flag inexistente, argumento faltando, formato inválido) terminam
  com código de saída 2.

[Unreleased]: https://github.com/edusouza/ctbz-cli/commits/main
