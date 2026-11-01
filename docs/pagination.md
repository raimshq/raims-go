# Paginação

`List` retorna somente uma página, sem iniciar outras chamadas:

```go
page, err := client.Customers.List(ctx, &raims.ListCustomersParams{
    PageSize: raims.Set[int64](20),
})
if err != nil {
    return err
}
for _, customer := range page.List {
    // Use customer.ID na sua aplicação.
}
next, err := page.NextPage(ctx)
```

`NextPage` e `PreviousPage` preservam filtros, ordenação, IDs e opções da chamada
original. Retornam nil quando não há continuação. Os cursores são opacos;
use os valores recebidos sem interpretá-los. Os parâmetros são copiados antes
da primeira chamada, para que alterações posteriores não troquem os filtros.

Para continuar automaticamente:

```go
for customer, err := range client.Customers.ListAll(ctx, nil).All() {
    if err != nil {
        return err
    }
    // Use customer.ID na sua aplicação.
}
```

O iterador busca páginas durante o consumo. `break` interrompe novos envios.
Também é possível usar `Next()`, `Current()` e `Err()`.
Não compartilhe um iterador entre goroutines.

`ListAll` usa a mesma validação e associação exata de nomes JSON de `List`.
Campos desconhecidos não substituem a lista, a continuação ou os cursores.

Uma página vazia pode ter continuação. O iterador segue esse cursor e recusa
cursores ausentes ou repetidos. O contexto permite cancelar o percurso ou
limitar seu prazo. Limite também a quantidade de recursos conforme a aplicação.

Na navegação manual, `NextPage` e `PreviousPage` também recusam ciclos com
mais de um cursor. O histórico acompanha cada sequência na mesma direção;
mudar de direção permite voltar às páginas já consultadas. Consultar novamente
a partir de uma página anterior mantém uma sequência independente.
