package webhook

import (
	"messaging-line/config"

	"go.uber.org/fx"
)

// ProvideWebhookHandler - Webhook 핸들러 생성
func ProvideWebhookHandler(accessToken string, env *config.Env) *WebhookHandler {
	lineClient := NewLineClient(accessToken, env.LineApiPrefix)
	webhookService := NewWebhookService(lineClient)
	return NewWebhookHandler(webhookService)
}

// Module - Webhook 도메인 모듈
var Module = fx.Module("webhook",
	fx.Provide(ProvideWebhookHandler),
)
