# Raims Go

Go SDK for server applications that use the Raims public API.
Requires Go 1.26.5 or later and uses only the standard library at runtime.

## Installation and updates

In a project with a `go.mod` file, run:

```shell
go get github.com/raimshq/raims-go@latest
```

To update, run the same command and test your application's integration.
Review the [release notes](https://github.com/raimshq/raims-go/releases)
before updating. The package manager records the selected version in your
project; future updates are not applied automatically.

## First use

Import the `github.com/raimshq/raims-go` module in your application:

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

Reuse the client across requests and goroutines. The key determines the company,
environment and permissions. The base URL is `https://api.raims.com`.
Keep credentials, personal identity documents and temporary links out of the
browser, version control and logs.

## Create a charge

With the client configured, use a customer ID from the account and the order data:

```go
charge, err := client.Charges.Create(ctx, &raims.CreateChargeParams{
    CustomerID:  raims.Set(customerID),
    AmountCents: raims.Set[int64](15000),
    DueDate:     raims.Set(dueDateUnix),
    Title:       raims.Set("Order 1234"),
}, raims.WithIdempotencyKey(persistedOrderKey))
if err != nil {
    return err
}
requestID := charge.LastResponse.RequestID
pending := charge.LastResponse.StatusCode == 202
```

`customerID` is the identifier returned by the API, `dueDateUnix` is an `int64`
with the selected due date in Unix seconds UTC, and `persistedOrderKey` identifies
the intent to create the charge. Persist the key and parameters before sending
the request. `15000` represents BRL 150.00. A `202` response means the result is
pending; check the resource before treating the operation as complete or
retrying the intent.

## Services

| Service | Resources |
| --- | --- |
| `Accounts` | Context, account details, fees and checkout |
| `Balances` | Balance and forecast |
| `Customers` | Customer registration and lookup |
| `Charges` | Charges, payments, documents and refunds for a charge |
| `PaymentLinks` | Payment links and requests to send them by email |
| `Refunds` | Refunds and receipts |
| `Disputes` | Disputes |
| `Statements` | Statements |
| `Withdrawals` | Withdrawals, limits, summary and automatic settings |
| `WithdrawalDestinations` | Withdrawal destinations |
| `Webhooks` | Registration, pause, health, testing and secret rotation |
| `WebhookEventTypes` | Event types |
| `Simulations.Charges` and `Simulations.Refunds` | Simulations with a Sandbox credential |

Services provide methods such as `Create`, `Get`, `Update` and `List`, depending
on the resource. `ListAll` returns an iterator with `Next`, `Current`, `Err` and `All`.

## Guides

- [Configuration, OAuth and partial updates](docs/usage.md)
- [Pagination](docs/pagination.md)
- [Errors, idempotency and retries](docs/errors-and-retries.md)
- [Receiving webhooks](docs/webhooks.md)
- [Endpoint reference](docs/reference.md)
- [Contributing](CONTRIBUTING.md), [maintenance](docs/maintenance.md) and [security](SECURITY.md)
- [API documentation](https://docs.raims.com/en/api)

Distributed under the [MIT license](LICENSE).
