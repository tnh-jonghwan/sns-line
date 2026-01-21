package webhook

import (
	"log"
	"sns-line/domain/eventHub"
	"sns-line/domain/line"
)

type WebhookService struct {
	lineClient *line.LineClient
}

func NewWebhookService(lineClient *line.LineClient) *WebhookService {
	return &WebhookService{
		lineClient: lineClient,
	}
}

func (s *WebhookService) HandleEvents(events []Event, eventHub *eventHub.EventHub) error {
	for _, event := range events {
		log.Printf("Event type: %s, User ID: %s", event.Type, event.Source.UserID)

		switch event.Type {
		case "message":
			s.handleMessageEvent(event, eventHub)
		default:
			log.Printf("Unhandled event type: %s", event.Type)
		}
	}
	return nil
}

// handleMessageEvent - 메시지 이벤트 처리
func (s *WebhookService) handleMessageEvent(event Event, eventHub *eventHub.EventHub) {
	if event.Message != nil && event.Message.Type == "text" {
		userMessage := event.Message.Text
		log.Printf("Received message: %s", userMessage)

		// eventHub로 브로드캐스트
		eventHub.Broadcast(userMessage, event.Source.UserID)

		// 메시지 답장
		if err := s.lineClient.ReplyMessage(event.ReplyToken, "받은 메시지: "+userMessage); err != nil {
			log.Printf("Failed to reply message: %v", err)
		}
	}
}
