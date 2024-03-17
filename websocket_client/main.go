package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial failed:", err)
	}
	defer conn.Close()

	// 发送消息到服务器
	message := []byte("Hello, Server!")
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Write error:", err)
		return
	}

	// 接收服务器的响应
	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	fmt.Printf("Received message: %s\n", p)
}
