package webhook

import "log"

type WebhookService struct {
	lineClient *LineClient
}

func NewWebhookService(lineClient *LineClient) *WebhookService {
	return &WebhookService{
		lineClient: lineClient,
	}
}

func (s *WebhookService) HandleEvents(events []Event) error {
	for _, event := range events {
		log.Printf("Event type: %s, User ID: %s", event.Type, event.Source.UserID)

		switch event.Type {
		case "message":
			s.handleMessageEvent(event)
		// case "follow":
		// 	s.handleFollowEvent(event)
		// case "unfollow":
		// 	s.handleUnfollowEvent(event)
		default:
			log.Printf("Unhandled event type: %s", event.Type)
		}
	}
	return nil
}

// handleMessageEvent - 메시지 이벤트 처리
func (s *WebhookService) handleMessageEvent(event Event) {
	if event.Message != nil && event.Message.Type == "text" {
		log.Printf("Received message: %s", event.Message.Text)

		// 메시지 답장
		replyText := "받은 메시지: " + event.Message.Text
		if err := s.lineClient.ReplyMessage(event.ReplyToken, replyText); err != nil {
			log.Printf("Failed to reply message: %v", err)
		}
	}
}

// // handleFollowEvent - 팔로우 이벤트 처리
// func (s *WebhookService) handleFollowEvent(event Event) {
// 	log.Printf("New follower: %s", event.Source.UserID)

// 	if err := s.lineClient.ReplyMessage(event.ReplyToken, "친구 추가 감사합니다! 👋"); err != nil {
// 		log.Printf("Failed to reply follow event: %v", err)
// 	}
// }

// // handleUnfollowEvent - 언팔로우 이벤트 처리
// func (s *WebhookService) handleUnfollowEvent(event Event) {
// 	log.Printf("User unfollowed: %s", event.Source.UserID)
// 	// unfollow는 replyToken이 없으므로 답장 불가
// }
