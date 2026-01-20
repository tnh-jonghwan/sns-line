package main

import (
	"fmt"
	"messaging-line/jwt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Access Token 생성
	accessToken := jwt.GenerateAccessToken()
	fmt.Println("Generated Access Token:", accessToken)

	app := fiber.New()

	// 기본 라우트
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "LINE Webhook Server",
		})
	})

	// LINE Webhook 엔드포인트
	app.Post("/webhook", func(c *fiber.Ctx) error {
		// TODO: 1. Request Body 파싱
		// - LINE에서 보내는 JSON 데이터를 구조체로 파싱
		// - 구조체: destination, events[] 포함

		// TODO: 2. Events 순회하며 처리
		// for _, event := range events:
		//   - event.type 확인 (message, follow, unfollow, join, leave 등)
		//   - event.source.userId 확인 (사용자 정보)
		//   - event.replyToken 확인 (답장용 토큰)

		// TODO: 3. Message 이벤트 처리
		// if event.type == "message":
		//   - event.message.type 확인 (text, image, video 등)
		//   - event.message.text 확인 (텍스트 메시지인 경우)
		//   - 로직 처리 (DB 저장, 분석 등)

		// TODO: 4. 메시지 답장 (선택사항)
		// - POST https://api.line.me/v2/bot/message/reply
		// - Header: Authorization: Bearer {access_token}
		// - Body: { replyToken, messages[] }

		// TODO: 5. 200 OK 응답 (필수!)
		// - LINE은 반드시 200 응답 받아야 함
		// - 그렇지 않으면 재시도함
		return c.SendStatus(fiber.StatusOK)
	})

	// TODO: Access Token 갱신 API (선택사항)
	// JWT는 30분마다 만료되므로, 주기적으로 새로운 access token 발급 필요
	app.Get("/refresh-token", func(c *fiber.Ctx) error {
		// TODO: 새로운 JWT 생성 → access token 발급
		return c.JSON(fiber.Map{
			"message": "Token refreshed",
		})
	})

	// 서버 시작
	app.Listen(":3000")
}
