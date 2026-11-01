# Contributing

Read `AGENTS.md` and the [maintenance guide](docs/maintenance.md).
Use Go 1.26.5 or later and Python 3.8 or later for offline generation.

From the SDK root:

```powershell
python scripts/generate.py --check
go run ./scripts/check-api --check
go vet ./...
go build ./...
```

To regenerate the outputs, run the first command without `--check`. Review the
diff of types, methods and metadata before validating consumers.

## Workspace validation

Tests, corpora and external consumers belong exclusively in
`C:\server\raims\.tmp\raims-go-sdk`. The validation module uses `replace` to
point to the SDK checkout and does not access the live API. The suite checks
published endpoints and statuses, authentication, environments, serialization,
pagination, idempotency, cancellation, errors and webhook signatures.

```powershell
.\scripts\check.ps1
```

The script calls the canonical validator in the sibling `Developer-Toolkit`
checkout. It materializes a fresh suite in `.tmp` and checks the contract,
compatibility, regressions, the documented consumer, vulnerabilities and source
export. It does not rely on files from a previous run. A clone without the
Toolkit reports the missing dependency; the commands above remain available
for generation, compatibility, vet and build.
Also check the minimum and supported stable versions, concurrency with
`go test -race` in an environment with a C compiler, and `govulncheck`.

Stop temporary processes and containers before finishing. Preserve existing
processes. Do not include tests, corpora, reports, secrets or validation files
in the module, commits or distribution archives.

Leave local changes uncommitted and unstaged. Passing validation does not
authorize publication, tags, releases or infrastructure configuration.
