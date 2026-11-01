# Go SDK static analysis

The preflight uses Semgrep 1.181.0 and the official rules at revision
`9bbf021d0acffa3095fcffd17b2bc7787eb4796d`, together with the local rules in
`security/semgrep.yml`.

The analysis covers the entire SDK root, including generated Go files and
engineering scripts in Go and Python. The result must prove that every `.go`
and `.py` source was read; only artifact directories defined by the canonical
validator are excluded.

The local rules check TLS verification, isolated HTTP client configuration,
shell execution and embedded credentials. The official Go and Python rules
complement this analysis. Scanner errors, incomplete coverage or unresolved
findings block validation.

The canonical helper requires Python 3.12. In the example below, `python` must
point to that version. On Windows, the `Test-RaimsRelease.ps1` launcher also
accepts `-PythonPath` to select the executable.

Run the command from the workspace:

```powershell
python Developer-Toolkit/scripts/release-preflight/sdk_security.py --source Raims-Go --workspace . --repository Raims-Go --output .tmp/raims-go-sdk/security --prepare-tools
```

Use a new or empty output directory. Reports stay in `.tmp`. This check does
not create infrastructure, enable workflows or publish the SDK.
