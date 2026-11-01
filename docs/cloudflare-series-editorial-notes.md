# Cloudflare series, parts 2 through 5

## Scope

On November 1, 2026, the author requested sequential drafting of parts 2 through 5 with independent review of each part before advancing. These are conceptual Brazilian Portuguese articles, without installation tutorials. Each has approximately 1,200 body words, a distinct original cover, native series metadata, and `draft = true`. Part 1 was preserved unchanged.

The new articles cover:

- Part 2: domains, registrar versus DNS/hosting, proxy, HTTPS, cache and DDNS.
- Part 3: Tunnel behavior, protocol/client distinctions, workloads, public-media terms, upload/timeouts and Quick Tunnels.
- Part 4: Access, private routes, Gateway, shared Zero Trust seats, support and log-retention limits.
- Part 5: Workers, Static Assets, R2 allowances and billing, Turnstile and Email Routing/Sending.

No live Cloudflare resource, DNS record, access policy, subscription or deployment was changed. Other chats' exposure and Tailscale drafts remain untouched.

## Independent review

Two fresh reviewer tasks, `/root/parts_editorial` and `/root/parts_technical`, used the required `gpt-5.6-sol`/high reviewer profile without inherited producer history. They reviewed the Git worktree directly, scoped to each article, with before/after snapshots. They were prohibited from reading producer notes, receipts, transcripts, session history and memory. Each declared stable Git identities and read-only isolation. Both applied prose specialist criteria; changing technical claims were checked against primary Cloudflare sources.

Both reviewers accepted every final article with no unresolved blockers or deferred hardening. Exact base/HEAD and accepted scoped identities are recorded in [review receipts](cloudflare-series-review-receipts.json).

Part 2 passed its first review; a second pass accepted an optional accessibility improvement to the cover alt, identifying the Cloudflare-inspired character as editorial art and describing the tag's waving gesture. Parts 3 and 4 were accepted in one pass. Part 3's optional suggestion to simplify one terms sentence was not adopted because its surrounding examples already express the practical boundary without expanding the legal discussion. Part 5's technical review requested the current Beta designation for Email Sending. Both reviewers accepted its second pass after that change.

## Images

Four covers were generated with the integrated image-generation tool, one image per call, using the existing character sheets as attached visual references. The prompt set and reference identities are preserved in [cover prompts](cloudflare-series-cover-prompts.json). Each accepted output was copied into its article's `assets/cover-cartoon.png`; the originals remain in the task's generated-image directory and external visual-evidence directory.

All covers were opened at full size. Character silhouettes, faces, visible limbs, key shaft, tube continuity, contact points and storage-box geometry were inspected. The server's far-side arm and one cloud foot have natural occlusion in their scenes. No text, watermark, extra key leg or detached limb was observed. The artwork is an original editorial interpretation, not an official Cloudflare mascot.

Parts 2, 3 and 4 reuse already-redacted DNS, tunnel-list and Access-policy screenshots respectively. No new authenticated screenshots or account changes were needed.

## Validation

`make check` passed: normal publication validation, uncached Go tests, formatting and `go vet`. All five Cloudflare source articles remain drafts, excluded from the normal publication and RSS. The normal build contains 51 HTML pages and 9 BlogPosting blocks.

The isolated preview marks only its copies of the five series articles visible. Its final build validated 61 HTML pages and 14 BlogPosting blocks, local links/fragments, JSON-LD, images, feeds and legacy URLs. A separate assertion verified the complete five-part native series index in every new article and unchanged accepted review identities.

Browser inspection covered the new article covers and native series links, desktop dark and 390px mobile light layouts, and blog/home listing thumbnails. The temporary viewport and theme were restored and the separate QA tab was closed. No real-device touch or live deployment verification was performed. The preview remains available at `http://127.0.0.1:1314/` for author iteration.

## Subsequent complete-series review

The author requested an additional independent review of all five parts together. Three fresh reviewers accepted the corrected series after two passes. The final [cross-part review](cloudflare-series-final-review.md) records progression, terminology, plan/allowance boundaries, reader comprehension and technical consistency. It supersedes the earlier authoring-stage version coverage for part 5: the memory row was removed from its introductory table, and daily-quota behavior was qualified for origin routes. All article sources remain drafts.
