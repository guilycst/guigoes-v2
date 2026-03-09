# Guigoes (Hugo)

Este repositório foi migrado de uma aplicação Go + HTMX para um site estático com [Hugo](https://gohugo.io/) usando o tema [hugo-noir](https://github.com/prxshetty/hugo-noir).

## Stack atual

- Hugo (Extended)
- Tema: `themes/hugo-noir`
- Conteúdo em `content/`
- Dados do tema em `data/pt-br/`

## Estrutura principal

- `hugo.toml`: configuração principal do site
- `content/blogs/*/index.md`: posts migrados
- `content/blogs/*/assets/*`: assets dos posts
- `themes/hugo-noir`: tema aplicado

## Rodando localmente

Pré-requisito: Hugo Extended instalado.

```bash
make run
```

Ou diretamente:

```bash
hugo server -D --disableFastRender
```

## Build

```bash
make build
```

Gera o site estático em `public/`.

## Deploy

Deploy feito via **Cloudflare Pages** conectado ao GitHub.

Configuração recomendada no projeto Pages:

- Framework preset: `Hugo`
- Build command: `hugo --minify`
- Build output directory: `public`
- Root directory: `/`

Variáveis de ambiente:

- `HUGO_VERSION=0.150.0`
- `HUGO_ENV=production`

Depois de conectar a branch `main`, cada push gera um deploy automático.
