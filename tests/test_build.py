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
