package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源的WebSocket连接
		},
	}
	clients = make([]*websocket.Conn, 0) // 存储所有客户端的切片
	mutex   sync.Mutex                   // 定义互斥锁
)

// broadcast 发送消息给所有连接的客户端，但不包括发送消息的客户端
func broadcast(msg []byte, exclude *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, conn := range clients {
		// if conn != exclude { // 不包括发送消息的客户端
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			// 处理错误，例如关闭连接或记录错误
			fmt.Println("Error writing to client:", err)
		}
		// }
	}
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// 将新的客户端添加到客户端列表中
	mutex.Lock()
	clients = append(clients, conn)
	mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Error during read: %v", err)
			}
			break
		}
		fmt.Printf("Received message from %s: %s\n", conn.RemoteAddr(), string(message))

		// 广播消息到所有客户端，除了发送消息的客户端
		broadcast(message, conn)
	}

	// 从客户端列表中移除当前客户端
	mutex.Lock()
	for i, c := range clients {
		if c == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	mutex.Unlock()
}

func main() {
	http.HandleFunc("/ws", handleWebsocket)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
