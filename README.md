# Guigoes

Blog de Guilherme de Castro, publicado em https://guigoes.com com
[Gazeta v0.0.2](https://github.com/araihu/gazeta/releases/tag/v0.0.2) sobre
[Margo commit 882c1cc](https://github.com/araihu/margo/commit/882c1cc46e9651e419c7836606ac3716efe5dedf).

## Desenvolvimento

Requisitos: Go 1.27.3+ e Make. `make install` baixa as dependências fixadas em
`go.mod`. O comando local delega à biblioteca Gazeta, que renderiza pelo Margo.
Não há engine editorial copiada neste repositório nem instalação de outro binário.
Não é necessário instalar Node ou Hugo para gerar o site.

```sh
make check  # Build, testes e verificação de URLs, imagens, links, fragmentos e RSS
make run    # Preview em http://127.0.0.1:1313; refaça make build após editar
```

`public/` contém a publicação. `.build/` contém o staging local e é ignorado pelo Git.
O build renderiza e valida em staging antes de substituir `public/`; erros
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

Gazeta, chamada por `go run ./cmd/blog build`, converte TOML para o frontmatter YAML aceito pelo Margo e gera
índices, tags, RSS e redirects. Margo renderiza Markdown, metadados semânticos,
canonicals, sitemap, imagens e código. `site.yaml` define domínio, identidade visual e tema.
`blog.yaml` configura perfis de autores, textos editoriais em português, tipos
de página e arquivos/URLs legados. Novos blogs usam a mesma Gazeta com seus próprios
arquivos de configuração, conteúdo e assets.
Os arquivos em `content/` são a única fonte dos posts e da ordenação das listas.
Posts de uma série usam `series = "Homelab K8s"` e `seriesPart = 1`, com número
positivo e único dentro da série. O build gera o índice com a parte atual e os
links anterior/próximo apenas para artigos visíveis. Home e arquivo oferecem uma
entrada pela parte 1 quando o post mais recente pertence à série e essa primeira
parte está disponível. A listagem continua ordenada por data, sem alterar datas
para representar a ordem de leitura.
`go run ./cmd/blog check` valida uma saída existente; `serve` inicia o preview local.

Cada post publicado recebe um único JSON-LD `BlogPosting`, gerado no build:

- `title` e `description` geram `headline` e `description`.
- Cada entrada de `authors` gera uma entidade `Person`. Guilherme de Castro usa
  `https://guigoes.com/about/` como perfil, também ligado no cabeçalho do artigo.
- `publishedAt` e, quando presente, `modifiedAt` preservam o timestamp e o fuso.
- `cover` gera uma URL absoluta em `image`. Sem capa, o campo é omitido.
- URL, identificador e `mainEntityOfPage` usam o endereço canônico do post.

Gazeta declara `structured_data` no frontmatter intermediário. O Margo emite um
único script JSON-LD no head, com escape seguro e rejeição de scripts duplicados.
Não há mais injeção local de scripts no corpo do artigo.

- Home: `WebSite` e `WebPage`, com a lista dos quatro posts exibidos.
- Arquivo e tags: `CollectionPage`; as listagens de posts incluem `ItemList` na mesma ordem visível.
- Projetos: `CollectionPage`, sem inventar preços, avaliações ou versões de software.
- Sobre: `ProfilePage` com `Person`, usando a mesma identidade dos autores dos posts.
- Contato: `ContactPage`. Categorias, que apenas aponta para tags: `WebPage`.
- Aliases e 404: sem JSON-LD. Rascunhos e posts futuros ficam fora do build público.

O validador confere o HTML final de cada tipo de página, canonical, autoria,
datas, capas, idioma e itens das listas. `SearchAction` e breadcrumbs não são
anunciados porque o blog não oferece esses contratos de navegação.

Referências: [Google Article](https://developers.google.com/search/docs/appearance/structured-data/article)
e [Schema.org BlogPosting](https://schema.org/BlogPosting).

As quatro capas seguem a direção cartoon do Guigoes. Consulte
[o guia de mídia](docs/media-style.md) para paleta, composição e prompts utilizados.

Os seletores editoriais `.blog-*` pertencem à Gazeta; paleta, fontes e marca
continuam em `assets/blog.css`.

O visual adapta o exemplo editorial [`examples/blog` do Margo](https://github.com/araihu/margo/tree/main/examples/blog),
com licença preservada em `licenses/margo.txt`. A home mostra os quatro posts mais
recentes, com destaque para o primeiro. `cover = "./assets/capa.webp"` e `cover_alt`
definem a capa do artigo e das listagens. O cabeçalho nativo formata datas em
`DD/MM/AAAA` via `publication.date_format`. O frontmatter intermediário declara
`publication_links` para perfis de autores e arquivos de tags, sem criar HTML próprio.

O cabeçalho, a busca, o menu móvel e a alternância de tema são nativos do
Margo. `site.yaml` declara os links, as duas versões do logo, os ícones e o webmanifest. `assets/blog.css`
fornece a paleta e os estilos editoriais pelos tokens públicos do tema.

`blog.yaml` habilita `unsafe_html: true` porque os posts existentes incluem demos
HTML/JavaScript. Publique apenas conteúdo revisado. A demo Hyperscript preserva
sua dependência externa fixa em `hyperscript.org@0.9.12`.

Capas sociais são declaradas em `social_image` e renderizadas pelo Margo com
URL absoluta, dimensões, MIME e texto alternativo. Gazeta valida o resultado
contra `cover` e `cover_alt`. Não há reescrita do head após renderizar.

Gazeta gera RSS e redirects, e reúne os arquivos estáticos do blog. Depois
de reuni-los aos artefatos do Margo, `site.MarshalManifest` descreve os bytes
finais, preservando rotas e identidade do build sem copiar o algoritmo de hash.

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
- Go: usa o runtime da imagem do Cloudflare; `go.mod` define a versão mínima.

Branches de trabalho geram versões de preview; apenas `main` publica produção.
`wrangler.toml` mantém assets em `public/`, URLs com barra final e página 404.
O domínio e a credencial de deploy continuam na integração existente do Cloudflare.
Não é necessário copiar tokens para o GitHub.

`tests/legacy-urls.txt` registra URLs e feeds da versão anterior para impedir
regressões. `/posts/.../` redireciona para `/blogs/.../` com HTTP 301 no Cloudflare.

Para rollback, restaure a versão anterior em Workers Deployments. Reverter o
código para antes da migração também exige restaurar o antigo build command
`hugo --minify` no Cloudflare.

Todas as páginas usam o layout `docs` do Margo com `sidebar: false` e
`pagination: false`. Apenas os posts habilitam `toc: true`. O sumário lateral
nativo respeita o cabeçalho e acompanha a seção em leitura em telas largas.
A navegação da série continua definida pelo blog; não há outra navegação,
controle de tema ou rodapé injetado no corpo do artigo.

Os artigos da série Homelab habilitam o visualizador nativo com
`margo = { imageViewer = true }` no frontmatter TOML. Imagens Markdown
isoladas com título recebem zoom, deslocamento, tamanho original e expansão.
Capas sem título e miniaturas com links permanecem estáticas. Para habilitar
o controle em uma imagem, use `![Descrição](./assets/imagem.png "Legenda")`.

The Margo dependency is pinned to an upstream commit for native aggregate asset
budgets and authored Mermaid accessible titles. `site.yaml` sets a 64 MiB source
asset budget while individual input limits remain unchanged.
