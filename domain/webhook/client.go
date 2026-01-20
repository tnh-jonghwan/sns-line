package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// LineClient - LINE API와 통신하는 클라이언트
type LineClient struct {
	accessToken string
	apiURL      string
}

// NewLineClient - LineClient 생성자
func NewLineClient(accessToken, apiURL string) *LineClient {
	return &LineClient{
		accessToken: accessToken,
		apiURL:      apiURL,
	}
}

// ReplyMessage - LINE 메시지 답장 API 호출
func (c *LineClient) ReplyMessage(replyToken, text string) error {
	url := fmt.Sprintf("%s/v2/bot/message/reply", c.apiURL)

	replyData := ReplyRequest{
		ReplyToken: replyToken,
		Messages: []ReplyMessage{
			{
				Type: "text",
				Text: text,
			},
		},
	}

	jsonData, err := json.Marshal(replyData)
	if err != nil {
		return fmt.Errorf("failed to marshal reply data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Reply API error: %s", string(body))
		return fmt.Errorf("reply failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Reply sent successfully")
	return nil
}
