# Raims Go SDK

Leia o README, CONTRIBUTING.md e docs/maintenance.md antes de alterar o SDK.

- O módulo é `github.com/raimshq/raims-go`, com package `raims` e somente
  biblioteca padrão no runtime. Preserve o cliente reutilizável por goroutines,
  os serviços por recurso e `context.Context` em todas as chamadas.
- O SDK consome o OpenAPI público publicado. A origem imutável e o SHA-256 ficam
  em `contracts/source.json`. Não importe código dos hosts nem inclua endpoints
  disponíveis somente no checkout local da API.
- Execute `python scripts/generate.py`. Arquivos `*_generated.go`,
  `contracts/operations.json`, `contracts/schemas.json` e `docs/reference.md`
  são gerados; não os edite manualmente. O mapa precisa cobrir cada endpoint
  fixado exatamente uma vez. `--check` confere sem regravar.
- A base independente `contracts/go-api.json` é conferida por
  `go run ./scripts/check-api --check`. Somente `--write` atualiza a base;
  revise as diferenças antes dessa atualização. Não a regrave pelo gerador.
- Preserve int64, Unix segundos, campos futuros, omissão e null. Entradas usam
  `Optional[T]`; mapas vazios e itens nullable precisam manter seu significado.
- A chave de API determina o ambiente e não envia header de ambiente. OAuth
  exige ambiente explicitamente autorizado e compatível com o endpoint;
  derive a compatibilidade do contrato, antes de enviar HTTP.
- Não gere idempotência, repita escritas ou inicie polling após `202`.
  Impeça também replay de escrita pelo transporte HTTP. Retries opcionais
  são exclusivos de GET e respeitam Retry-After.
- O transporte autenticado não segue redirects, baixa URLs externas nem envia
  uploads. Não exponha um método público para endpoints ou headers arbitrários.
- Diagnósticos contêm somente método, caminho sem IDs, tentativa, duração,
  status e Request-Id. Nunca incluir query, corpo ou credencial.
- Testes, corpus, simuladores, dependências de validação e relatórios ficam
  exclusivamente em `C:\server\raims\.tmp\raims-go-sdk`. Não versioná-los.
- Antes de concluir, confira geração, formatação, vet, build e consumidores
  externos, incluindo contrato, OAuth, paginação, retries e webhooks.
- `scripts/check.ps1` chama o adapter canônico `go_sdk.py` do Developer-Toolkit.
  Ele deve reconstruir a suíte em uma saída nova dentro de `.tmp`, sem depender
  de testes preexistentes. Preserve os cenários físicos fora deste repositório.
- Deixe mudanças sem commit e sem conteúdo em staging. Configuração de runners,
  publicação no GitHub e tags/releases são etapas próprias. Use apenas runners
  CodeBuild autorizados quando o CI for habilitado.
- Descreva os recursos da Raims diretamente. Preserve licenças e procedência
  de código incorporado. Aplique a skill humanizer à prosa criada ou alterada.
