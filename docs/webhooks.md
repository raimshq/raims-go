# Receiving webhooks

Read the original body as bytes, with a size limit set by your server, before
parsing or modifying the JSON. Pass the received HTTP headers to the SDK:

```go
event, err := raims.ConstructWebhookEvent(payload, request.Header, signingSecret,
    raims.DefaultWebhookTolerance, time.Now())
if err != nil {
    // Reject the delivery according to your server's policy.
    return err
}
```

`ConstructWebhookEvent` authenticates the signature and returns an object with
`json.RawMessage` values, preserving integers and future fields. To authenticate
without parsing the body, use `VerifyWebhookSignature` with the same arguments.

The SDK checks `webhook-id`, `webhook-timestamp` and `webhook-signature`, rejects
missing or repeated headers, and verifies HMAC-SHA256 over the original bytes.
The secret uses the `whsec_` format returned by Raims; do not decode it before
passing it to the method.

The default tolerance is five minutes for past or future timestamps.
Zero tolerance requires the same second; it does not disable the time check.
A `time.Time{}` value uses the current clock.

After authentication, check the event type and expected environment. Persist
and deduplicate its identifier durably before applying the business effect.
The signature proves authenticity; it does not replace deduplication or, by
itself, confirm that a resource's financial operation is complete.

Keep the secret out of logs and version control. During rotation, use only valid
secrets and follow the period specified by the API.

## Payment instruction events

For new subscriptions, use `charge.payment_instructions_available` and
`charge.payment_instructions_failed`. Read `data.id` for the charge ID and
`data.payment_method` for the payment method. These event codes replace
`pix.available`, `pix.issuance_failed`, `boleto.available`, and
`boleto.issuance_failed` in the subscription catalog. Historical deliveries
can retain the earlier codes.

Read the event catalog for the authenticated environment before updating a
subscription. For refund and dispute events, use `data.charge_id` when the
event includes it.

Use the event and delivery history services to inspect saved events, deliveries,
and attempts. A response that accepts a resend or an earlier attempt does not
confirm that the destination received it.
