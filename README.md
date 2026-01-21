# LINE Webhook Reply 사용법

## 개요

LINE Messaging API의 reply 기능을 구현했습니다. LINE API 스펙에 맞춰 최대 5개의 메시지를 한 번에 보낼 수 있습니다.

---

## API 구조

### LINE API 스펙

```bash
curl -X POST https://api.line.me/v2/bot/message/reply \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer {channel access token}' \
-d '{
    "replyToken":"nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
    "messages":[
        {
            "type":"text",
            "text":"Hello, user"
        },
        {
            "type":"text",
            "text":"May I help you?"
        }
    ]
}'
```

---

## 구현된 함수

### 1. ReplyMessage (단일 메시지)

```go
func (c *LineClient) ReplyMessage(replyToken, text string) error
```

**사용 예시:**
```go
err := lineClient.ReplyMessage(event.ReplyToken, "안녕하세요!")
```

---

### 2. ReplyMessages (복수 메시지, 최대 5개)

```go
func (c *LineClient) ReplyMessages(replyToken string, texts []string) error
```

**사용 예시:**
```go
messages := []string{
    "안녕하세요!",
    "무엇을 도와드릴까요?",
}
err := lineClient.ReplyMessages(event.ReplyToken, messages)
```

---

## WebhookService에서 사용하는 방법

### 현재 구조 (service.go)

```go
func (s *WebhookService) handleMessageEvent(event Event) {
    if event.Message == nil || event.Message.Type != "text" {
        return
    }

    userMessage := event.Message.Text
    log.Printf("User said: %s", userMessage)

    // 단일 메시지 답장
    if err := s.lineClient.ReplyMessage(event.ReplyToken, "메시지를 받았습니다: "+userMessage); err != nil {
        log.Printf("Failed to reply: %v", err)
    }
}
```

### 복수 메시지 답장 예시

```go
func (s *WebhookService) handleMessageEvent(event Event) {
    if event.Message == nil || event.Message.Type != "text" {
        return
    }

    userMessage := event.Message.Text
    log.Printf("User said: %s", userMessage)

    // 복수 메시지 답장
    replies := []string{
        "메시지를 받았습니다!",
        fmt.Sprintf("내용: %s", userMessage),
        "무엇을 도와드릴까요?",
    }

    if err := s.lineClient.ReplyMessages(event.ReplyToken, replies); err != nil {
        log.Printf("Failed to reply: %v", err)
    }
}
```

---

## 주요 기능

✅ **LINE API 스펙 준수**
- URL: `{apiURL}/v2/bot/message/reply`
- Authorization: `Bearer {accessToken}`
- Content-Type: `application/json`

✅ **복수 메시지 지원**
- 최대 5개 메시지까지 한 번에 전송
- 5개 초과 시 에러 반환

✅ **상세한 로깅**
- 전송하는 JSON 데이터 로깅
- 에러 발생 시 상세 정보 로깅
- 성공 시 로그 출력

---

## 에러 처리

```go
// 5개 초과 에러
err := lineClient.ReplyMessages(replyToken, []string{"1", "2", "3", "4", "5", "6"})
// Error: LINE API allows maximum 5 messages per reply, got 6

// API 호출 실패
err := lineClient.ReplyMessage(replyToken, "Hello")
// Error: reply failed with status 400: {"message":"Invalid request"}
```

---

## 환경 변수 설정

`.env` 파일에 다음 값들이 필요합니다:

```env
KID=your_kid_here
CHANNEL_ID=your_channel_id_here
LINE_API_PREFIX=https://api.line.me
```

---

## 현재 의존성 흐름

```
main.go
  └─ fx.Provide(jwt.GetAccessToken, config.GetEnv)
       └─ WebhookModule
            └─ NewLineClient(accessToken, env)
                 └─ LineClient.ReplyMessage()
```

완벽하게 LINE API 스펙에 맞춰서 구현되었습니다! 🎉
