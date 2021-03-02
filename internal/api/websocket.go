package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// UpdaterWebsocket handles broadcasting updates to the database;
// Note: this just means querying the collection every
// `api.Cfg.PollingInterval` seconds and returning the *entire*
// contents of the collection.
func (a *Api) UpdaterWebsocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
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
		for range time.Tick(delay) {
			go func() {
				// Get ALL records every N seconds
				statusList, err := a.DB.GetAllStatus()
				if err != nil {
					c.Logger().Error("Error getting all status from database")
				}

				// Broadcast result to all connections
				err = a.Websocket.SendMessageToPool(statusList)
				if err != nil {
					c.Logger().Error(err)
					c.Logger().Error("Error sending message to pool")
				}

				// Read inbound websocket messages;
				// if we get an error back, that means the client
				// closed the connection
				msg := ""
				err = websocket.Message.Receive(ws, &msg)
				if err != nil {
					msg := fmt.Sprintf(
						"Client %s quitting. %d connections remain.",
						c.Request().Host,
						len(a.Websocket.Connections),
					)
					c.Logger().Error(err)
					c.Logger().Info(msg)
					return
				}
			}()
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

// SendMessageToPool sends a message to every connection
func (pool *WebsocketConnectionPool) SendMessageToPool(message interface{}) error {
	jsonMsg, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	pool.RLock()
	defer pool.RUnlock()
	for connection := range pool.Connections {
		// TODO: What to do when this errors?
		if err := websocket.Message.Send(connection, string(jsonMsg)); err != nil {
			continue
		}
	}
	return nil
}

// CloseWebsocketConnection closes a single websocket connection and
// deletes it from our map
func (pool *WebsocketConnectionPool) CloseWebsocketConnection(connection *websocket.Conn) {
	pool.Lock()
	connection.Close()
	delete(pool.Connections, connection)
	pool.Unlock()
}
