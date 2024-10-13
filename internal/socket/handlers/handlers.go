package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
}

type ClientManager struct {
	Clients   map[*Client]bool
	BroadCast chan []byte
	Mutex     sync.Mutex
}

var Manager = ClientManager{
	Clients:   make(map[*Client]bool),
	BroadCast: make(chan []byte),
}

func (m *ClientManager) AddClient(c *Client) {
	m.Mutex.Lock()
	m.Clients[c] = true
	m.Mutex.Unlock()
}

func (m *ClientManager) DeleteClient(c *Client) {
	m.Mutex.Lock()
	delete(m.Clients, c)
	m.Mutex.Unlock()
}

func (m *ClientManager) BroadCastMessage(data []byte) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	fmt.Println("Number of clients before broadcasting:", len(m.Clients))

	for client := range m.Clients {
		fmt.Println("Broadcasting to client:", client)
		err := client.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Error broadcasting to client:", err)
			client.conn.Close()
			m.DeleteClient(client)
		}
	}
}

func SocketHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade to websocket", err)
		}

		client := &Client{conn: conn}
		Manager.AddClient(client)

		defer func() {
			client.conn.Close()
			Manager.DeleteClient(client)
		}()

		fmt.Println("Number of clients:", len(Manager.Clients))
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				break
			}
			log.Printf("Received message: %s", message)

			// Echo the message back to the client
			if err := conn.WriteMessage(messageType, message); err != nil {
				log.Printf("Error writing message: %v", err)
				break
			}

		}

	}
}

func BroadcastHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error while reading body: %v", err)
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		Manager.BroadCastMessage(data)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message broadcasted successfully"))
	}
}
