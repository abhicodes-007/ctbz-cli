#!/usr/bin/env python3
"""Gera o catálogo de endpoints da plataforma Contabilizei a partir do front.

Baixa o painel (/painel-de-controle/) e todos os chunks JavaScript usando a
sessão salva pelo `ctbz login`, encontra as chamadas axios e grava um Markdown
agrupado por base de API e por domínio.

Uso:
    ctbz login
    scripts/extrair-endpoints.py > docs/api/catalogo.md

Variáveis: CTBZ_HOME (diretório da sessão, padrão ~/.config/ctbz).
Só usa a biblioteca padrão do Python.
"""
import collections
import concurrent.futures
import datetime
import gzip
import json
import os
import re
import sys
import urllib.request

BASE = "https://app.contabilizei.com.br"
PAINEL = BASE + "/painel-de-controle/"
UA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0 Safari/537.36"

# Exportações do módulo que cria as instâncias axios (ver docs/frontend).
INSTANCIAS = {
    "c": "/api/plataforma/",
    "d": "/api/legado/",
    "b": "/api/legado/",
    "e": "/api/multiusuario/",
    "f": "/api/fintech/",
    "a": "/api/leads/hubspot/",
}

RE_CHAMADA = re.compile(r'\w+\["([a-f])"\]\.(get|post|put|delete|patch)\(\s*["`]([^"`]+)["`]')
# Chamadas cujo caminho está numa variável: var e="/novo-emissor/...";return u["c"].get(e)
RE_CHAMADA_VAR = re.compile(r'var (\w)="(/?[a-z][^"]+)";return \w+\["([a-f])"\]\.(get|post|put|delete|patch)\(\1\b')
RE_BASE_DECLARADA = re.compile(r'baseURL:"(/api/[^"]+)"')


def sessao_cookie():
    home = os.environ.get("CTBZ_HOME") or os.path.join(
        os.environ.get("XDG_CONFIG_HOME") or os.path.expanduser("~/.config"), "ctbz")
    with open(os.path.join(home, "session.json")) as f:
        sess = json.load(f)
    return "; ".join(f"{k}={v['Value']}" for k, v in sess["cookies"].items())


def baixar(url, cookie):
    req = urllib.request.Request(url, headers={"User-Agent": UA, "Cookie": cookie, "Accept-Encoding": "gzip"})
    with urllib.request.urlopen(req, timeout=60) as r:
        data = r.read()
        if r.headers.get("Content-Encoding") == "gzip":
            data = gzip.decompress(data)
    return data.decode("utf-8", errors="ignore")


def main():
    cookie = sessao_cookie()
    html = baixar(PAINEL, cookie)
    if "form-login" in html or 'location.replace("/login' in html:
        sys.exit("sessão expirada: rode `ctbz login`")
    scripts = sorted(set(re.findall(r'(?:src|href)="(js/[^"]+\.js)"', html)))
    with concurrent.futures.ThreadPoolExecutor(8) as ex:
        fontes = list(ex.map(lambda s: baixar(PAINEL + s, cookie), scripts))

    chamadas = collections.defaultdict(set)
    bases = set()
    for js in fontes:
        bases.update(RE_BASE_DECLARADA.findall(js))
        for inst, metodo, caminho in RE_CHAMADA.findall(js):
            chamadas[INSTANCIAS[inst]].add((metodo.upper(), caminho))
        for _, caminho, inst, metodo in RE_CHAMADA_VAR.findall(js):
            chamadas[INSTANCIAS[inst]].add((metodo.upper(), caminho))

    out = sys.stdout
    total = sum(len(v) for v in chamadas.values())
    out.write("# Catálogo de endpoints (gerado)\n\n")
    out.write(f"> Gerado por `scripts/extrair-endpoints.py` em {datetime.date.today().isoformat()} "
              f"a partir de {len(scripts)} arquivos JavaScript do painel. {total} chamadas encontradas.\n>\n"
              "> A base de cada chamada é inferida pela instância axios usada no código; "
              "caminhos terminados em `/` ou `=` recebem parâmetros concatenados pelo front.\n"
              "> Endpoints já testados estão em [endpoints-verificados.md](endpoints-verificados.md).\n\n")
    out.write("Bases declaradas no front: " + ", ".join(f"`{b}`" for b in sorted(bases)) + "\n\n")
    for base in sorted(chamadas, key=lambda b: (-len(chamadas[b]), b)):
        itens = chamadas[base]
        out.write(f"## `{base}` ({len(itens)})\n\n")
        por_dominio = collections.defaultdict(list)
        for metodo, caminho in itens:
            dominio = caminho.lstrip("/").split("/")[0].split("?")[0] or "(raiz)"
            por_dominio[dominio].append((metodo, caminho))
        for dominio in sorted(por_dominio):
            out.write(f"### {dominio}\n\n| Método | Caminho |\n|---|---|\n")
            for metodo, caminho in sorted(por_dominio[dominio], key=lambda t: (t[1], t[0])):
                out.write(f"| {metodo} | `{caminho}` |\n")
            out.write("\n")


if __name__ == "__main__":
    main()
