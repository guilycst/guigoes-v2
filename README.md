# Guigoes

Blog de Guilherme de Castro, publicado em https://guigoes.com com
[Margo v0.0.25](https://github.com/araihu/margo/releases/tag/v0.0.25).

## Desenvolvimento

Requisitos: Python 3.11+, Make e curl. `make install` baixa o binário oficial
do Margo para Linux ou macOS e verifica SHA-256 antes de instalá-lo em `.tools/`.
Não é necessário instalar Go, Node ou Hugo para gerar o site.

```sh
make check  # Build, testes e verificação de URLs, imagens, links, fragmentos e RSS
make run    # Preview em http://127.0.0.1:1313; refaça make build após editar
```

`public/` contém a publicação. `.build/` e `.tools/` são locais e ignorados pelo Git.
O build renderiza em staging antes de substituir `public/`; erros de renderização
preservam a saída anterior. `make rundocker` testa a imagem estática com Nginx.

## Publicar um post

Crie `content/blogs/meu-post/index.md` com metadados TOML entre `+++`:

```markdown
+++
title = "Título do post"
description = "Resumo do artigo."
language = "pt-BR"
authors = ["Guilherme de Castro"]
publishedAt = "2026-10-27T09:00:00-03:00"
tags = ["go", "kubernetes"]
draft = false
+++
Texto do artigo em Markdown.
```

Imagens ficam em `content/blogs/meu-post/assets/` e usam links `./assets/imagem.png`.
A pasta define a URL `/blogs/meu-post/`. `draft = true` e datas futuras excluem
o post de páginas, assets, listagens, sitemap e feeds. Datas futuras passam a
aparecer no próximo build após o horário informado, sem agendamento automático.

`scripts/build.py` converte TOML para o frontmatter YAML aceito pelo Margo e gera
índices, tags, RSS e redirects. Margo renderiza Markdown, metadados semânticos,
canonicals, sitemap, imagens e código. `site.yaml` define identidade e tema.
Os arquivos em `content/` são a única fonte dos posts e da ordenação das listas.

O visual adapta o exemplo editorial [`examples/blog` do Margo](https://github.com/araihu/margo/tree/main/examples/blog),
com licença preservada em `licenses/margo.txt`. A home mostra os cinco posts mais
recentes, com destaque para o primeiro. `cover = "./assets/capa.webp"` e `cover_alt`
definem a capa do artigo e das listagens. O cabeçalho próprio formata datas em
`DD/MM/AAAA` e liga tags aos arquivos, mantendo os metadados semânticos do Margo.

O build habilita `--allow-unsafe-html` porque os posts existentes incluem demos
HTML/JavaScript. Publique apenas conteúdo revisado. A demo Hyperscript preserva
sua dependência externa fixa em `hyperscript.org@0.9.12`.

## CI/CD

GitHub Actions executa `make check` em pull requests e pushes em `main`, e guarda
`public/` como artifact. Cloudflare Workers Builds usa a conexão GitHub existente
para executar os mesmos checks antes de publicar o Worker `guigoes-v2`.

Configuração de Workers Builds em Settings > Builds:

- Repositório: `guilycst/guigoes-v2`, branch de produção `main`.
- Build command: `make check`.
- Deploy command: `npx --yes wrangler@4.153.0 deploy`.
- Version command: `npx --yes wrangler@4.153.0 versions upload`.
- Root directory: `/`.
- Python: `3.14.8`, já disponível na imagem do Cloudflare e selecionado por `.python-version`.

Branches de trabalho geram versões de preview; apenas `main` publica produção.
`wrangler.toml` mantém assets em `public/`, URLs com barra final e página 404.
O domínio e a credencial de deploy continuam na integração existente do Cloudflare.
Não é necessário copiar tokens para o GitHub.

`tests/legacy-urls.txt` registra URLs e feeds da versão anterior para impedir
regressões. `/posts/.../` redireciona para `/blogs/.../` com HTTP 301 no Cloudflare.

Para rollback, restaure a versão anterior em Workers Deployments. Reverter o
código para antes da migração também exige restaurar o antigo build command
`hugo --minify` no Cloudflare.
