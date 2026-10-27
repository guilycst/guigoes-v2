from datetime import datetime, timezone
import importlib.util
from pathlib import Path
import tempfile
import unittest
import xml.etree.ElementTree as ET

spec = importlib.util.spec_from_file_location("build", Path(__file__).parents[1] / "scripts/build.py")
build = importlib.util.module_from_spec(spec)
spec.loader.exec_module(build)


class PublicationTests(unittest.TestCase):
    def test_publication_details_preserve_dates_and_escape_text(self):
        output = build.publication_details({"authors": ["Ana & Rui"], "publishedAt": "2024-04-30T01:37:09-03:00", "modifiedAt": "2024-05-01T10:00:00Z", "tags": ["<SEO>"]})
        self.assertIn("Ana &amp; Rui", output)
        self.assertIn('datetime="2024-04-30T01:37:09-03:00">30/04/2024</time>', output)
        self.assertIn('datetime="2024-05-01T10:00:00Z">01/05/2024</time>', output)
        self.assertIn('href="/tags/seo/">&lt;SEO&gt;</a>', output)
        self.assertIn('data-margo-publication-metadata="true"', output)

    def test_editorial_listing_escapes_metadata(self):
        output = build.listing([{"title": '<script>alert("title")</script>', "route": "/blogs/test/", "publishedAt": "2024-04-30T01:37:09-03:00", "description": "Go & HTML", "cover": './assets/cover.webp', "cover_alt": 'A "cover"'}])
        self.assertNotIn("<script>", output)
        self.assertIn("&lt;script&gt;", output)
        self.assertIn("Go &amp; HTML", output)
        self.assertIn('src="blogs/test/assets/cover.webp"', output)
        self.assertIn('alt="A &quot;cover&quot;"', output)

    def test_drafts_and_future_posts_stay_private(self):
        now = datetime(2026, 10, 27, tzinfo=timezone.utc)
        self.assertFalse(build.published({"draft": True}, now))
        self.assertFalse(build.published({"publishedAt": "2026-10-27T00:00:00-03:00"}, now))
        self.assertTrue(build.published({"publishedAt": "2026-10-26T21:00:00-03:00"}, now))
        with self.assertRaises(ValueError):
            build.published({"publishedAt": "2026-10-27T00:00:00"}, now)

    def test_existing_tag_routes(self):
        self.assertEqual(build.tag_slug("SEO"), "seo")
        self.assertEqual(build.tag_slug("social media"), "social-media")
        self.assertEqual(build.tag_slug("Automação"), "automacao")
        with self.assertRaises(ValueError):
            build.tag_slug("///")

    def test_rss_escapes_content_and_preserves_timezone(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "index.xml"
            build.feed(path, [{"title": "Go & <HTML>", "route": "/blogs/test/", "publishedAt": "2024-04-30T01:37:09-03:00"}], "Guigoes")
            item = ET.parse(path).find("channel/item")
            self.assertEqual(item.findtext("title"), "Go & <HTML>")
            self.assertEqual(item.findtext("guid"), "https://guigoes.com/blogs/test/")
            self.assertTrue(item.findtext("pubDate").endswith("-0300"))
