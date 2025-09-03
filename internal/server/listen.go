package server

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	socket "github.com/khabib-developer/hydra-server/pkg/websocket"
)

func (server *Server) listen(connID string) {
	sender := server.Users[connID]
	defer func() {
		sender.SafeConn.Conn.Close()
		delete(server.Users, connID)

		server.destroyUserChannels(sender)

		fmt.Println("connection closed for", connID)
	}()

	for {
		messageType, msg, err := sender.SafeConn.Conn.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}

		if messageType == websocket.TextMessage {

			var payload socket.WebsocketDto

			if err := json.Unmarshal(msg, &payload); err != nil {
				server.sendMessage(sender.SafeConn, socket.MessageTypeError, []byte(`"wrong type of command"`))
				return
			}

			if handler, ok := server.handlers[payload.MessageType]; ok {
				handler(payload.Payload, sender)
			} else {
				server.sendRawMessage(sender.SafeConn, socket.MessageTypeError, "unknown command")
			}

		}
	}
}
