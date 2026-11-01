# Raims Go

SDK Go para integrar aplicações de servidor à API pública da Raims.
Requer Go 1.26.5 ou superior e usa somente a biblioteca padrão no runtime.

## Instalação e atualização

Em um projeto com `go.mod`, execute:

```shell
go get github.com/raimshq/raims-go@latest
```

Para atualizar, execute o mesmo comando e teste a integração da sua aplicação.
Revise as [notas das releases](https://github.com/raimshq/raims-go/releases)
antes de atualizar. O gerenciador registra a versão escolhida no projeto;
atualizações futuras não são aplicadas automaticamente.

## Primeiro uso

Importe o módulo `github.com/raimshq/raims-go` em sua aplicação:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    raims "github.com/raimshq/raims-go"
)

func main() {
    client, err := raims.NewClient(os.Getenv("RAIMS_API_KEY"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.CloseIdleConnections()

    page, err := client.Customers.List(context.Background(), &raims.ListCustomersParams{
        PageSize: raims.Set[int64](20),
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, customer := range page.List {
        fmt.Println(customer.ID)
    }
}
```

Reutilize o cliente entre requisições e goroutines. A chave define empresa,
ambiente e permissões. A URL-base é `https://api.raims.com`.
Mantenha credenciais, documentos pessoais e links temporários fora do navegador,
do controle de versão e dos logs.

## Criar uma cobrança

Com o cliente configurado, use o ID de um cliente da conta e os dados do pedido:

```go
charge, err := client.Charges.Create(ctx, &raims.CreateChargeParams{
    CustomerID:  raims.Set(customerID),
    AmountCents: raims.Set[int64](15000),
    DueDate:     raims.Set(dueDateUnix),
    Title:       raims.Set("Pedido 1234"),
}, raims.WithIdempotencyKey(persistedOrderKey))
if err != nil {
    return err
}
requestID := charge.LastResponse.RequestID
pending := charge.LastResponse.StatusCode == 202
```

`customerID` é o identificador retornado pela API, `dueDateUnix` é um `int64`
com o vencimento escolhido em segundos Unix UTC e `persistedOrderKey` identifica
a intenção de criar a cobrança. Persista a chave e os parâmetros antes do envio.
`15000` representa R$ 150,00. Uma resposta `202` indica resultado pendente;
acompanhe o recurso antes de concluir ou repetir a intenção.

## Serviços

| Serviço | Recursos |
| --- | --- |
| `Accounts` | Contexto, dados da conta, tarifas e checkout |
| `Balances` | Saldo e previsão |
| `Customers` | Cadastro e consulta de clientes |
| `Charges` | Cobranças, pagamento, documentos e estornos por cobrança |
| `PaymentLinks` | Links de pagamento e solicitação de envio por e-mail |
| `Refunds` | Estornos e recibos |
| `Disputes` | Contestações |
| `Statements` | Extratos |
| `Withdrawals` | Saques, limites, resumo e configuração automática |
| `WithdrawalDestinations` | Destinos de saque |
| `Webhooks` | Cadastro, pausa, saúde, teste e rotação de segredo |
| `WebhookEventTypes` | Tipos de evento |
| `Simulations.Charges` e `Simulations.Refunds` | Simulações com credencial Sandbox |

Os serviços oferecem métodos como `Create`, `Get`, `Update` e `List`, conforme
o recurso. `ListAll` retorna um iterador com `Next`, `Current`, `Err` e `All`.

## Guias

- [Configuração, OAuth e atualizações parciais](docs/usage.md)
- [Paginação](docs/pagination.md)
- [Erros, idempotência e retries](docs/errors-and-retries.md)
- [Recebimento de webhooks](docs/webhooks.md)
- [Referência por endpoint](docs/reference.md)
- [Contribuição](CONTRIBUTING.md), [manutenção](docs/maintenance.md) e [segurança](SECURITY.md)
- [Documentação da API](https://docs.raims.com/br/api)

Distribuído sob a licença [MIT](LICENSE).
