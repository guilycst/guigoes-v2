# Composed navigation verification

The blog consumes published Go module revisions. No local replacement, consumer
CSS, consumer JavaScript or generated-HTML rewrite implements this behavior.

- App Shells adds typed global destinations and composes them with optional
  contextual documentation entries in one mobile drawer. Its existing
  HeaderNavigation slot remains an opt-in render-once alternative.
- Margo maps navigation.links to that contract. sidebar:false disables contextual
  entries, while global destinations remain reachable before the brand through
  the mobile hamburger. Language controls retain their configured labels/flags.
- Existing component browser checks and new tests cover 390, 1023, 1024 and
  1440 px, Escape, focus return/trapping, JavaScript-disabled native links,
  HTMX route changes and Back/Forward. Margo site tests also exercise local and
  inline assets, publication metadata and pages that toggle contextual navigation.
  History restoration covers body swaps and ignores late metadata for a route
  that is no longer visible.

Upstream work:

- [App Shells PR 61](https://github.com/araihu/goshtoso-app-shells/pull/61)
- [Margo PR 115](https://github.com/araihu/margo/pull/115)

The local blog build, URL/JSON-LD validation, tests and vet pass with GOWORK=off.
Physical touch devices and assistive-technology certification are outside this
check. The production blog is unchanged; the existing article PR remains a draft.
