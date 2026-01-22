package line

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sns/config"
)

type LineClient struct {
	LineAccessToken string
	LineApiURL      string
}

func NewLineClient(env *config.Env) *LineClient {
	return &LineClient{
		LineAccessToken: env.LineAccessToken,
		LineApiURL:      env.LineApiPrefix,
	}
}

// ReplyMessage - LINE 메시지 답장 API 호출 (단일 메시지)
func (c *LineClient) ReplyMessage(replyToken, text string) error {
	return c.ReplyMessages(replyToken, []string{text})
}

// ReplyMessages - LINE 메시지 답장 API 호출 (복수 메시지, 최대 5개)
func (c *LineClient) ReplyMessages(replyToken string, texts []string) error {
	url := fmt.Sprintf("%s/v2/bot/message/reply", c.LineApiURL)

	// 최대 5개 메시지만 허용
	if len(texts) > 5 {
		return fmt.Errorf("LINE API allows maximum 5 messages per reply, got %d", len(texts))
	}

	// 메시지 배열 생성
	messages := make([]ReplyMessage, len(texts))
	for i, text := range texts {
		messages[i] = ReplyMessage{
			Type: "text",
			Text: text,
		}
	}

	replyData := ReplyRequest{
		ReplyToken: replyToken,
		Messages:   messages,
	}

	jsonData, err := json.Marshal(replyData)
	if err != nil {
		return fmt.Errorf("failed to marshal reply data: %w", err)
	}

	log.Printf("Sending reply to LINE API: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.LineAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Printf("Reply API error (status %d): %s", resp.StatusCode, string(body))
		return fmt.Errorf("reply failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Println("Reply sent successfully")
	return nil
}

// BroadcastMessage - LINE 브로드캐스트 API 호출 (모든 친구에게 메시지 전송)
func (c *LineClient) BroadcastMessage(text string) error {
	url := fmt.Sprintf("%s/v2/bot/message/broadcast", c.LineApiURL)

	// 브로드캐스트 메시지 생성
	broadcastData := map[string]interface{}{
		"messages": []map[string]string{
			{
				"type": "text",
				"text": text,
			},
		},
	}

	jsonData, err := json.Marshal(broadcastData)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast data: %w", err)
	}

	log.Printf("Sending broadcast to LINE API: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.LineAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Printf("Broadcast API error (status %d): %s", resp.StatusCode, string(body))
		return fmt.Errorf("broadcast failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Println("Broadcast sent successfully")
	return nil
}
