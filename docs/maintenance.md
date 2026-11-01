# Contract and maintenance

## Source and generation

The contract pinned in `contracts/openapi.json` comes from the public
[raimshq/openapi](https://github.com/raimshq/openapi) mirror.
`contracts/source.json` records the immutable revision, URL, SHA-256, API version
and endpoint count. The canonical source remains in Public-API; the published
mirror defines the surface consumed by this SDK.

The generator runs offline. It does not read credentials, download contracts or
call endpoints. `contracts/service-map.json` explicitly maps each endpoint to
a method. The generator rejects mismatched hashes, missing or repeated endpoints
and compositions without an implemented representation.

| File | Responsibility |
| --- | --- |
| `client.go` | Credentials, composition, transport, timeout, retries and diagnostics |
| `optional.go` | Omission, values and null in inputs |
| `resource.go` | Metadata, original fields and errors |
| `contract.go` | Prepared metadata, validation, exact reading and response serialization |
| `pagination.go` | Cursor navigation and iteration |
| `webhook.go` | Delivery authentication |
| `*_generated.go` | Services, parameters, models and constants |
| `contracts/operations.json` and `contracts/schemas.json` | Metadata used at runtime |
| `scripts/generate.py` | Deterministic generation |

To update, select an approved public revision, verify its manifest and checksum,
and update the snapshot and source metadata. Review changes to authentication,
environments, fields, nullability, statuses, pagination and idempotency. Adjust
the map and run generation. Do not add routes available only in local code.

Responses accept future fields and enum values. Inputs retain their own types.
Unions of multiple concrete types, inline responses and changes to the
combination of authentication mechanisms require explicit support before the
contract is updated. The generator rejects these forms when it cannot preserve
their meaning.

`contracts/go-api.json` is the independent baseline of public symbols. Run
`go run ./scripts/check-api --check` to verify package identity, types, methods,
fields, constants and JSON tags. A change in signature, type or meaning requires
a compatibility review before choosing a version. OpenAPI generation and the
compatibility check do not rewrite the baseline. Update it explicitly with
`go run ./scripts/check-api --write` after reviewing the differences.

The module uses only the standard library. Validate the minimum supported
version and the current stable version, including Go's vulnerability tool.
The version in `client.go` identifies the SDK; it does not change the HTTP
version `/v1`.

## Distribution and support

Local preparation does not confirm publication. Before distributing, verify the
remote repository's identity, module metadata, license and included files.
Tags, releases and publication require separate authorization. Do not configure
CI before authorizing the repository and its workflow on the CodeBuild runners.
Preserve security checks and their link to the published commit.

After confirming distribution, prepare the BR/EN guides in
`raims-developer-docs` and the changelog entry with the actual availability date.
The guide must cover installation and updates through the Go module,
requirements, credentials, permissions, environments, pagination, null,
idempotency, errors, pending responses and webhooks. Use
`go get github.com/raimshq/raims-go@latest` to obtain the current stable version;
do not pin a release in installation examples.

This guidance belongs to the existing `raims-developer-docs` catalog used by
Raims-Chat. The SDK does not change authenticated context, Dashboard permissions
or financial query tools. This library does not require a context schema change.
Keep the TODO item open until the sources are published, the catalog is
regenerated and the guidance is validated in both languages. Do not present a
local checkout as an available package. Do not index this repository's internal
documents in the chat.
