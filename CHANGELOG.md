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
- Script `scripts/otp-gmail-gws.sh` para ler o OTP do Gmail com o Google Workspace CLI.
- Documentação da engenharia reversa, ADRs e roadmap.

[Unreleased]: https://github.com/edusouza/ctbz-cli/commits/main
