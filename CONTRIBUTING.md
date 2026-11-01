# Contribuir

Leia `AGENTS.md` e o [guia de manutenção](docs/maintenance.md).
Use Go 1.26.5 ou superior e Python 3.8 ou superior para a geração offline.

Na raiz do SDK:

```powershell
python scripts/generate.py --check
go run ./scripts/check-api --check
go vet ./...
go build ./...
```

Para regenerar, execute o primeiro comando sem `--check`. Confira o diff dos
tipos, métodos e metadados antes de validar os consumidores.

## Validação no workspace

Testes, corpus e consumidores externos ficam exclusivamente em
`C:\server\raims\.tmp\raims-go-sdk`. Esse módulo usa `replace` para apontar
para o checkout do SDK e não acessa a API real. A bateria verifica os endpoints
e status publicados, autenticação, ambientes, serialização, paginação,
idempotência, cancelamento, erros e assinatura de webhooks.

```powershell
.\scripts\check.ps1
```

O script chama o validador canônico do checkout vizinho `Developer-Toolkit`,
que materializa uma suíte nova em `.tmp` e confere contrato, compatibilidade,
regressões, consumidor documentado, vulnerabilidades e exportação de fontes.
Ele não depende de arquivos de uma execução anterior. Um clone sem o Toolkit
informa a dependência ausente; os comandos acima continuam disponíveis para
geração, compatibilidade, vet e build.
Confira também as versões mínima e estável suportadas, concorrência com
`go test -race` em ambiente com compilador C e `govulncheck`.

Encerre processos e containers temporários antes de concluir. Preserve processos
preexistentes. Não inclua testes, corpus, relatórios, segredos ou arquivos de
validação no módulo, em commits ou em arquivos de distribuição.

Deixe alterações locais sem commit e sem conteúdo em staging. Uma validação
aprovada não autoriza publicação, tag, release ou configuração de infraestrutura.
