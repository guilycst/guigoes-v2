#!/usr/bin/env python3
"""Check published routes, local references, metadata, and feeds."""
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import unquote, urljoin, urlsplit
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parents[1]
PUBLIC = ROOT / "public"


class Page(HTMLParser):
    def __init__(self, text):
        super().__init__()
        self.refs = []
        self.ids = set()
        self.canonical = None
        self.language = None
        self.feed(text)

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if "id" in attrs:
            self.ids.add(attrs["id"])
        if tag == "html":
            self.language = attrs.get("lang")
        if tag == "link" and attrs.get("rel") == "canonical":
            self.canonical = attrs.get("href")
        for name in ("href", "src"):
            if attrs.get(name):
                self.refs.append(attrs[name])


def target(url):
    path = PUBLIC / unquote(urlsplit(url).path).lstrip("/")
    if path.is_dir():
        return path / "index.html"
    if not path.suffix and path.with_suffix(".html").is_file():
        return path.with_suffix(".html")
    return path


pages = {path: Page(path.read_text()) for path in PUBLIC.rglob("*.html")}
errors = []
for path, page in pages.items():
    route = "/" + path.relative_to(PUBLIC).as_posix().removesuffix("index.html")
    if page.language != "pt-BR":
        errors.append(f"{route}: missing pt-BR language")
    if not page.canonical or not page.canonical.startswith("https://guigoes.com/"):
        errors.append(f"{route}: invalid canonical")
    for ref in page.refs:
        url = urlsplit(urljoin("https://guigoes.com" + route, ref))
        if url.netloc != "guigoes.com" or url.scheme != "https":
            continue
        dest = target(url.geturl())
        if not dest.is_file():
            errors.append(f"{route}: missing {ref}")
        elif url.fragment and dest in pages and unquote(url.fragment) not in pages[dest].ids:
            errors.append(f"{route}: missing fragment {ref}")
for route in (ROOT / "tests/legacy-urls.txt").read_text().splitlines():
    if not target(route).is_file():
        errors.append(f"Legacy URL missing: {route}")
for path in PUBLIC.rglob("*.xml"):
    ET.parse(path)
posts = [p for p in pages if p.parent.parent.name == "blogs"]
items = ET.parse(PUBLIC / "index.xml").findall("channel/item")
if len(posts) != len(items):
    errors.append("RSS item count differs from published posts")
if errors:
    raise SystemExit("\n".join(errors))
print(f"Checked {len(pages)} HTML pages, local links/fragments, legacy URLs, and {len(items)} RSS items")
