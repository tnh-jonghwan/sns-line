package eventHub

import (
	"bufio"
	"fmt"
	"log"
	"sns/domain/line"

	"github.com/gofiber/fiber/v2"
)

func RegistereventHubRoutes(app *fiber.App, eventHub *EventHub, lineClient *line.LineClient) {
	// eventHub 엔드포인트
	app.Get("/events", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("X-Accel-Buffering", "no")

		// 새 클라이언트 채널 생성
		client := make(chan string, 10)
		eventHub.AddClient(client)

		log.Println("eventHub 클라이언트 연결됨")

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			defer func() {
				eventHub.RemoveClient(client)
				log.Println("eventHub 클라이언트 연결 해제됨")
			}()

			// 초기 연결 메시지
			fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\"}\n\n")
			w.Flush()

			// 메시지 수신 루프
			for message := range client {
				fmt.Fprintf(w, "data: %s\n\n", message)
				if err := w.Flush(); err != nil {
					log.Printf("eventHub flush error: %v", err)
					return
				}
			}
		})

		return nil
	})

	// 메시지 전송 API
	app.Post("/api/send", func(c *fiber.Ctx) error {
		var req struct {
			Text string `json:"text"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid request",
			})
		}

		log.Printf("메시지 전송 요청: %s", req.Text)

		// LINE Broadcast API로 메시지 전송
		if err := lineClient.BroadcastMessage(req.Text); err != nil {
			log.Printf("Broadcast 실패: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
		})
	})
}
