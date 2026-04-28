# PentAGI — Web assessment (compact) — OWASP Juice Shop @ localhost

Use this template for **grey-box / lab** testing when the app is already running (e.g. Docker) at **`http://localhost:3000`**. The UI is an Angular SPA; the shell route is **`http://localhost:3000/#/`** — fragments are **not** sent to the server; keep `base_url` as **origin only** and put SPA notes in `notes`.

---

## A) COPY-PASTE — paste everything below into the PentAGI task (single block)

```text
=== WEB_ASSESSMENT_INPUT ===
execution_profile: BALANCED
auth_required: no
base_url: http://localhost:3000
notes: OWASP Juice Shop lab; UI http://localhost:3000/#/ (Angular SPA); APIs under /rest/ and /api/; hash fragment is not sent to server; use README/SOLUTIONS lab accounts only if auth branch needed—no brute force.
dast_allowed: yes
dast_primary_tool: httpx
dast_secondary_tool: nuclei
in_scope_hosts: localhost,127.0.0.1
out_of_scope_hosts:
max_requests_per_minute: 120
max_parallel_requests: 3
upgrade_on_signal: true
roe_environment: lab
roe_exploit_level: poc_read_only
roe_emergency_contact: N/A
=== END WEB_ASSESSMENT_INPUT ===

=== ROE_INPUT ===
roe_environment: lab
roe_exploit_level: poc_read_only
roe_emergency_contact: N/A
allowed_methods: GET,POST,PUT,PATCH,DELETE,OPTIONS
disallowed_actions: destructive_data_change,credential_stuffing,out_of_scope_hosts,exfiltration_to_third_party
=== END ROE_INPUT ===
```

Optional — **only** if you already have a valid session after manual login in a real browser (paste at runtime; never commit tokens to git):

```text
=== BROWSER_SESSION ===
access_token: <PASTE_JWT_OR_LEAVE_EMPTY>
refresh_token:
expires_at:
token_type: bearer
=== END BROWSER_SESSION ===
```

If you use `BROWSER_SESSION`, set in the first block: `auth_required: yes` and complete token fields.

---

## B) Juice Shop–specific hints (do not echo this file into the report)

- **Stack:** Express backend + Angular frontend; many challenges tie to **`/rest/*`**, uploads, auth, and static files under `/ftp`, `/assets`, etc.
- **Scope:** only **`localhost`** / **`127.0.0.1`** on the port you mapped (default **3000**). If the container uses another host port, change `base_url` and retest.
- **DAST:** one primary pass (**httpx** fingerprint + **nuclei** with **tagged** templates matching Node/Express/exposure—not full chaos templates). Respect `max_requests_per_minute` / `max_parallel_requests`.
- **Auth:** default run is **`auth_required: no`** (anonymous + public API). Upgrade to authenticated branches only with explicit RoE + session/JWT from operator.
- **Stability:** tool-call text stays **short, English, ASCII-heavy**; if `MALFORMED_FUNCTION_CALL` repeats twice, stop the branch and use `ask`/`done` with a one-line blocker.

---

## C) Execution profile (choose one; must match block in §A)

| Profile | When to use |
|--------|-------------|
| `FAST` | Quick smoke + top P0 only; one DAST tool; tight time. |
| `BALANCED` | Default lab run (above). |
| `MAX` | Deep coverage; more P1 and second DAST pass if RoE allows. |

Hard limits (mandatory):

| Item | FAST | BALANCED | MAX |
|------|-----:|---------:|----:|
| P0 flows deep-dived (Pass B) | 5 | 10 | 18 |
| Deep-read routes/APIs | 8 | 16 | 30 |
| Payload variants per branch | 2 | 3 | 5 |
| Memory queries (whole run) | 2 | 3 | 6 |
| Global DAST scanner runs | 1 | 1 | 2 |
| DAST wall-clock (minutes) | 8 | 15 | 30 |
| Retries same error | 2 | 3 | 3 |
| Parallel branches | 1 | 1 | 2 |

---

## D) Pipeline order (reduce loops)

1. One fingerprint: tech, auth style, main API prefixes, cookies.
2. One inventory pass (httpx / light crawl / manual `rest` discovery) — dedupe URLs.
3. Pass A map P0/P1/P2; Pass B only P0 first, then P1 if time.
4. `search` (OSINT): **one** pass after fingerprint (framework + known Juice Shop patterns), not per endpoint.
5. `browser` only when DOM/JS is needed; else `terminal` (curl + jq) with rate caps.
6. `coder` only for >5 repetitive structured checks; never for a single curl.
7. No parallel subtasks for the same goal.

### D.1) Speed + token optimization (mandatory)

- **One inventory source of truth:** write normalized URL list once (file under `/work` if available); reuse it—do not re-crawl the same paths for every branch.
- **DAST:** run **httpx** first (passive fingerprint only). Run **nuclei** only with **narrow tags** (e.g. `technologies`, `exposures`, `misconfiguration`, `panel`—match observed stack); **never** full `-t all` / community mega-lists on localhost lab.
- **Wall-clock:** stop each scanner at the profile cap in §C (`DAST wall-clock`); record `Untested` + reason instead of extending.
- **Requests:** honor `max_requests_per_minute` and `max_parallel_requests`; on 429/slowdown, halve parallelism and continue P0 only.
- **Tool calls:** keep each tool `question`/`message` **under ~1.2KB**; put long payloads in files and reference paths—reduces `MALFORMED_FUNCTION_CALL` and log bloat.
- **Per-branch cap:** max **2** expensive tool attempts (same error or no new signal); then `ask`/`done` or pivot—no blind retries.
- **Report:** dedupe findings by root cause; one representative endpoint + blast radius note—no duplicate prose for similar routes.

---

## E) Safety & truthfulness

- No finding without **tool-backed** evidence (request/response snippet, log, or reproducible step).
- Any redirect or call off `in_scope_hosts` → stop branch, record residual risk.
- Redact secrets in evidence; never paste real production credentials.

---

## F) Report output contract (exact heading order)

1. `# Executive Summary`
2. `## One-page Risk Snapshot`
3. `## Engagement Metadata` (must include: profile, DAST tool(s) actually run, scope hosts, rate limits)
4. `## Attack Surface Inventory`
5. `## Findings`
6. `## Phase 3 Coverage Matrix`
7. `## Residual Risks & Assumptions`
8. `## Kill Chains (if any)`
9. `## Remediation Roadmap (30/60/90 days)`
10. `## Retest Scope`

Missing data: `N/A - <reason>`.

Attack Surface Inventory columns:

`ID | Method | Path/Flow | Auth Type | Role | Source (crawl/js/wordlist/manual) | Priority (P0/P1/P2) | Test Status`

Phase 3 Coverage Matrix columns:

`Group | Test Item | Status (PoC verified / Negative tested / N/A / Untested) | Evidence Ref | Notes`

Coverage groups: **3.A** AuthZ · **3.B** Injection/input · **3.C** Modern API · **3.D** Business logic · **3.E** Client/browser edge · **3.F** Headers/CSP/HSTS (signal-based).

---

## G) Finding schema (each finding)

`Finding <ID>: <Title>`  
`- Severity: Critical|High|Medium|Low|Info`  
`- CVSS 3.1: ... or N/A`  
`- CWE: ... or N/A`  
`- OWASP: Axx:2021 / APIx:2023 or N/A`  
`- WSTG: WSTG-.... or N/A`  
`- ATT&CK: Txxxx or N/A`  
`- Affected: method/path/component/role`  
`- Prerequisites: ...`  
`- Non-technical summary: ...`  
`- Technical description: ...`  
`- Steps to reproduce: 1..2..3..`  
`- Evidence (redacted): ...`  
`- Impact: Business ... / Technical (CIA) ...`  
`- Remediation: priority 1; priority 2`  
`- Likelihood/FP checks: ...`  
`- Retest criteria: [ ] ...`

---

## H) Done criteria

- Every **P0** flow: Pass A + Pass B complete, or explicit blocker in Residual Risks (within profile limits).
- Phase 3 matrix: row per group **3.A–3.F** or `Skipped (FAST)` / `N/A` with reason.
- Each finding: evidence + retest criteria.
- No open branch with strong signal left **untested** without documenting why.
