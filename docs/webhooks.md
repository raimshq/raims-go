# Receber webhooks

Leia o corpo original como bytes, com limite definido pelo seu servidor, antes
de interpretar ou modificar o JSON. Passe os headers HTTP recebidos ao SDK:

```go
event, err := raims.ConstructWebhookEvent(payload, request.Header, signingSecret,
    raims.DefaultWebhookTolerance, time.Now())
if err != nil {
    // Recuse a entrega conforme a política do seu servidor.
    return err
}
```

`ConstructWebhookEvent` autentica a assinatura e retorna um objeto com valores
`json.RawMessage`, preservando inteiros e campos futuros. Para autenticar sem
interpretar o corpo, use `VerifyWebhookSignature` com os mesmos argumentos.

O SDK confere `webhook-id`, `webhook-timestamp` e `webhook-signature`, recusa
headers ausentes ou repetidos e verifica HMAC-SHA256 sobre os bytes originais.
O segredo tem o formato `whsec_` retornado pela Raims; não o decodifique antes
de passá-lo ao método.

A tolerância padrão é cinco minutos, para timestamps passados ou futuros.
Tolerância zero exige o mesmo segundo; ela não desliga a verificação de tempo.
Uma data `time.Time{}` usa o relógio atual.

Depois de autenticar, confira o tipo do evento e o ambiente esperado. Grave e
deduplique seu identificador de forma durável antes de aplicar o efeito de
negócio. A assinatura comprova autenticidade; não substitui essa deduplicação
nem confirma, sozinha, a conclusão financeira de um recurso.

Mantenha o segredo fora dos logs e do controle de versão. Durante uma rotação,
use somente os segredos válidos e siga o período informado pela API.
