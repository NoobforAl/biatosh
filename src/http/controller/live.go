package controller

import (
	"biatosh/contract"
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type client struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	c *websocket.Conn
}

type wsStore struct {
	mux     sync.RWMutex
	clients map[string]*client
	log     contract.Logger
}

func newWsStore(log contract.Logger) *wsStore {
	return &wsStore{
		mux:     sync.RWMutex{},
		clients: make(map[string]*client, 1024),
		log:     log,
	}
}

func (w *wsStore) getAllClients() []*client {
	w.mux.RLock()
	defer w.mux.RUnlock()

	clients := make([]*client, 0, len(w.clients))
	for _, c := range w.clients {
		if c != nil {
			clients = append(clients, c)
		}
	}

	return clients
}

func (w *wsStore) addClient(id, name string, c *websocket.Conn) {
	w.mux.Lock()
	defer w.mux.Unlock()

	w.clients[id] = &client{
		Id:   id,
		Name: name,
		c:    c,
	}
}

func (w *wsStore) removeClient(id string) {
	w.mux.Lock()
	defer w.mux.Unlock()

	delete(w.clients, id)
}

func (w *wsStore) sendUsersForAllClient() {
	w.mux.RLock()
	defer w.mux.RUnlock()

	for _, client := range w.clients {
		if err := client.c.WriteJSON(w.getAllClients()); err != nil {
			w.log.Error("Failed to write JSON (to send all users):", err)
			return
		}
	}
}

func (w *wsStore) notifyClient(id string, name string) {
	w.mux.RLock()
	defer w.mux.RUnlock()

	if client, ok := w.clients[id]; ok {
		message := fmt.Sprintf(`{"message": "%s is looking at you"}`, name)
		err := client.c.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			w.log.Error("Failed to notify client:", err)
		}
	}
}
