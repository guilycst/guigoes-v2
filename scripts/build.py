#!/usr/bin/env python3
"""Prepare blog indexes and feeds; Margo owns Markdown and HTML rendering."""
from datetime import datetime, timezone
from email.utils import format_datetime
from html import escape
from pathlib import Path
import json
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
NAV = '<nav class="guigoes-nav" aria-label="Navegação principal"><a href="/">Guigoes</a><a href="/blogs/">Blog</a><a href="/about/">Sobre</a><a href="/contact/">Contato</a><a href="/index.xml">RSS</a></nav>'
FOOTER = '<footer class="guigoes-footer">Guilherme de Castro · Publicado com <a href="https://github.com/araihu/margo">Margo</a></footer>'


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
    heading = "" if body.startswith("# ") else f"# {metadata['title']}\n\n"
    target.write_text("---\n" + json.dumps(metadata, ensure_ascii=False) + "\n---\n\n"
                      + NAV + "\n\n" + heading + body + "\n\n" + FOOTER + "\n")


def listing(posts):
    return "\n\n".join(
        f"## [{p['title']}]({p['route']})\n\n"
        f"{datetime.fromisoformat(p['publishedAt']).strftime('%d/%m/%Y')}\n\n{p.get('description', '')}"
        for p in posts
    ) + "\n"


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
            tags = " · ".join(f"[{t}](/tags/{tag_slug(t)}/)" for t in metadata.get("tags", []))
            body = f"Explorar tags: {tags}\n\n" + body
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
            body += "\n\n## Publicações\n\n" + listing(posts)
        write_page(relative, metadata, body)
    page_meta = {"language": "pt-BR"}
    write_page("blogs/index.md", dict(page_meta, title="Blog"), listing(posts))
    tags = {}
    for post in posts:
        for tag in post.get("tags", []):
            tags.setdefault(tag_slug(tag), {"label": tag, "posts": []})["posts"].append(post)
    write_page("tags/index.md", dict(page_meta, title="Tags"), "\n".join(f"- [{v['label']}](/tags/{k}/)" for k, v in sorted(tags.items())))
    write_page("categories/index.md", dict(page_meta, title="Categorias"), "As publicações estão organizadas por [tags](/tags/).")
    for slug, tag in tags.items():
        write_page(f"tags/{slug}/index.md", dict(page_meta, title=tag["label"]), listing(tag["posts"]))
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
