# Raims Go SDK

Read the README, CONTRIBUTING.md and docs/maintenance.md before changing the SDK.

- Keep repository descriptions and Markdown documentation in English.
- Keep internal TODO lists, backlogs, and roadmaps out of this public repository
  and its packages. Track operational work in the private workspace backlog
  maintained outside Git. Do not publish its contents, location, or links.
- The module is `github.com/raimshq/raims-go`, with package `raims` and only
  standard library dependencies at runtime. Preserve the client that can be
  shared by goroutines, services grouped by resource and `context.Context` on
  every call.
- The SDK consumes the published public OpenAPI. Its immutable source and
  SHA-256 are recorded in `contracts/source.json`. Do not import host code or
  include endpoints available only in the API's local checkout.
- Run `python scripts/generate.py`. Files named `*_generated.go`,
  `contracts/operations.json`, `contracts/schemas.json` and `docs/reference.md`
  are generated; do not edit them manually. The map must cover each pinned
  endpoint exactly once. `--check` verifies the outputs without rewriting them.
- The independent baseline `contracts/go-api.json` is checked by
  `go run ./scripts/check-api --check`. Only `--write` updates it; review the
  differences before doing so. The generator must not rewrite this baseline.
- Preserve int64, Unix seconds, future fields, omission and null. Inputs use
  `Optional[T]`; empty maps and nullable items must retain their meaning.
- An API key determines the environment and does not send an environment header.
  OAuth requires an explicitly authorized environment compatible with the
  endpoint; derive compatibility from the contract before sending HTTP.
- Do not generate idempotency keys, retry writes or start polling after `202`.
  The HTTP transport must also prevent write replay. Optional retries apply
  only to GET and honor Retry-After.
- Authenticated transport must not follow redirects, download external URLs or
  upload files. Do not expose a public method for arbitrary endpoints or headers.
- Diagnostics contain only the method, path without IDs, attempt, duration,
  status and Request-Id. Never include query strings, bodies or credentials.
- Tests, corpora, simulators, validation dependencies and reports belong
  exclusively in `C:\server\raims\.tmp\raims-go-sdk`. Do not version them.
- Before finishing, check generation, formatting, vet, build and external
  consumers, covering the contract, OAuth, pagination, retries and webhooks.
- `scripts/check.ps1` calls the canonical `go_sdk.py` adapter in Developer-Toolkit.
  It must reconstruct the suite in a new output directory under `.tmp`, without
  relying on existing tests. Keep physical test files outside this repository.
- Leave changes uncommitted and unstaged. Runner configuration, publication on
  GitHub and tags/releases are separate steps. Use only authorized CodeBuild
  runners when CI is enabled.
- Describe Raims features directly. Preserve licenses and the provenance of
  incorporated code. Apply the humanizer skill to new or changed prose.
