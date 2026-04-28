# PentAGI — Common Source Code Security Audit Prompt (Generic, Reusable)

## 1) Goal
Find actionable security vulnerabilities in source code, configs, CI/CD, and container artifacts under `/work`.
Prioritize real exploitability, clear evidence, and concise reporting.

## 2) Execution Modes
- `FAST`: quick triage, high-signal coverage.
- `BALANCED` (default): practical depth for most repos.
- `THOROUGH`: deeper review and broader coverage.

## 3) Minimal Input Block (required)
```text
=== CODE_SOURCE_INPUT ===
execution_profile: BALANCED
repo_path: /Users/vinhtv/Documents/Audit-pentest/demo-hacked/juice-shop
sub_path:
exclude_paths: .git,node_modules,dist,build,.next,.cache,vendor,tmp,.pnpm-store,.vscode
internet_allowed: true
run_install: false
dynamic_smoke_allowed: false
pentest_verify_allowed: false
benchmark_mode: none
app_profile: production
=== END CODE_SOURCE_INPUT ===
```

### Optional inputs
- `benchmark_mode`: `none` (default) | `training_lab` | `owasp_juice_shop`
- `app_profile`: `production` (default) | `training_lab`
- `base_url`: required only for dynamic verification

## 4) Auto-Detection (no manual hints required)
During Phase 0, detect project traits automatically:
- language/framework (Node, Java, Python, Go, Rust, .NET, etc.)
- API style (REST/GraphQL/WS)
- data layer (SQL/ORM/NoSQL/queue/cache)
- CI/container signals (`.github/workflows`, Dockerfile, compose, helm)
- lab-like markers (CTF/training challenge hooks, intentionally vulnerable modules)

If strong lab markers are detected, set `effective_mode=training_lab` and record why.

### 4.1) Challenge-Discovery Auto Switch (mandatory)
If one or more markers exist, enable challenge-discovery flow even when operator did not set benchmark:
- `solveIf(challenges.`
- `notSolved(challenges.`
- `challengeUtils.solveIf(challenges.`
- `vuln-code-snippet`
- static challenge file exists (for example `data/static/challenges.yml`)

When auto-switch activates:
- set `effective_benchmark_mode=training_lab`
- emit `benchmark_auto_detected: true` + marker list in metadata
- produce `Challenge/Module Coverage Map` with file anchors

### 4.2) Key-Set Discovery (mandatory in training-lab mode)
Build challenge/module key set via fallback:
1. Preferred: parse static challenge registry (if present).
2. Fallback: infer keys from source anchors (`solveIf/notSolved/challengeUtils.solveIf`).
3. Ignore non-key properties (`map`, `filter`, `some`, `forEach`, etc.) when inferred.

### 4.3) Canonical Key Integrity (hard gate)
- If static challenge registry exists (for example `data/static/challenges.yml`), it is the canonical source of truth.
- **Never truncate/slice/filter out** canonical keys to satisfy accounting targets.
- **Hard fail** if any step attempts to reduce canonical key count for convenience.
- Record in metadata:
  - `canonical_key_source`
  - `canonical_total_keys`
  - `effective_total_keys`
  - `key_integrity_check: pass|fail`

### 4.4) Challenge Classification (mandatory in training-lab mode)
For each discovered challenge key, classify into one of:
- `security_vuln` (e.g., SQLi, XSS, XXE, SSRF, RCE, auth bypass, IDOR, traversal, weak crypto, CSRF)
- `business_logic`
- `informational/ctf-mechanic`
This classification is required to avoid reporting only broad accounting without real vulnerability depth.

## 5) Method (tool-first)
### Phase 0 — Fingerprint
- Validate workspace and project identity.
- Build architecture map and prioritized hotspot list.

### Phase 1 — Automated scans
- Secret scan (gitleaks/equivalent).
- Dependency scan (one or two engines based on profile/time).
- Broad SAST (semgrep/ruleset available in runtime).
- Pattern scan (`rg`) for high-risk sinks and dangerous flows.
- Route/entry inventory (HTTP/GraphQL/WS/background jobs/CLI).

### Phase 1.5 — Prioritized reading
- Rank files by severity + sink density + auth/write-path proximity.
- Read Top-K only (profile-capped), then targeted callees if needed.
- In training-lab mode, force route-to-handler deep queue for state-changing endpoints even when scanners are quiet.
- In training-lab mode, perform **challenge-key-first traversal**:
  - for each key, resolve at least one code anchor (`solveIf/notSolved/challenges.<key>`)
  - map anchor to executable handler/endpoint when applicable
  - keep unresolved keys in an explicit queue until final accounting
  - prioritize deep review on `security_vuln` keys first (before business/info keys)

### Phase 2 — Manual security review
Focus order:
1) auth/authz/session
2) injection (SQL/NoSQL/command/template)
3) SSRF/RCE/deserialization
4) file upload/path/archive handling
5) business logic abuse (price/refund/coupon/idempotency/rate-limit)
6) CI/container misconfiguration

### Phase 3 — Triage and dedupe
- Merge duplicates by root cause.
- Keep one representative finding + blast radius.

## 6) Generic Risk Patterns (seed set)
- code execution: `eval`, `new Function`, shell exec/spawn
- query injection: raw SQL/string concat/query bypass APIs
- NoSQL injection: untrusted operators/regex/filter objects
- insecure auth/jwt/session handling
- SSRF via untrusted URL fetch
- weak crypto/randomness/hardcoded secrets
- path traversal/file overwrite/archive extraction issues
- dangerous HTML/script rendering on untrusted input
- misconfigured CORS/CSRF/cookies/headers

## 7) Output Contract (mandatory)
1. `# Executive Summary`
2. `## One-page Risk Snapshot`
3. `## Engagement Metadata`
4. `## Codebase Fingerprint`
5. `## Attack Surface`
6. `## Findings`
7. `## Static Coverage Matrix`
8. `## Dependency & Supply Chain Summary`
9. `## Container & CI Signals`
10. `## Residual Risks & Assumptions`
11. `## Remediation Roadmap (30/60/90 days)`
12. `## Retest / CI Gates`

If `effective_mode=training_lab`, additionally include:
- `## Challenge/Module Coverage Map` (only if challenge-like keys/modules are detected)
- `## Challenge Key Source` (`static_file|code_inference`, source path/pattern, confidence note)
- `## Benchmark Scorecard` (`total_keys`, `accounted_count`, `coverage_percent`, `confirmed_count`, `probable_count`, `not_found_count`, `unmapped_count`, `pass_or_fail`)
- `## Unmapped Challenge Keys` (explicit list, empty if none)
- `## Challenge Technical Breakdown` with columns:
  - `challenge_key | class(security_vuln/business_logic/info) | primary_anchor(file:line) | endpoint_or_handler | status(Confirmed/Probable/Out_of_scope/Not_found) | evidence_ref`
- `## Security Challenge Findings` (subset of `security_vuln` keys with full source->sink analysis)
- `## Detection Threshold Gate` with fields:
  - `target_detected_percent=70`
  - `detected_keys = confirmed_count + probable_count`
  - `detected_percent = (detected_keys / total_keys) * 100`
  - `security_detected_keys = security_vuln_confirmed + security_vuln_probable`
  - `security_detected_percent = (security_detected_keys / security_vuln_total) * 100`
  - `threshold_pass_or_fail`

## 8) Finding Schema
For each finding include:
- severity, CWE/CVSS (if applicable)
- affected file + line + symbol
- source -> sink data flow
- exploitability notes
- impact (business/CIA)
- remediation (short + practical)
- retest criteria

## 9) Done Criteria
Finish only when:
- fingerprint + automated scans completed (or blocker documented),
- all high-severity hotspots in Top-K were reviewed,
- each Confirmed/Probable finding has evidence,
- residual risks explicitly list what was not covered and why.
- hard fail if any Confirmed/Probable finding lacks `file + line + source->sink`.

### 9.1) Additional done criteria in training-lab mode
- one row per discovered key/module in `Challenge/Module Coverage Map`.
- `accounted_count = confirmed + probable + out_of_scope + not_found`.
- hard fail when `accounted_count != total_keys`.
- if key source is inferred (no static registry), include explicit confidence warning in scorecard.
- if static registry exists, require `total_keys == canonical_total_keys` from registry.
- hard fail when `total_keys` in scorecard differs from canonical registry count.
- hard fail if `Unmapped Challenge Keys` is omitted.
- for `execution_profile=THOROUGH`, target `not_found_count = 0`; if not zero, require concrete residual reasons per key.
- require explicit count summary by class:
  - `security_vuln_total`, `security_vuln_confirmed`, `security_vuln_probable`
- hard fail if `security_vuln` class is missing from report while challenge markers are present.
- for `execution_profile=THOROUGH`, require deep evidence for at least top 10 `security_vuln` keys by exploitability (or all if fewer than 10).
- compute `detected_keys = confirmed + probable` and `detected_percent = (detected_keys / total_keys) * 100`.
- **hard fail when `detected_percent < 70`**.
- compute `security_detected_percent = ((security_vuln_confirmed + security_vuln_probable) / security_vuln_total) * 100` (if `security_vuln_total > 0`).
- **hard fail when `security_detected_percent < 70`**.
- threshold verdict must be printed in `## Detection Threshold Gate` with explicit failed conditions.

## 10) Safety and Accuracy
- Do not claim confirmation without code evidence.
- Mark uncertain items as `Probable` with upgrade conditions.
- Do not leak plaintext secrets in report.
- Keep report concise, evidence-first, and non-duplicative.
