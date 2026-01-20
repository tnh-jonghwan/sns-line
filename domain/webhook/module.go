package webhook

import "go.uber.org/fx"

// WebhookModule - Webhook 관련 모든 의존성을 제공하는 FX 모듈
// LineClient -> WebhookService -> WebhookHandler 의존성 체인을 자동으로 구성
var WebhookModule = fx.Module("webhook",
	fx.Provide(
		NewLineClient,     // LINE API 클라이언트 생성
		NewWebhookService, // Webhook 비즈니스 로직 서비스 생성
		NewWebhookHandler, // HTTP 핸들러 생성
	),
)
