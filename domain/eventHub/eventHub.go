package eventHub

import (
	"encoding/json"
	"sync"
)

type Message struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type EventHub struct {
	clients map[chan Message]bool
	mu      sync.RWMutex
}

func NewEventHub() *EventHub {
	return &EventHub{
		clients: make(map[chan Message]bool),
	}
}

// ---

func (h *EventHub) Register(client chan Message) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
}

func (h *EventHub) Unregister(client chan Message) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client)
	}
}

func (h *EventHub) Broadcast(text, userID string) {
	message := Message{
		Text:   text,
		UserID: userID,
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
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
