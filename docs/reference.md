# Referência do SDK Go
Gerada do contrato público fixado. `Params` usa `Set(valor)` para enviar um campo e `Null[T]()` para enviar null. Todos os métodos recebem `context.Context`, os IDs indicados no caminho e opções de requisição opcionais. Consulte os campos dos tipos com `go doc`.
| Serviço e método | Endpoint | Parâmetros | Retorno | Status | Ambientes OAuth | Idempotência |
| --- | --- | --- | --- | --- | --- | --- |
| `Accounts.GetFees` | `GET /v1/accounts/fees` | `nenhum` | `AccountFees` | 200 | Sandbox, Production | Não |
| `Accounts.GetNumber` | `GET /v1/accounts/account_number` | `nenhum` | `AccountNumber` | 200 | Production | Não |
| `Accounts.GetCheckoutSettings` | `GET /v1/accounts/checkout_settings` | `nenhum` | `CheckoutSettings` | 200 | Sandbox, Production | Não |
| `Accounts.UpdateCheckoutSettings` | `PATCH /v1/accounts/checkout_settings` | `UpdateCheckoutSettingsParams` | `CheckoutSettings` | 200, 202 | Sandbox, Production | Sim |
| `Accounts.CreateCheckoutLogoUpload` | `POST /v1/accounts/checkout_settings/logo_uploads` | `CreateCheckoutLogoUploadParams` | `CheckoutAssetUpload` | 201 | Sandbox, Production | Sim |
| `Accounts.CreateCheckoutFaviconUpload` | `POST /v1/accounts/checkout_settings/favicon_uploads` | `CreateCheckoutFaviconUploadParams` | `CheckoutAssetUpload` | 201 | Sandbox, Production | Sim |
| `Accounts.DeleteCheckoutAssetUpload` | `DELETE /v1/accounts/checkout_settings/uploads/{id}` | `nenhum` | `CheckoutAssetUploadDeleted` | 200 | Sandbox, Production | Sim |
| `Accounts.GetContext` | `GET /v1/accounts/context` | `nenhum` | `AccountContext` | 200 | Sandbox, Production | Não |
| `Balances.Get` | `GET /v1/balances` | `nenhum` | `Balance` | 200 | Sandbox, Production | Não |
| `Balances.Forecast` | `GET /v1/balances/forecast` | `GetBalanceForecastParams` | `BalanceForecast` | 200 | Sandbox, Production | Não |
| `Charges.Create` | `POST /v1/charges` | `CreateChargeParams` | `ChargeCreated` | 201, 202 | Sandbox, Production | Sim |
| `Charges.List` | `GET /v1/charges` | `ListChargesParams` | `ChargeListResponse` | 200 | Sandbox, Production | Não |
| `Charges.Get` | `GET /v1/charges/{id}` | `nenhum` | `Charge` | 200 | Sandbox, Production | Não |
| `Charges.Update` | `PATCH /v1/charges/{id}` | `UpdateChargeParams` | `Charge` | 200 | Sandbox, Production | Sim |
| `Charges.Revoke` | `DELETE /v1/charges/{id}` | `nenhum` | `RevokeChargeResult` | 200, 202 | Sandbox, Production | Sim |
| `Charges.GetPayment` | `GET /v1/charges/{id}/payment` | `nenhum` | `PaymentReceipt` | 200 | Sandbox, Production | Não |
| `Charges.GetReceipt` | `GET /v1/charges/{id}/receipt` | `nenhum` | `FinancialDocument` | 200 | Sandbox, Production | Não |
| `Charges.GetBoleto` | `GET /v1/charges/{id}/boleto` | `nenhum` | `FinancialDocument` | 200 | Sandbox, Production | Não |
| `Charges.CreateRefund` | `POST /v1/charges/{id}/refunds` | `CreateRefundParams` | `Refund` | 200, 202 | Sandbox, Production | Sim |
| `Charges.ListRefunds` | `GET /v1/charges/{id}/refunds` | `ListChargeRefundsParams` | `RefundListResponse` | 200 | Sandbox, Production | Não |
| `Charges.GetRefundEligibility` | `GET /v1/charges/{id}/refund_eligibility` | `nenhum` | `RefundEligibility` | 200 | Sandbox, Production | Não |
| `Charges.Statistics` | `GET /v1/charges/statistics` | `GetChargeStatisticsParams` | `ChargeStatistics` | 200 | Sandbox, Production | Não |
| `Customers.Create` | `POST /v1/customers` | `CreateCustomerParams` | `CustomerCreated` | 200, 201 | Sandbox, Production | Sim |
| `Customers.List` | `GET /v1/customers` | `ListCustomersParams` | `CustomerListResponse` | 200 | Sandbox, Production | Não |
| `Customers.Get` | `GET /v1/customers/{id}` | `nenhum` | `CustomerListItem` | 200 | Sandbox, Production | Não |
| `Customers.Update` | `PATCH /v1/customers/{id}` | `UpdateCustomerParams` | `CustomerListItem` | 200 | Sandbox, Production | Sim |
| `Customers.Delete` | `DELETE /v1/customers/{id}` | `nenhum` | `CustomerDeleted` | 200 | Sandbox, Production | Não |
| `Customers.DeleteBatch` | `DELETE /v1/customers` | `RemoveCustomersParams` | `CustomerBatchDeleted` | 200 | Sandbox, Production | Sim |
| `Disputes.List` | `GET /v1/disputes` | `ListDisputesParams` | `DisputeListResponse` | 200 | Sandbox, Production | Não |
| `Disputes.Get` | `GET /v1/disputes/{id}` | `nenhum` | `Dispute` | 200 | Sandbox, Production | Não |
| `PaymentLinks.Create` | `POST /v1/payment_links` | `CreatePaymentLinkParams` | `PaymentLink` | 200, 201 | Sandbox, Production | Sim |
| `PaymentLinks.List` | `GET /v1/payment_links` | `ListPaymentLinksParams` | `PaymentLinkListResponse` | 200 | Sandbox, Production | Não |
| `PaymentLinks.Get` | `GET /v1/payment_links/{id}` | `nenhum` | `PaymentLinkDetails` | 200 | Sandbox, Production | Não |
| `PaymentLinks.Update` | `PATCH /v1/payment_links/{id}` | `UpdatePaymentLinkParams` | `PaymentLink` | 200 | Sandbox, Production | Sim |
| `PaymentLinks.SendEmail` | `POST /v1/payment_links/{id}/send_email` | `nenhum` | `PaymentLinkEmailRequestAccepted` | 202 | Sandbox, Production | Sim |
| `Refunds.List` | `GET /v1/refunds` | `ListRefundsParams` | `RefundListResponse` | 200 | Sandbox, Production | Não |
| `Refunds.Get` | `GET /v1/refunds/{id}` | `nenhum` | `Refund` | 200 | Sandbox, Production | Não |
| `Refunds.GetReceipt` | `GET /v1/refunds/{id}/receipt` | `nenhum` | `FinancialDocument` | 200 | Sandbox, Production | Não |
| `Statements.List` | `GET /v1/statements` | `ListStatementsParams` | `StatementListResponse` | 200 | Sandbox, Production | Não |
| `Statements.Get` | `GET /v1/statements/{id}` | `nenhum` | `Statement` | 200 | Sandbox, Production | Não |
| `WebhookEventTypes.List` | `GET /v1/webhook_event_types` | `nenhum` | `WebhookEventTypeListResponse` | 200 | Sandbox, Production | Não |
| `Webhooks.Create` | `POST /v1/webhooks` | `CreateWebhookParams` | `WebhookCreated` | 201 | Sandbox, Production | Sim |
| `Webhooks.List` | `GET /v1/webhooks` | `ListWebhooksParams` | `WebhookListResponse` | 200 | Sandbox, Production | Não |
| `Webhooks.Get` | `GET /v1/webhooks/{webhook_id}` | `nenhum` | `Webhook` | 200 | Sandbox, Production | Não |
| `Webhooks.Update` | `PATCH /v1/webhooks/{webhook_id}` | `UpdateWebhookParams` | `Webhook` | 200 | Sandbox, Production | Sim |
| `Webhooks.Delete` | `DELETE /v1/webhooks/{webhook_id}` | `nenhum` | `Response` | 204 | Sandbox, Production | Sim |
| `Webhooks.GetHealth` | `GET /v1/webhooks/{webhook_id}/health` | `nenhum` | `WebhookHealth` | 200 | Sandbox, Production | Não |
| `Webhooks.SetPause` | `PATCH /v1/webhooks/{webhook_id}/pause` | `SetWebhookPauseParams` | `WebhookPause` | 200 | Sandbox, Production | Sim |
| `Webhooks.RotateSecret` | `POST /v1/webhooks/{webhook_id}/rotate_secret` | `RotateWebhookSecretParams` | `WebhookSecretRotation` | 200 | Sandbox, Production | Sim |
| `Webhooks.Test` | `POST /v1/webhooks/{webhook_id}/test` | `nenhum` | `WebhookTestAccepted` | 202 | Sandbox, Production | Sim |
| `WithdrawalDestinations.List` | `GET /v1/withdrawal_destinations` | `ListWithdrawalDestinationsParams` | `WithdrawalDestinationListResponse` | 200 | Production | Não |
| `WithdrawalDestinations.Get` | `GET /v1/withdrawal_destinations/{id}` | `nenhum` | `WithdrawalDestination` | 200 | Production | Não |
| `Withdrawals.Create` | `POST /v1/withdrawals` | `CreateWithdrawalParams` | `Withdrawal` | 200, 202 | Production | Sim |
| `Withdrawals.List` | `GET /v1/withdrawals` | `ListWithdrawalsParams` | `WithdrawalListResponse` | 200 | Sandbox, Production | Não |
| `Withdrawals.Get` | `GET /v1/withdrawals/{id}` | `nenhum` | `Withdrawal` | 200 | Sandbox, Production | Não |
| `Withdrawals.GetLimits` | `GET /v1/withdrawals/limits` | `nenhum` | `WithdrawalLimits` | 200 | Sandbox, Production | Não |
| `Withdrawals.GetSummary` | `GET /v1/withdrawals/summary` | `nenhum` | `WithdrawalSummary` | 200 | Production | Não |
| `Withdrawals.GetAutomaticSettings` | `GET /v1/withdrawals/automatic_settings` | `nenhum` | `AutomaticWithdrawalSettings` | 200 | Production | Não |
| `Withdrawals.UpdateAutomaticSettings` | `PATCH /v1/withdrawals/automatic_settings` | `UpdateAutomaticWithdrawalSettingsParams` | `AutomaticWithdrawalSettings` | 200 | Production | Sim |
| `Simulations.Charges.ConfirmPayment` | `POST /v1/simulations/charges/{id}/confirm_payment` | `nenhum` | `ChargeSimulation` | 200 | Sandbox | Sim |
| `Simulations.Charges.FailPayment` | `POST /v1/simulations/charges/{id}/fail_payment` | `nenhum` | `ChargeSimulation` | 200 | Sandbox | Sim |
| `Simulations.Charges.ConfigureScenario` | `POST /v1/simulations/charges/{id}/scenario` | `ChargeScenarioParams` | `ChargeScenario` | 200 | Sandbox | Sim |
| `Simulations.Charges.ResolveScenario` | `POST /v1/simulations/charges/{id}/resolve` | `ResolveChargeScenarioParams` | `ChargeScenario` | 200 | Sandbox | Sim |
| `Simulations.Charges.IssuePayment` | `POST /v1/simulations/charges/{id}/issue_payment` | `IssueSimulationPaymentParams` | `Charge` | 200, 202 | Sandbox | Sim |
| `Simulations.Charges.CompensateBoleto` | `POST /v1/simulations/charges/{id}/compensate_boleto` | `nenhum` | `ChargeBoletoCompensation` | 200 | Sandbox | Sim |
| `Simulations.Charges.CreateRefund` | `POST /v1/simulations/charges/{id}/refunds` | `CreateRefundParams` | `Refund` | 200, 202 | Sandbox | Sim |
| `Simulations.Refunds.Complete` | `POST /v1/simulations/refunds/{id}/complete` | `nenhum` | `Refund` | 200 | Sandbox | Sim |
| `Simulations.Refunds.Fail` | `POST /v1/simulations/refunds/{id}/fail` | `nenhum` | `Refund` | 200 | Sandbox | Sim |
