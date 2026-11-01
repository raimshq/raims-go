# Go SDK reference

Generated from the pinned public contract. `Params` uses `Set(value)` to send a field and `Null[T]()` to send null. All methods accept `context.Context`, the IDs shown in the path and optional request options. Inspect type fields with `go doc`.

| Service and method | Endpoint | Parameters | Return type | Status | OAuth environments | Idempotency |
| --- | --- | --- | --- | --- | --- | --- |
| `Accounts.GetFees` | `GET /v1/accounts/fees` | `none` | `AccountFees` | 200 | Sandbox, Production | No |
| `Accounts.GetNumber` | `GET /v1/accounts/account_number` | `none` | `AccountNumber` | 200 | Production | No |
| `Accounts.GetCheckoutSettings` | `GET /v1/accounts/checkout_settings` | `none` | `CheckoutSettings` | 200 | Sandbox, Production | No |
| `Accounts.UpdateCheckoutSettings` | `PATCH /v1/accounts/checkout_settings` | `UpdateCheckoutSettingsParams` | `CheckoutSettings` | 200, 202 | Sandbox, Production | Yes |
| `Accounts.CreateCheckoutLogoUpload` | `POST /v1/accounts/checkout_settings/logo_uploads` | `CreateCheckoutLogoUploadParams` | `CheckoutAssetUpload` | 201 | Sandbox, Production | Yes |
| `Accounts.CreateCheckoutFaviconUpload` | `POST /v1/accounts/checkout_settings/favicon_uploads` | `CreateCheckoutFaviconUploadParams` | `CheckoutAssetUpload` | 201 | Sandbox, Production | Yes |
| `Accounts.DeleteCheckoutAssetUpload` | `DELETE /v1/accounts/checkout_settings/uploads/{id}` | `none` | `CheckoutAssetUploadDeleted` | 200 | Sandbox, Production | Yes |
| `Accounts.GetContext` | `GET /v1/accounts/context` | `none` | `AccountContext` | 200 | Sandbox, Production | No |
| `Balances.Get` | `GET /v1/balances` | `none` | `Balance` | 200 | Sandbox, Production | No |
| `Balances.Forecast` | `GET /v1/balances/forecast` | `GetBalanceForecastParams` | `BalanceForecast` | 200 | Sandbox, Production | No |
| `Charges.Create` | `POST /v1/charges` | `CreateChargeParams` | `ChargeCreated` | 201, 202 | Sandbox, Production | Yes |
| `Charges.List` | `GET /v1/charges` | `ListChargesParams` | `ChargeListResponse` | 200 | Sandbox, Production | No |
| `Charges.Get` | `GET /v1/charges/{id}` | `none` | `Charge` | 200 | Sandbox, Production | No |
| `Charges.Update` | `PATCH /v1/charges/{id}` | `UpdateChargeParams` | `Charge` | 200 | Sandbox, Production | Yes |
| `Charges.Revoke` | `DELETE /v1/charges/{id}` | `none` | `RevokeChargeResult` | 200, 202 | Sandbox, Production | Yes |
| `Charges.GetPayment` | `GET /v1/charges/{id}/payment` | `none` | `PaymentReceipt` | 200 | Sandbox, Production | No |
| `Charges.GetReceipt` | `GET /v1/charges/{id}/receipt` | `none` | `FinancialDocument` | 200 | Sandbox, Production | No |
| `Charges.GetBoleto` | `GET /v1/charges/{id}/boleto` | `none` | `FinancialDocument` | 200 | Sandbox, Production | No |
| `Charges.CreateRefund` | `POST /v1/charges/{id}/refunds` | `CreateRefundParams` | `Refund` | 200, 202 | Sandbox, Production | Yes |
| `Charges.ListRefunds` | `GET /v1/charges/{id}/refunds` | `ListChargeRefundsParams` | `RefundListResponse` | 200 | Sandbox, Production | No |
| `Charges.GetRefundEligibility` | `GET /v1/charges/{id}/refund_eligibility` | `none` | `RefundEligibility` | 200 | Sandbox, Production | No |
| `Charges.Statistics` | `GET /v1/charges/statistics` | `GetChargeStatisticsParams` | `ChargeStatistics` | 200 | Sandbox, Production | No |
| `Customers.Create` | `POST /v1/customers` | `CreateCustomerParams` | `CustomerCreated` | 200, 201 | Sandbox, Production | Yes |
| `Customers.List` | `GET /v1/customers` | `ListCustomersParams` | `CustomerListResponse` | 200 | Sandbox, Production | No |
| `Customers.Get` | `GET /v1/customers/{id}` | `none` | `CustomerListItem` | 200 | Sandbox, Production | No |
| `Customers.Update` | `PATCH /v1/customers/{id}` | `UpdateCustomerParams` | `CustomerListItem` | 200 | Sandbox, Production | Yes |
| `Customers.Delete` | `DELETE /v1/customers/{id}` | `none` | `CustomerDeleted` | 200 | Sandbox, Production | No |
| `Customers.DeleteBatch` | `DELETE /v1/customers` | `RemoveCustomersParams` | `CustomerBatchDeleted` | 200 | Sandbox, Production | Yes |
| `Disputes.List` | `GET /v1/disputes` | `ListDisputesParams` | `DisputeListResponse` | 200 | Sandbox, Production | No |
| `Disputes.Get` | `GET /v1/disputes/{id}` | `none` | `Dispute` | 200 | Sandbox, Production | No |
| `PaymentLinks.Create` | `POST /v1/payment_links` | `CreatePaymentLinkParams` | `PaymentLink` | 200, 201 | Sandbox, Production | Yes |
| `PaymentLinks.List` | `GET /v1/payment_links` | `ListPaymentLinksParams` | `PaymentLinkListResponse` | 200 | Sandbox, Production | No |
| `PaymentLinks.Get` | `GET /v1/payment_links/{id}` | `none` | `PaymentLinkDetails` | 200 | Sandbox, Production | No |
| `PaymentLinks.Update` | `PATCH /v1/payment_links/{id}` | `UpdatePaymentLinkParams` | `PaymentLink` | 200 | Sandbox, Production | Yes |
| `PaymentLinks.SendEmail` | `POST /v1/payment_links/{id}/send_email` | `none` | `PaymentLinkEmailRequestAccepted` | 202 | Sandbox, Production | Yes |
| `Refunds.List` | `GET /v1/refunds` | `ListRefundsParams` | `RefundListResponse` | 200 | Sandbox, Production | No |
| `Refunds.Get` | `GET /v1/refunds/{id}` | `none` | `Refund` | 200 | Sandbox, Production | No |
| `Refunds.GetReceipt` | `GET /v1/refunds/{id}/receipt` | `none` | `FinancialDocument` | 200 | Sandbox, Production | No |
| `Statements.List` | `GET /v1/statements` | `ListStatementsParams` | `StatementListResponse` | 200 | Sandbox, Production | No |
| `Statements.Get` | `GET /v1/statements/{id}` | `none` | `Statement` | 200 | Sandbox, Production | No |
| `WebhookEventTypes.List` | `GET /v1/webhook_event_types` | `none` | `WebhookEventTypeListResponse` | 200 | Sandbox, Production | No |
| `Webhooks.Create` | `POST /v1/webhooks` | `CreateWebhookParams` | `WebhookCreated` | 201 | Sandbox, Production | Yes |
| `Webhooks.List` | `GET /v1/webhooks` | `ListWebhooksParams` | `WebhookListResponse` | 200 | Sandbox, Production | No |
| `Webhooks.Get` | `GET /v1/webhooks/{webhook_id}` | `none` | `Webhook` | 200 | Sandbox, Production | No |
| `Webhooks.Update` | `PATCH /v1/webhooks/{webhook_id}` | `UpdateWebhookParams` | `Webhook` | 200 | Sandbox, Production | Yes |
| `Webhooks.Delete` | `DELETE /v1/webhooks/{webhook_id}` | `none` | `Response` | 204 | Sandbox, Production | Yes |
| `Webhooks.GetHealth` | `GET /v1/webhooks/{webhook_id}/health` | `none` | `WebhookHealth` | 200 | Sandbox, Production | No |
| `Webhooks.SetPause` | `PATCH /v1/webhooks/{webhook_id}/pause` | `SetWebhookPauseParams` | `WebhookPause` | 200 | Sandbox, Production | Yes |
| `Webhooks.RotateSecret` | `POST /v1/webhooks/{webhook_id}/rotate_secret` | `RotateWebhookSecretParams` | `WebhookSecretRotation` | 200 | Sandbox, Production | Yes |
| `Webhooks.Test` | `POST /v1/webhooks/{webhook_id}/test` | `none` | `WebhookTestAccepted` | 202 | Sandbox, Production | Yes |
| `WithdrawalDestinations.List` | `GET /v1/withdrawal_destinations` | `ListWithdrawalDestinationsParams` | `WithdrawalDestinationListResponse` | 200 | Production | No |
| `WithdrawalDestinations.Get` | `GET /v1/withdrawal_destinations/{id}` | `none` | `WithdrawalDestination` | 200 | Production | No |
| `Withdrawals.Create` | `POST /v1/withdrawals` | `CreateWithdrawalParams` | `Withdrawal` | 200, 202 | Production | Yes |
| `Withdrawals.List` | `GET /v1/withdrawals` | `ListWithdrawalsParams` | `WithdrawalListResponse` | 200 | Sandbox, Production | No |
| `Withdrawals.Get` | `GET /v1/withdrawals/{id}` | `none` | `Withdrawal` | 200 | Sandbox, Production | No |
| `Withdrawals.GetLimits` | `GET /v1/withdrawals/limits` | `none` | `WithdrawalLimits` | 200 | Sandbox, Production | No |
| `Withdrawals.GetSummary` | `GET /v1/withdrawals/summary` | `none` | `WithdrawalSummary` | 200 | Production | No |
| `Withdrawals.GetAutomaticSettings` | `GET /v1/withdrawals/automatic_settings` | `none` | `AutomaticWithdrawalSettings` | 200 | Production | No |
| `Withdrawals.UpdateAutomaticSettings` | `PATCH /v1/withdrawals/automatic_settings` | `UpdateAutomaticWithdrawalSettingsParams` | `AutomaticWithdrawalSettings` | 200 | Production | Yes |
| `Simulations.Charges.ConfirmPayment` | `POST /v1/simulations/charges/{id}/confirm_payment` | `none` | `ChargeSimulation` | 200 | Sandbox | Yes |
| `Simulations.Charges.FailPayment` | `POST /v1/simulations/charges/{id}/fail_payment` | `none` | `ChargeSimulation` | 200 | Sandbox | Yes |
| `Simulations.Charges.ConfigureScenario` | `POST /v1/simulations/charges/{id}/scenario` | `ChargeScenarioParams` | `ChargeScenario` | 200 | Sandbox | Yes |
| `Simulations.Charges.ResolveScenario` | `POST /v1/simulations/charges/{id}/resolve` | `ResolveChargeScenarioParams` | `ChargeScenario` | 200 | Sandbox | Yes |
| `Simulations.Charges.IssuePayment` | `POST /v1/simulations/charges/{id}/issue_payment` | `IssueSimulationPaymentParams` | `Charge` | 200, 202 | Sandbox | Yes |
| `Simulations.Charges.CompensateBoleto` | `POST /v1/simulations/charges/{id}/compensate_boleto` | `none` | `ChargeBoletoCompensation` | 200 | Sandbox | Yes |
| `Simulations.Charges.CreateRefund` | `POST /v1/simulations/charges/{id}/refunds` | `CreateRefundParams` | `Refund` | 200, 202 | Sandbox | Yes |
| `Simulations.Refunds.Complete` | `POST /v1/simulations/refunds/{id}/complete` | `none` | `Refund` | 200 | Sandbox | Yes |
| `Simulations.Refunds.Fail` | `POST /v1/simulations/refunds/{id}/fail` | `none` | `Refund` | 200 | Sandbox | Yes |
