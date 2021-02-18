package api

import (
	"github.com/gorilla/websocket"
)

func (pool *WebsocketConnectionPool) SendMessageToPool(message interface{}) error {
	pool.RLock()
	defer pool.RUnlock()
	for connection := range pool.Connections {
		if err := connection.WriteJSON(&message); err != nil {
			return err
		}
	}
	return nil
}

func (pool *WebsocketConnectionPool) CloseWebsocketConnection(connection *websocket.Conn) {
	pool.Lock()
	connection.Close()
	delete(pool.Connections, connection)
	pool.Unlock()
}
