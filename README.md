# Guigoes

Guilherme de Castro's blog at [guigoes.com](https://guigoes.com), built with the published [Gazeta](https://github.com/araihu/gazeta) engine on [Margo](https://github.com/araihu/margo). `go.mod` pins both versions. This repository owns content, branding, configuration and deployment, while `cmd/blog` is a thin CLI adapter.

## Development

Use Go 1.27.3+ and Make. Node.js 24+ is needed only for deployment packaging.

```sh
npm ci
make check dry-run  # Build, validate schema, test and package without publishing
make run            # Preview at http://127.0.0.1:1313
```

`make install` downloads pinned Go dependencies. `public/` contains the generated publication; `.build/` contains temporary build files. Both are ignored by Git. Rebuild after editing. A failed build preserves the previous publication. `make rundocker` tests the static site through Nginx.

`make check` runs the Gazeta publication validator, Go tests without cached results, formatting checks and `go vet`. It also requires the sitemap to contain every indexable canonical page exactly once, excluding aliases and the error document. The independent `internal/schema` contract builds successive domains and author identities and reads their final JSON-LD. It covers article/profile identity, canonical URLs, covers, dates, collections and deliberate metadata corruption.

## Content and identity

Create `content/blogs/my-article/index.md` with TOML frontmatter:

```toml
+++
title = "Article title"
description = "Article summary."
language = "pt-BR"
authors = ["Guilherme de Castro"]
publishedAt = "2026-10-30T09:00:00-03:00"
tags = ["go", "kubernetes"]
draft = true
+++
```

The directory defines `/blogs/my-article/`. Keep article images under its `assets/` directory and use relative links. Review the article, provide the real publication time and set `draft` to `false` to publish. Drafts and future posts, including their assets, stay out of the publication and feeds. A future post needs another build after its publication time; there is no implicit schedule.

`site.yaml` controls domain, visual identity, icons, navigation, theme and the `DD/MM/YYYY` publication date format. `blog.yaml` declares author profiles, localized editorial labels, page types and legacy files/URLs. `content/` is the only source for posts and list ordering. Preserve `publishedAt` during revisions; use `modifiedAt` for relevant editorial changes.

Series use `series` and a positive, unique `seriesPart`. Gazeta generates the index and previous/next links from visible articles. Home and archives link to part 1 when appropriate, while publication dates retain their original meaning.

The editorial `.blog-*` classes belong to Gazeta. Guigoes supplies its palette, fonts and branding in `assets/blog.css`. Header, search, mobile navigation, theme controls, table of contents and viewers are native Margo components. The Homelab series enables `margo = { imageViewer = true }`; a standalone Markdown image with a title receives native zoom/pan controls. Covers and linked thumbnails remain static.

The branding adapts Margo's [editorial example](https://github.com/araihu/margo/tree/main/examples/blog), with its license preserved in `licenses/margo.txt`. The [media guide](docs/media-style.md) describes Guigoes artwork. `blog.yaml` explicitly enables reviewed HTML/JavaScript because existing articles contain demos, including the pinned `hyperscript.org@0.9.12` dependency. Publish only reviewed source.

## Structured data

Gazeta declares metadata through Margo's native configuration. Margo emits one safely serialized JSON-LD script in the document head. There is no generated-head rewriting or local schema renderer.

- Articles use `BlogPosting`, with source title, description, authors, dates and canonical identity.
- Guilherme de Castro's author URL is `/about/`, also linked visibly in article headers. Guests retain their own configured identities.
- A real `cover` becomes the article image and social preview, with decoded dimensions/MIME and `cover_alt`. Articles without covers omit JSON-LD `image`; the site-wide social preview remains a separate fallback.
- Home uses `WebSite` and `WebPage`; article lists use `ItemList` in visible order.
- Blog/tag archives and projects use `CollectionPage`; about uses `ProfilePage` with `Person`; contact uses `ContactPage`. The categories landing page uses `WebPage`.
- Alias HTML and 404 pages omit JSON-LD. The generated 404 document declares native `noindex, follow` metadata and stays out of the sitemap and `llms.txt`. No unsupported search action or breadcrumb contract is advertised.

Gazeta validates final HTML, links, fragments, feeds, authorship, dates, images and list contents before replacing `public/`. Margo's manifest describes the final artifact bytes. Use `go run ./cmd/blog check` to validate an existing artifact without rebuilding.

## CI/CD

GitHub Actions is the production publication pipeline. Pull requests run content/schema checks and a Wrangler dry-run, then validate the downloaded artifact in a separate job. No deployment secret is available during those steps.

On `main`, the deployment workflow builds and validates a publication, uploads `public/` and `wrangler.toml`, then downloads that exact artifact in a separate job. It validates the downloaded output without rebuilding, checks that the commit is still current `main`, and publishes Worker `guigoes-v2`. Production runs are serialized; an older commit cannot replace a newer publication.

The existing `wrangler.toml` retains the `guigoes.com` Custom Domain, Workers.dev access, preview URLs, directory URL handling and custom 404 behavior. Wrangler is pinned in `package-lock.json`. Only the real deploy step receives Cloudflare credentials.

Configure these repository settings:

| Type | Name | Purpose |
| --- | --- | --- |
| Variable | `CLOUDFLARE_ACCOUNT_ID` | Cloudflare account hosting `guigoes-v2` |
| Variable | `CLOUDFLARE_DEPLOY_ENABLED` | `true` only after Workers Builds stops publishing production |
| Secret | `CLOUDFLARE_API_TOKEN` | Workers deployment token restricted to the hosting account |

Before enabling Actions production deployment, replace the Workers Builds **Deploy default branch** deploy command with `npx --yes wrangler@4.156.0 versions upload`. This follows the [Cloudflare procedure for disabling automatic deployments](https://developers.cloudflare.com/workers/ci-cd/builds/#disconnecting-builds): builds may still save versions, but only Actions promotes production. Keep **Deploy non-production branches** for native branch previews. Two production publishers must not run concurrently. A migration can leave Actions publication disabled while PR checks and manual dry-runs are reviewed.

For a manual check, run **Deploy blog** with `dry_run = true` on `main`. This exercises build, artifact transfer and packaging without publishing. A real manual deployment also requires current `main`, configured credentials and the enable variable.

## Production verification

After a successful deployment, the workflow calls the published Gazeta verifier against `site.base_url`, using the downloaded artifact. It never rebuilds the expected site. Verification checks:

- matching local/public Margo manifest before and after the complete pass;
- publication content, canonical, locale, JSON-LD and social metadata;
- exact public asset bytes and suitable HTML/image/CSS/JavaScript MIME types;
- legacy HTTP redirects and a missing route returning the correct 404 page/status.

The check retries propagation failures within a three-minute deadline and fails the job if production does not match. `production-verification-<commit>` stores a JSON receipt with origin, manifest digest, counts and attempts. Known Cloudflare protection injection is ignored only in comparison trees; publication HTML is never rewritten.

Run `make verify` against an already-built local artifact to repeat this check. The [Rich Results Test](https://search.google.com/test/rich-results) remains an external validation tool; the CI verifier checks the publication contract rather than calling Google during every deployment.

## Rollback and repository language

Revert the source commit and push the revert to `main` to create another validated publication. For an emergency rollback, restore a previous Worker deployment through Cloudflare, then align Git before the next automatic deployment. Do not restore the retired Hugo build command.

`tests/legacy-urls.txt` protects existing public URLs and feeds. Legacy `/posts/.../` routes redirect to `/blogs/.../` with HTTP 301. Preserve those paths when editing content.

Write repository documentation, GoDoc, source identifiers, code comments and developer-facing messages in English. Published content, author identities, UI labels and intentional localization/Unicode fixtures keep their configured language, Brazilian Portuguese for Guigoes.
