package instagram

import (
	"log"
	"sns-line/config"
	"sns-line/domain/eventHub"

	"github.com/gofiber/fiber/v2"
)

type InstagramHandler struct {
	verifyToken string
	eventHub    *eventHub.EventHub
}

func NewInstagramHandler(env *config.Env, eventHub *eventHub.EventHub) *InstagramHandler {
	return &InstagramHandler{
		verifyToken: env.InstagramVerifyToken,
		eventHub:    eventHub,
	}
}

// HandleVerify - GET 웹훅 검증 핸들러
func (h *InstagramHandler) HandleVerify(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	log.Printf("Instagram webhook verify: mode=%s, token=%s", mode, token)

	if mode == "subscribe" && token == h.verifyToken {
		log.Println("Instagram webhook verified!")
		return c.SendString(challenge)
	}

	log.Println("Instagram webhook verification failed")
	return c.SendStatus(fiber.StatusBadRequest)
}

// HandleWebhook - POST 웹훅 수신 핸들러
func (h *InstagramHandler) HandleWebhook(c *fiber.Ctx) error {
	var req WebhookRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing Instagram webhook: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	log.Printf("Instagram webhook received: %+v", req)

	// 메시지 처리
	for _, entry := range req.Entry {
		for _, messaging := range entry.Messaging {
			if messaging.Message != nil && messaging.Message.Text != "" {
				log.Printf("Instagram message from %s: %s", messaging.Sender.ID, messaging.Message.Text)

				// SSE로 브로드캐스트
				h.eventHub.Broadcast(messaging.Message.Text, messaging.Sender.ID)
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
