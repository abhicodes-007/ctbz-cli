#!/usr/bin/env sh
# Imprime o código OTP mais recente da Contabilizei lido do Gmail usando o
# Google Workspace CLI (gws). Feito para ser usado como CTBZ_OTP_CMD:
#
#   export CTBZ_OTP_CMD="$PWD/scripts/otp-gmail-gws.sh"
#   ctbz login
#
# A CLI define CTBZ_OTP_SINCE (epoch em segundos) com o instante em que o
# código foi pedido, então só e-mails a partir dali são considerados. Se nada
# chegou ainda, o script sai com status 1 e a CLI tenta de novo em alguns segundos.
#
# Variáveis opcionais:
#   CTBZ_OTP_QUERY  filtro de busca do Gmail (padrão: from:seguranca@contabilizei.com.br).
#                   Se você encaminha o e-mail para outra caixa, ajuste aqui,
#                   ex.: 'subject:"código de verificação" to:bot@seudominio.com'
#   GWS             caminho do binário gws (padrão: gws)
#
# Para testar sozinho (sem a CLI), rode o script logo depois de receber um código:
# ele procura e-mails dos últimos 10 minutos e imprime o código mais recente.
set -eu

GWS="${GWS:-gws}"
for bin in "$GWS" jq; do
	command -v "$bin" >/dev/null 2>&1 || {
		echo "comando não encontrado: $bin (veja docs/otp-automatico)" >&2
		exit 127
	}
done
SINCE="${CTBZ_OTP_SINCE:-$(($(date +%s) - 600))}"
QUERY="${CTBZ_OTP_QUERY:-from:seguranca@contabilizei.com.br} after:${SINCE}"

id=$("$GWS" gmail users messages list \
	--params "$(jq -nc --arg q "$QUERY" '{userId: "me", q: $q, maxResults: 1}')" |
	jq -r '.messages[0].id // empty')

if [ -z "$id" ]; then
	echo "nenhum e-mail de OTP encontrado ainda (busca: $QUERY)" >&2
	exit 1
fi

# O snippet normalmente já contém o código; o corpo decodificado é o plano B.
"$GWS" gmail users messages get \
	--params "$(jq -nc --arg id "$id" '{userId: "me", id: $id, format: "full"}')" |
	jq -r '
		.snippet,
		(.. | objects | select(.mimeType? == "text/plain" or .mimeType? == "text/html")
			| .body.data // empty | gsub("-"; "+") | gsub("_"; "/") | @base64d)
	' |
	grep -oE '\b[0-9]{6}\b' | head -n 1 | grep . || {
	echo "e-mail $id encontrado, mas sem código de 6 dígitos" >&2
	exit 1
}
