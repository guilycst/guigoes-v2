# Guigoes identity

`logo-v2.svg` is the outlined master. The symbol and `go` letters are cyan, `#007d9c`.
The other letters are charcoal, `#17131b`. `logo-light-v2.svg`
uses cream lettering for dark backgrounds. The standalone icon is cyan.
No SVG embeds bitmap images or depends on installed fonts.

Exports:

- Browser icons: SVG and 16, 32, 48, 96 px PNG; ICO contains 16, 32, 48 px.
- Apple touch icons: opaque cream background at 152, 167, 180 px.
- Android: 192 and 512 px, with separate opaque maskable versions.
- Share card: `social-v3.png`, 1280×640, used by Open Graph and X.
- Square media card: `social-square-v3.png`, 1080×1080, for manual sharing.
- Cards also have outlined SVG masters.

The icon masters use no pre-rounded background: devices apply their own mask.
`static/site.webmanifest` declares exact dimensions, MIME types and icon purposes.
The blog remains a normal website, with `display: browser`; no offline app is implied.

`cmd/blog/branding.go` adds device metadata to each generated head, copies assets,
and refreshes artifact hashes. Margo still supplies canonical and social metadata.
The visible page title and description remain specific to each route; the card
artwork is shared brand artwork. The square card is not used as a large-image card.

The artwork began with the built-in image generator, then was traced to vector
curves and recolored to isolate `go`. PNGs were exported from these vectors.
The `v3` card URLs avoid reusing cached previews with the previous dark symbol.
