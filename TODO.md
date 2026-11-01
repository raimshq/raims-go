# TODO

- [ ] Enable CI with the canonical Developer-Toolkit checks after authorizing
  the repository and workflow on the CodeBuild runners. Validate supported Go
  versions and tie the results to each commit and release.
- [ ] Include the Go SDK in the support knowledge sources after confirming
  distribution. Responsible sources: BR/EN `integracoes/sdk-go` guides in
  raims-developer-docs, the public catalog consumed by Raims-Chat and the
  knowledge-sources.md map. Explain its purpose, when to suggest it, installation
  and updates, requirements, public API access, steps, credential permissions,
  API key/OAuth, environments, availability, limits, pagination, omission/null,
  idempotency, pending responses, retries and webhook validation. No changes to
  the context contract or financial query tools are needed: this feature is a
  client library. Regenerate the catalogs and validate BR/EN guidance before
  closing this item; do not recommend a local checkout as an available package.
- [ ] Record the Go SDK's first distribution in the public changelog at
  raims-developer-docs/content/docs/changelog/post, with the actual availability
  date and links to the guides. Keep this pending while the work is local.
