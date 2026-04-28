# Juice Shop Benchmark Scorecard (PentAGI)

Use this template after each benchmark run to prove measurable performance.

## Run Metadata
- run_id: `<timestamp-or-uuid>`
- date: `<YYYY-MM-DD>`
- model: `<e.g. Gemini 3.1>`
- prompt_file: `template_local_code_en_juice_shop.md`
- repo: `/work` (OWASP Juice Shop)
- benchmark_mode: `owasp_juice_shop`
- execution_profile: `THOROUGH`

## Key Metrics
- total_keys: `<from challenges.yml>`
- accounted_count: `<Confirmed + Probable + Out of scope + Not found in code>`
- coverage_percent: `<(accounted_count / total_keys) * 100>`
- confirmed_count: `<int>`
- probable_count: `<int>`
- out_of_scope_count: `<int>`
- not_found_count: `<int>`
- detection_percent: `<((confirmed_count + probable_count) / total_keys) * 100>`
- unmapped_keys_count: `<keys missing handler_files + endpoints + evidence_refs>`

## Hard Gates (Benchmark)
- gate_1_accounting_complete: `PASS|FAIL` (`accounted_count == total_keys`)
- gate_2_full_coverage: `PASS|FAIL` (`coverage_percent == 100`)
- gate_3_detection_floor: `PASS|FAIL` (`detection_percent >= 85`)
- gate_4_mapping_quality: `PASS|FAIL` (`unmapped_keys_count == 0`)
- overall_pass_or_fail: `PASS|FAIL`

## Failure Reasons (if any)
- `<reason 1>`
- `<reason 2>`

## Evidence Quality Check
- rows_with_valid_status: `<int>`
- rows_missing_status: `<int>`
- rows_missing_evidence_refs: `<int>`
- rows_missing_handler_and_endpoint: `<int>`
- sample_missing_rows: `<comma-separated keys>`

## Reproducibility (multi-run)
| run_id | model | total_keys | accounted_count | coverage_percent | confirmed_count | probable_count | detection_percent | unmapped_keys_count | overall |
|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| `<run-1>` | `<model>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<PASS/FAIL>` |
| `<run-2>` | `<model>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<PASS/FAIL>` |
| `<run-3>` | `<model>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<PASS/FAIL>` |

## Demo Narrative (for stakeholders)
- PentAGI accounted for all Juice Shop challenge keys with traceable evidence per key.
- Hard benchmark gates were enforced automatically; any missing challenge or weak mapping failed the run.
- Results are reproducible across repeated runs and suitable for audit/pentest readiness demonstrations.
