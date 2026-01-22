package eventHub

import (
	"encoding/json"
	"log"
	"sync"
)

type Message struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type EventHub struct {
	clients map[chan string]bool // Changed from chan Message to chan string
	mu      sync.RWMutex
}

type BroadcastMessage struct {
	Text     string `json:"text"`
	UserId   string `json:"userId"`
	Platform string `json:"platform"` // "line" or "instagram"
}

func NewEventHub() *EventHub {
	return &EventHub{
		clients: make(map[chan string]bool), // Changed from chan Message to chan string
	}
}

// ---

func (e *EventHub) AddClient(client chan string) { // Renamed from Register, changed signature
	e.mu.Lock()
	defer e.mu.Unlock()
	e.clients[client] = true
	log.Println("eventHub 클라이언트 연결됨") // Added log
}

func (e *EventHub) RemoveClient(client chan string) { // Renamed from Unregister, changed signature
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.clients[client]; ok {
		delete(e.clients, client)
		close(client)
	}
	log.Println("eventHub 클라이언트 연결 해제됨") // Added log
}

func (e *EventHub) Broadcast(text string, userId string, platform string) { // Updated signature
	e.mu.RLock()
	defer e.mu.RUnlock()

	msg := BroadcastMessage{ // Changed to BroadcastMessage
		Text:     text,
		UserId:   userId,
		Platform: platform,
	}

	jsonData, err := json.Marshal(msg) // Marshal BroadcastMessage to JSON
	if err != nil {
		log.Printf("JSON 마샬링 에러: %v", err) // Added error handling
		return
	}

	for client := range e.clients {
		select {
		case client <- string(jsonData): // Send JSON string
		default:
			// 클라이언트가 응답하지 않으면 스킵
		}
	}
}

func (m Message) ToJSON() string {
	data, _ := json.Marshal(m)
	return string(data)
}
