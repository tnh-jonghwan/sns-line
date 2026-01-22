package instagram

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"sns-line/config"
	"sns-line/domain/eventHub"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type InstagramHandler struct {
	verifyToken string
	appSecret   string
	eventHub    *eventHub.EventHub
}

func NewInstagramHandler(env *config.Env, eventHub *eventHub.EventHub) *InstagramHandler {
	return &InstagramHandler{
		verifyToken: env.InstagramVerifyToken,
		appSecret:   env.InstagramAppSecret,
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
	// 서명 검증 (Facebook 예시 코드의 verifyRequestSignature와 동일)
	signature := c.Get("x-hub-signature")
	if signature == "" {
		log.Println("Missing x-hub-signature header")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// "sha1=hash" 형식에서 hash 추출
	parts := strings.Split(signature, "=")
	if len(parts) != 2 || parts[0] != "sha1" {
		log.Println("Invalid signature format")
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	signatureHash := parts[1]

	// HMAC-SHA1으로 예상 해시 생성
	body := c.Body()
	mac := hmac.New(sha1.New, []byte(h.appSecret))
	mac.Write(body)
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	// 서명 비교
	if signatureHash != expectedHash {
		log.Printf("Signature verification failed: got %s, expected %s", signatureHash, expectedHash)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	log.Println("Signature verified successfully")

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
