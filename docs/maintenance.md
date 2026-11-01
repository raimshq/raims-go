# Contrato e manutenção

## Fonte e geração

O contrato fixado em `contracts/openapi.json` vem do espelho público
[raimshq/openapi](https://github.com/raimshq/openapi). `contracts/source.json`
registra a revisão imutável, URL, SHA-256, versão da API e quantidade de endpoints.
A fonte canônica permanece na Public-API; o espelho publicado delimita a
superfície consumida por este SDK.

O gerador é offline. Não lê credenciais, baixa contratos nem chama endpoints.
`contracts/service-map.json` mapeia explicitamente cada endpoint para um método.
O gerador recusa hash divergente, endpoints ausentes ou repetidos e composições
que ainda não possuem representação implementada.

| Arquivo | Responsabilidade |
| --- | --- |
| `client.go` | Credenciais, composição, transporte, timeout, retries e diagnósticos |
| `optional.go` | Omissão, valor e null em entradas |
| `resource.go` | Metadados, campos originais e erros |
| `contract.go` | Metadados preparados, validação, leitura exata e serialização das respostas |
| `pagination.go` | Navegação e iteração por cursor |
| `webhook.go` | Autenticação de entregas |
| `*_generated.go` | Serviços, parâmetros, modelos e constantes |
| `contracts/operations.json` e `contracts/schemas.json` | Metadados usados no runtime |
| `scripts/generate.py` | Geração determinística |

Para atualizar, selecione uma revisão pública aprovada, confira manifesto e
checksum e atualize o snapshot e a origem. Revise mudanças de autenticação,
ambiente, campos, nullability, status, paginação e idempotência. Ajuste o mapa
e execute a geração. Não acrescente rotas disponíveis somente no código local.

As respostas aceitam campos e valores de enum futuros. Entradas mantêm tipos
próprios. Uniões com vários tipos concretos, respostas inline e alterações na
combinação dos mecanismos de autenticação exigem suporte explícito antes de
atualizar o contrato. O gerador recusa essas formas quando não consegue
preservar seu significado.

`contracts/go-api.json` é a base independente dos símbolos públicos. Execute
`go run ./scripts/check-api --check` para conferir a identidade do pacote,
tipos, métodos, campos, constantes e tags JSON. Uma mudança de assinatura,
tipo ou significado exige análise de compatibilidade antes de definir a versão.
A geração do OpenAPI e
o check de compatibilidade não regravam a base. Sua atualização é explícita,
com `go run ./scripts/check-api --write`, após revisar as diferenças.

O módulo usa somente a biblioteca padrão. Valide a versão mínima suportada e
a estável atual, incluindo a ferramenta de vulnerabilidades de Go. A versão
em `client.go` identifica o SDK; ela não altera a versão HTTP `/v1`.

## Distribuição e atendimento

A preparação local não confirma publicação. Antes de distribuir, confira a
identidade do repositório remoto, os metadados do módulo, a licença e os arquivos
incluídos. Tags, releases e publicação seguem autorização própria.
Não configure CI antes de autorizar o repositório e seu workflow nos runners
CodeBuild. Preserve os checks de segurança e o vínculo com o commit publicado.

Após confirmar a distribuição, prepare os guias BR/EN em
`raims-developer-docs` e a nota de changelog com a data real de disponibilização.
O guia deve cobrir instalação e atualização pelo módulo Go, requisitos,
credenciais, permissões, ambiente, paginação, null, idempotência, erros, respostas
pendentes e webhooks. Use `go get github.com/raimshq/raims-go@latest` para obter
a versão estável atual; não fixe uma release nos exemplos de instalação.

Essa orientação pertence ao catálogo existente `raims-developer-docs` usado
pelo Raims-Chat. O SDK não altera contexto autenticado, permissões do Dashboard
nem ferramentas de consulta financeira. Não há necessidade de mudar o schema
do contexto por causa desta biblioteca.
Mantenha o item do TODO aberto até publicar as fontes, regenerar o catálogo e
validar a orientação nos dois idiomas. Não anuncie um checkout local como
pacote disponível. Não indexe documentos internos deste repositório no chat.
