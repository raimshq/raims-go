# Pagination

`List` returns a single page without starting further calls:

```go
page, err := client.Customers.List(ctx, &raims.ListCustomersParams{
    PageSize: raims.Set[int64](20),
})
if err != nil {
    return err
}
for _, customer := range page.List {
    // Use customer.ID in your application.
}
next, err := page.NextPage(ctx)
```

`NextPage` and `PreviousPage` preserve the filters, ordering, IDs and options of
the original call. They return nil when there is no continuation. Cursors are
opaque; use the received values without interpreting them. Parameters are copied
before the first call so later changes cannot replace the filters.

To continue automatically:

```go
for customer, err := range client.Customers.ListAll(ctx, nil).All() {
    if err != nil {
        return err
    }
    // Use customer.ID in your application.
}
```

The iterator fetches pages as they are consumed. `break` stops new requests.
You can also use `Next()`, `Current()` and `Err()`.
Do not share an iterator between goroutines.

`ListAll` uses the same validation and exact JSON field name matching as `List`.
Unknown fields cannot replace the list, continuation information or cursors.

An empty page can have a continuation. The iterator follows that cursor and
rejects missing or repeated cursors. The context lets you cancel iteration or
set a deadline. Also limit the number of resources to suit your application.

During manual navigation, `NextPage` and `PreviousPage` also reject cycles with
more than one cursor. History follows each sequence in the same direction;
changing direction lets you return to pages already visited. Starting again
from an earlier page keeps an independent sequence.
