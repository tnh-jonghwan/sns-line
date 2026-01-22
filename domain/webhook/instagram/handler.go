package instagram

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
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

// verifySignature - 웹훅 요청 서명 검증
func (h *InstagramHandler) verifySignature(c *fiber.Ctx) error {
	signature := c.Get("x-hub-signature")
	if signature == "" {
		log.Println("Missing x-hub-signature header")
		return fiber.NewError(fiber.StatusUnauthorized, "Missing signature")
	}

	// "sha1=hash" 형식에서 hash 추출
	parts := strings.Split(signature, "=")
	if len(parts) != 2 || parts[0] != "sha1" {
		log.Println("Invalid signature format")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid signature format")
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
		return fiber.NewError(fiber.StatusUnauthorized, "Signature mismatch")
	}

	log.Println("Signature verified successfully")
	return nil
}

// shouldSkipEvent - 불필요한 이벤트인지 확인
func (h *InstagramHandler) shouldSkipEvent(messaging *Messaging) bool {
	if messaging.Read != nil {
		log.Println("Got a read event - skipping")
		return true
	}
	if messaging.Delivery != nil {
		log.Println("Got a delivery event - skipping")
		return true
	}
	if messaging.Message != nil && messaging.Message.IsEcho {
		log.Printf("Got an echo of our send, mid = %s - skipping", messaging.Message.Mid)
		return true
	}
	return false
}

// handleMessage - 사용자 메시지 처리
func (h *InstagramHandler) handleMessage(messaging *Messaging) {
	if messaging.Message == nil || messaging.Message.Text == "" {
		return
	}

	senderID := messaging.Sender.ID
	messageText := messaging.Message.Text

	log.Printf("📩 Instagram message from %s: %s", senderID, messageText)

	// SSE로 브로드캐스트
	h.eventHub.Broadcast(messageText, senderID)
}

// handlePostback - Postback 이벤트 처리
func (h *InstagramHandler) handlePostback(messaging *Messaging) {
	if messaging.Postback == nil {
		return
	}

	senderID := messaging.Sender.ID
	payload := messaging.Postback.Payload

	log.Printf("🔘 Instagram postback from %s: %s", senderID, payload)

	// Postback도 브로드캐스트
	h.eventHub.Broadcast(payload, senderID)
}

// HandleWebhook - POST 웹훅 수신 핸들러
func (h *InstagramHandler) HandleWebhook(c *fiber.Ctx) error {
	// 서명 검증
	if err := h.verifySignature(c); err != nil {
		return err
	}

	// RAW Body 로그 출력 (디버깅용)
	bodyBytes := c.Body()

	// 요청 파싱
	var req WebhookRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		log.Printf("Error parsing Instagram webhook: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// log.Printf("Instagram webhook received: %+v", req)

	// Instagram (또는 page) 이벤트인지 확인
	if req.Object != "instagram" && req.Object != "page" {
		log.Printf("Unsupported object type: %s", req.Object)
		return c.SendStatus(fiber.StatusNotFound)
	}

	// 각 엔트리 처리 (배치로 여러 개 올 수 있음)
	for _, entry := range req.Entry {
		// Messaging 처리 (DM)
		for _, messaging := range entry.Messaging {
			// 불필요한 이벤트 필터링
			if h.shouldSkipEvent(&messaging) {
				continue
			}

			// 메시지 처리
			h.handleMessage(&messaging)

			// Postback 처리
			h.handlePostback(&messaging)
		}

		// Changes 처리 (댓글, 좋아요 등)
		for _, change := range entry.Changes {
			log.Printf("📨 Instagram change event: field=%s", change.Field)

			switch change.Field {
			case "messages":
				// DM 메시지 (changes로 올 수도 있음)
				log.Printf("💬 Instagram DM from %s: %s", change.Value.From.ID, change.Value.Text)
				// EventHub로 브로드캐스트
				h.eventHub.Broadcast(change.Value.Text, change.Value.From.ID)
			case "comments":
				log.Printf("💬 Comment from %s: %s", change.Value.From.Username, change.Value.Text)
				// 댓글 이벤트도 EventHub로 브로드캐스트
				h.eventHub.Broadcast(change.Value.Text, change.Value.From.ID)
			case "mentions":
				log.Printf("@️⃣ Mention event")
				// 멘션 이벤트 처리
			case "feed":
				log.Printf("📰 Feed event")
				// 피드 이벤트 처리
			default:
				log.Printf("🔔 Other change event: %s", change.Field)
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
