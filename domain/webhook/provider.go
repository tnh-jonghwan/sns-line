package webhook

import (
	"messaging-line/config"

	"go.uber.org/fx"
)

// ProvideLineClient - LINE API 클라이언트 Provider
func ProvideLineClient(accessToken string, env *config.Env) *LineClient {
	return NewLineClient(accessToken, env.LineApiPrefix)
}

// ProvideWebhookService - Webhook 서비스 Provider
func ProvideWebhookService(client *LineClient) *WebhookService {
	return NewWebhookService(client)
}

// ProvideWebhookHandler - Webhook 핸들러 Provider
func ProvideWebhookHandler(service *WebhookService) *WebhookHandler {
	return NewWebhookHandler(service)
}

// Module - Webhook 모듈
var Module = fx.Options(
	fx.Provide(
		ProvideLineClient,
		ProvideWebhookService,
		ProvideWebhookHandler,
	),
)
