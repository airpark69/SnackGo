package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// 클라이언트와의 연결을 추적
var clients = make(map[*websocket.Conn]bool)

// 메시지 브로드캐스트 채널
var broadcast = make(chan Message)

// 업그레이드 구성을 위해 사용
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message struct to hold the message data
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// 최근 30개의 메시지를 저장하는 슬라이스
var messageHistory []Message

// 초기 ws 연결 시 클라이언트와 소통하는 부분
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to upgrade to websocket: %v", err)
	}
	defer ws.Close()

	// 연결된 클라이언트에 메시지 히스토리 전송
	for _, msg := range messageHistory {
		if err := ws.WriteJSON(msg); err != nil {
			log.Printf("error: %v", err)
		}
	}
	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

// 연결 이후 클라이언트와 소통하는 부분
func HandleMessages() {
	for {
		msg := <-broadcast

		// 메시지 히스토리에 추가하고 30개로 제한
		messageHistory = append(messageHistory, msg)

		if len(messageHistory) > 30 {
			messageHistory = messageHistory[1:]
		}

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
