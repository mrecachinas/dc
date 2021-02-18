package api

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

// UpdaterWebsocket handles broadcasting updates to the database;
// Note: this just means querying the collection every
// `api.Cfg.PollingInterval` seconds and returning the *entire*
// contents of the collection.
func (a *Api) UpdaterWebsocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer ws.Close()

	// Create a new connection entry in our connection map
	a.Websocket.Lock()
	a.Websocket.Connections[ws] = struct{}{}
	defer a.Websocket.CloseWebsocketConnection(ws)
	a.Websocket.Unlock()

	msg := fmt.Sprintf(
		"Client %s joined. %d total connections.",
		c.Request().Host,
		len(a.Websocket.Connections),
	)
	c.Logger().Info(msg)

	delay := time.Duration(a.Cfg.PollingInterval) * time.Second

	// TODO: Make error handling/logging more useful? Break if error occurs?
	// TODO: Figure out a way to break gracefully on ctrl-c?
	for {
		select {
		case <-time.After(delay):
			// Get ALL records every N seconds
			statusList, err := GetAllStatus(a.DB.Database(a.Cfg.MongoDatabaseName))
			if err != nil {
				c.Logger().Error("Error getting all status from database")
			}

			// Broadcast result to all connections
			err = a.Websocket.SendMessageToPool(statusList)
			if err != nil {
				c.Logger().Error("Error sending message to pool")
			}

			// Read inbound websocket messages;
			// if we get an error back, that means the client
			// closed the connection
			_, _, err = ws.ReadMessage()
			if err != nil {
				msg := fmt.Sprintf(
					"Client %s quitting. %d connections remain.",
					c.Request().Host,
					len(a.Websocket.Connections),
				)
				c.Logger().Info(msg)
				break
			}
		}
	}
}

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
