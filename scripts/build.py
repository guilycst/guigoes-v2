#!/usr/bin/env python3
"""Prepare blog indexes and feeds; Margo owns Markdown and HTML rendering."""
from datetime import datetime, timezone
from email.utils import format_datetime
from html import escape
from pathlib import Path
import json
import posixpath
import re
import shutil
import subprocess
import tomllib
import unicodedata
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parents[1]
BASE = "https://guigoes.com"
STAGE = ROOT / ".build/content"
RENDER = ROOT / ".build/render"
NAV = '<nav class="guigoes-nav" aria-label="Navegação principal"><a class="guigoes-brand" href="/">guigoes<span aria-hidden="true">.</span></a><span class="guigoes-nav-links"><a href="/blogs/">Artigos</a><a href="/about/">Sobre</a><a href="/contact/">Contato</a><a href="/index.xml">RSS</a></span></nav>'
FOOTER = '<footer class="guigoes-footer"><span>Guilherme de Castro<br>Blumenau, Santa Catarina</span><span><a href="https://github.com/guilycst">GitHub</a> · <a href="https://www.linkedin.com/in/guicas">LinkedIn</a><br>Publicado com <a href="https://github.com/araihu/margo">Margo</a></span></footer>'


def read_post(path):
    _, front, body = path.read_text().split("+++", 2)
    return tomllib.loads(front), body.lstrip("\n")


def tag_slug(tag):
    text = unicodedata.normalize("NFKD", tag).encode("ascii", "ignore").decode().lower()
    slug = re.sub(r"[^a-z0-9]+", "-", text).strip("-")
    if not slug:
        raise ValueError(f"Tag without a URL slug: {tag!r}")
    return slug


def published(metadata, now):
    date = metadata.get("publishedAt")
    timestamp = datetime.fromisoformat(date) if date else None
    if timestamp is not None and timestamp.tzinfo is None:
        raise ValueError("publishedAt must include a timezone")
    return not metadata.get("draft", False) and (timestamp is None or timestamp <= now)


def write_page(path, metadata, body):
    target = STAGE / path
    target.parent.mkdir(parents=True, exist_ok=True)
    if body.startswith("# "):
        heading, _, body = body.partition("\n")
        heading += "\n\n"
    else:
        heading = f"# {metadata['title']}\n\n"
    if metadata.get("publishedAt"):
        body = publication_details(metadata) + "\n\n" + body
    page_class = "guigoes-home" if str(path) == "index.md" else "guigoes-content"
    target.write_text("---\n" + json.dumps(metadata, ensure_ascii=False) + "\n---\n\n"
                      + NAV + f'\n\n<div class="{page_class}">\n\n' + heading + body
                      + "\n\n</div>\n\n" + FOOTER + "\n")


def publication_details(metadata):
    authors = ", ".join(escape(author) for author in metadata.get("authors", []))
    dates = []
    for key, label in (("publishedAt", "Publicado"), ("modifiedAt", "Atualizado")):
        if metadata.get(key):
            value = metadata[key]
            date = datetime.fromisoformat(value).strftime("%d/%m/%Y")
            dates.append(f'{label} <time datetime="{escape(value, quote=True)}">{date}</time>')
    tags = "".join(f'<li><a href="/tags/{tag_slug(tag)}/">{escape(tag)}</a></li>' for tag in metadata.get("tags", []))
    # Margo's marker preserves this header while still generating head metadata.
    return ('<header class="margo-document__metadata" data-margo-publication-metadata="true" aria-label="Sobre este artigo">'
            f'<address>{authors}</address><span>{" · ".join(dates)}</span>'
            f'<ul aria-label="Tags">{tags}</ul></header>')


def listing(posts, source="index.md"):
    cards = []
    for post in posts:
        route = escape(post["route"], quote=True)
        cover = post.get("cover")
        image = ""
        if cover:
            asset = post["route"].lstrip("/") + cover.removeprefix("./")
            src = escape(posixpath.relpath(asset, posixpath.dirname(source) or "."), quote=True)
            alt = escape(post.get("cover_alt", ""), quote=True)
            image = f'<a class="guigoes-post-image" href="{route}" tabindex="-1" aria-hidden="true"><img src="{src}" alt="{alt}" loading="lazy" width="640" height="360"></a>'
        date = datetime.fromisoformat(post["publishedAt"]).strftime("%d/%m/%Y")
        cards.append(f'<li class="guigoes-post">{image}<div class="guigoes-post-text">'
                     f'<time datetime="{escape(post["publishedAt"], quote=True)}">{date}</time>'
                     f'<h2><a href="{route}">{escape(post["title"])}</a></h2>'
                     f'<p>{escape(post.get("description", ""))}</p></div></li>')
    return '<ol class="guigoes-posts" aria-label="Artigos">' + "".join(cards) + "</ol>\n"


def feed(path, posts, title):
    rss = ET.Element("rss", version="2.0")
    channel = ET.SubElement(rss, "channel")
    for key, value in {"title": title, "link": BASE + "/", "description": "Blog de Guilherme de Castro", "language": "pt-BR"}.items():
        ET.SubElement(channel, key).text = value
    for post in posts:
        item = ET.SubElement(channel, "item")
        for key, value in {"title": post["title"], "link": BASE + post["route"], "guid": BASE + post["route"], "description": post.get("description", ""), "pubDate": format_datetime(datetime.fromisoformat(post["publishedAt"]))}.items():
            ET.SubElement(item, key).text = value
    path.parent.mkdir(parents=True, exist_ok=True)
    ET.ElementTree(rss).write(path, encoding="utf-8", xml_declaration=True)


def build():
    # Only disposable build directories are removed before a successful render.
    for path in (STAGE, RENDER):
        if path.exists():
            shutil.rmtree(path)
    STAGE.mkdir(parents=True)
    posts = []
    documents = []
    now = datetime.now(timezone.utc)
    for source in sorted((ROOT / "content").rglob("*.md")):
        metadata, body = read_post(source)
        date = metadata.get("publishedAt")
        if not published(metadata, now):
            continue
        relative = source.relative_to(ROOT / "content")
        if relative.parts[0] == "blogs":
            if not date:
                raise ValueError(f"Missing publishedAt: {source}")
            post = dict(metadata, route="/" + relative.parent.as_posix() + "/")
            posts.append(post)
            cover = metadata.get("cover")
            if cover:
                body = f"![{metadata.get('cover_alt', metadata['title'])}]({cover})\n\n" + body
        documents.append((relative, metadata, body))
        for asset in (source.parent / "assets").rglob("*"):
            if asset.is_file():
                target = STAGE / asset.relative_to(ROOT / "content")
                target.parent.mkdir(parents=True, exist_ok=True)
                shutil.copy2(asset, target)
    posts.sort(key=lambda p: (datetime.fromisoformat(p["publishedAt"]), p["route"]), reverse=True)
    for relative, metadata, body in documents:
        if relative == Path("index.md"):
            body += "\n\n" + listing(posts[:5]) + '\n\n[Todos os artigos](/blogs/)\n'
        write_page(relative, metadata, body)
    page_meta = {"language": "pt-BR"}
    write_page("blogs/index.md", dict(page_meta, title="Blog"), listing(posts, "blogs/index.md"))
    tags = {}
    for post in posts:
        for tag in post.get("tags", []):
            tags.setdefault(tag_slug(tag), {"label": tag, "posts": []})["posts"].append(post)
    write_page("tags/index.md", dict(page_meta, title="Tags"), "\n".join(f"- [{v['label']}](/tags/{k}/)" for k, v in sorted(tags.items())))
    write_page("categories/index.md", dict(page_meta, title="Categorias"), "As publicações estão organizadas por [tags](/tags/).")
    for slug, tag in tags.items():
        write_page(f"tags/{slug}/index.md", dict(page_meta, title=tag["label"]), listing(tag["posts"], f"tags/{slug}/index.md"))
    write_page("404.md", dict(page_meta, title="Página não encontrada"), "Este endereço não existe. [Voltar ao blog](/blogs/).")
    # Published articles contain reviewed HTML/JS demos. Never use untrusted Markdown here.
    subprocess.run([str(ROOT / ".tools/margo/margo"), "site", "site.yaml", "--output-dir", str(RENDER), "--allow-unsafe-html"], cwd=ROOT, check=True)
    # Keep old image URLs even when an article no longer references its thumbnail.
    for asset in STAGE.rglob("*"):
        if asset.is_file() and "assets" in asset.relative_to(STAGE).parts:
            target = RENDER / asset.relative_to(STAGE)
            target.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(asset, target)
    shutil.copytree(ROOT / "static", RENDER, dirs_exist_ok=True)
    feed(RENDER / "index.xml", posts, "Guigoes")
    feed(RENDER / "blogs/index.xml", posts, "Guigoes · Blog")
    feed(RENDER / "tags/index.xml", [], "Guigoes · Tags")
    feed(RENDER / "categories/index.xml", [], "Guigoes · Categorias")
    for slug, tag in tags.items():
        feed(RENDER / f"tags/{slug}/index.xml", tag["posts"], "Guigoes · " + tag["label"])
    redirects = []
    for post in posts:
        for alias in post.get("aliases", []):
            if not re.fullmatch(r"/[a-z0-9_/-]+/", alias) or ".." in alias:
                raise ValueError(f"Unsafe alias: {alias}")
            target = RENDER / alias.strip("/") / "index.html"
            if target.exists():
                raise ValueError(f"Alias collides with published page: {alias}")
            target.parent.mkdir(parents=True, exist_ok=True)
            url = escape(BASE + post["route"], quote=True)
            target.write_text(f'<!doctype html><html lang="pt-BR"><head><meta charset="utf-8"><title>Redirecionando</title><link rel="canonical" href="{url}"><meta http-equiv="refresh" content="0; url={url}"></head><body><a href="{url}">Abrir artigo</a></body></html>')
            redirects.extend([f"{alias} {post['route']} 301", f"{alias.rstrip('/')} {post['route']} 301"])
    (RENDER / "_redirects").write_text("\n".join(redirects) + "\n")
    (RENDER / "robots.txt").write_text(f"User-agent: *\nAllow: /\nSitemap: {BASE}/sitemap.xml\n")
    public = ROOT / "public"
    if public.exists():
        shutil.rmtree(public)
    RENDER.rename(public)
    print(f"Built {len(posts)} posts with Margo v0.0.25 into public/")


if __name__ == "__main__":
    build()
