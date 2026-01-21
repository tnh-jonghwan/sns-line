package sse

import (
	"encoding/json"
	"sync"
)

// Message SSE로 전송할 메시지
type Message struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

// Broadcaster SSE 클라이언트들을 관리
type Broadcaster struct {
	clients map[chan Message]bool
	mu      sync.RWMutex
}

// NewBroadcaster Broadcaster 생성자
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan Message]bool),
	}
}

func (b *Broadcaster) Register(client chan Message) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client] = true
}

func (b *Broadcaster) Unregister(client chan Message) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[client]; ok {
		delete(b.clients, client)
		close(client)
	}
}

func (b *Broadcaster) Broadcast(text, userID string) {
	message := Message{
		Text:   text,
		UserID: userID,
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	for client := range b.clients {
		select {
		case client <- message:
		default:
			// 클라이언트가 응답하지 않으면 스킵
		}
	}
}

func (m Message) ToJSON() string {
	data, _ := json.Marshal(m)
	return string(data)
}
