// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"messaging-line/jwt"
// 	"net/http"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// )

// // LINE Webhook 요청 구조체
// type WebhookRequest struct {
// 	Destination string  `json:"destination"`
// 	Events      []Event `json:"events"`
// }

// type Event struct {
// 	Type       string   `json:"type"`
// 	Timestamp  int64    `json:"timestamp"`
// 	Source     Source   `json:"source"`
// 	ReplyToken string   `json:"replyToken,omitempty"`
// 	Message    *Message `json:"message,omitempty"`
// }

// type Source struct {
// 	Type   string `json:"type"`
// 	UserID string `json:"userId,omitempty"`
// }

// type Message struct {
// 	ID   string `json:"id"`
// 	Type string `json:"type"`
// 	Text string `json:"text,omitempty"`
// }

// // Reply API 요청 구조체
// type ReplyRequest struct {
// 	ReplyToken string         `json:"replyToken"`
// 	Messages   []ReplyMessage `json:"messages"`
// }

// type ReplyMessage struct {
// 	Type string `json:"type"`
// 	Text string `json:"text"`
// }

// // Access Token 응답 구조체
// type TokenResponse struct {
// 	AccessToken string `json:"access_token"`
// 	TokenType   string `json:"token_type"`
// 	ExpiresIn   int    `json:"expires_in"`
// 	KeyID       string `json:"key_id"`
// }

// var accessToken string

// func main_answer() {
// 	// .env 파일 로드
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// JWT 토큰 생성
// 	jwtToken := jwt.GetJWT()
// 	fmt.Println("Generated JWT:", jwtToken)

// 	// Access Token 발급
// 	accessToken = getAccessToken(jwtToken)
// 	fmt.Println("Access Token:", accessToken)

// 	app := fiber.New()

// 	// 기본 라우트 (테스트용)
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{
// 			"message": "LINE Webhook Server",
// 			"status":  "running",
// 		})
// 	})

// 	// LINE Webhook 엔드포인트
// 	app.Post("/webhook", webhookHandler)

// 	// Access Token 갱신 API
// 	app.Get("/refresh-token", func(c *fiber.Ctx) error {
// 		jwtToken := jwt.GetJWT()
// 		accessToken = getAccessToken(jwtToken)
// 		return c.JSON(fiber.Map{
// 			"message":      "Token refreshed",
// 			"access_token": accessToken,
// 		})
// 	})

// 	// 서버 시작
// 	log.Fatal(app.Listen(":3000"))
// }

// // Access Token 발급 함수
// func getAccessToken(jwtToken string) string {
// 	url := "https://api.line.me/oauth2/v2.1/token"

// 	// Form data 생성
// 	data := map[string]string{
// 		"grant_type":            "client_credentials",
// 		"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
// 		"client_assertion":      jwtToken,
// 	}

// 	// JSON body 생성
// 	jsonData, _ := json.Marshal(data)

// 	// HTTP 요청
// 	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("Failed to get access token: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// 응답 파싱
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	var tokenResp TokenResponse
// 	json.Unmarshal(body, &tokenResp)

// 	return tokenResp.AccessToken
// }

// // Webhook 핸들러
// func webhookHandler(c *fiber.Ctx) error {
// 	// 1. Request Body 파싱
// 	var req WebhookRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		log.Printf("Error parsing request: %v", err)
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request",
// 		})
// 	}

// 	// 2. Events 순회하며 처리
// 	for _, event := range req.Events {
// 		log.Printf("Event type: %s, User ID: %s", event.Type, event.Source.UserID)

// 		// 3. Message 이벤트 처리
// 		if event.Type == "message" && event.Message != nil {
// 			if event.Message.Type == "text" {
// 				log.Printf("Received message: %s", event.Message.Text)

// 				// 4. 메시지 답장
// 				replyMessage(event.ReplyToken, "받은 메시지: "+event.Message.Text)
// 			}
// 		}

// 		// Follow 이벤트 처리
// 		if event.Type == "follow" {
// 			log.Printf("New follower: %s", event.Source.UserID)
// 			replyMessage(event.ReplyToken, "친구 추가 감사합니다! 👋")
// 		}

// 		// Unfollow 이벤트 처리
// 		if event.Type == "unfollow" {
// 			log.Printf("User unfollowed: %s", event.Source.UserID)
// 			// unfollow는 replyToken이 없음
// 		}
// 	}

// 	// 5. 200 OK 응답
// 	return c.SendStatus(fiber.StatusOK)
// }

// // 메시지 답장 함수
// func replyMessage(replyToken, text string) error {
// 	url := "https://api.line.me/v2/bot/message/reply"

// 	replyData := ReplyRequest{
// 		ReplyToken: replyToken,
// 		Messages: []ReplyMessage{
// 			{
// 				Type: "text",
// 				Text: text,
// 			},
// 		},
// 	}

// 	jsonData, _ := json.Marshal(replyData)

// 	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+accessToken)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("Failed to send reply: %v", err)
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		log.Printf("Reply API error: %s", string(body))
// 		return fmt.Errorf("reply failed with status %d", resp.StatusCode)
// 	}

// 	log.Printf("Reply sent successfully")
// 	return nil
// }
