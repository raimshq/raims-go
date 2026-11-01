# Erros, idempotência e retries

## Erros

Use `errors.As` para conferir a categoria e `errors.Is` para cancelamentos:

```go
var apiError *raims.APIError
if errors.As(err, &apiError) {
    status := apiError.StatusCode
    code := apiError.Code
    requestID := apiError.Response.RequestID
    // Decida a ação pelo status e pelo código.
}
```

| Tipo | Significado |
| --- | --- |
| `ValidationError` | Entrada recusada antes de HTTP |
| `APIError` | Resposta HTTP de falha, incluindo 429 e 5xx |
| `ConnectionError` | Resposta interrompida, conexão ou timeout |
| `ProtocolError` | Redirect, status inesperado, resposta incompatível ou cursor inválido |
| `SignatureError` | Falha na autenticação de uma entrega webhook |

`APIError.Detail` contém código, mensagem e detalhes de validação quando a API
retorna um envelope válido. Proxies podem responder sem esse envelope.
`ProtocolError.Response` preserva os metadados HTTP recebidos.
Os textos de `Error()` das falhas HTTP e de conexão não incluem o payload.
Não registre a exceção desembrulhada, structs completos ou mensagens da API sem
conferir os dados que contêm. Uma conexão perdida durante uma escrita deixa
seu resultado incerto; ela pode ter sido processada.

## Idempotência

Passe `WithIdempotencyKey(persistedKey)` somente em endpoints que aceitam o
header, identificados na [referência](reference.md). A chave deve conter entre
1 e 255 caracteres ASCII visíveis, sem espaços.

Crie e persista a chave antes do envio, junto dos parâmetros da intenção.
Se a resposta se perder, conserve ambos. O SDK não gera chaves nem repete
escritas. `IDEMPOTENCY_REQUEST_IN_PROGRESS` pede espera conforme `Retry-After`;
`IDEMPOTENCY_RESULT_UNKNOWN` exige consultar e reconciliar antes de repetir.
Consulte as [regras da API](https://docs.raims.com/br/api/idempotencia).

## Retries

O padrão é zero retries. `WithMaxNetworkRetries` aceita de zero a cinco
tentativas adicionais, exclusivamente para GET. São elegíveis falhas de
conexão e HTTP 408, 429, 500, 502, 503 e 504.

O backoff tem jitter e respeita `Retry-After`, em segundos ou data HTTP.
Uma espera superior a 30 segundos devolve a falha para a aplicação reagendar.
Cada tentativa tem seu timeout; um contexto com prazo limita a duração total.

POST, PATCH e DELETE têm um único envio, inclusive com chave de idempotência.
O transporte padrão também impede replay dessas requisições ao perder uma
conexão reutilizada. Um transporte personalizado precisa manter essa garantia.
