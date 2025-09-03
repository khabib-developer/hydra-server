package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	messagePkg "github.com/khabib-developer/hydra-server/pkg/message"
	socket "github.com/khabib-developer/hydra-server/pkg/websocket"
)

func (server *Server) sendMessage(safeConn *socket.SafeConn, messageType socket.MessageType, payload json.RawMessage) {
	safeConn.Mutex.Lock()
	defer safeConn.Mutex.Unlock()
	data := socket.WebsocketDto{
		MessageType: messageType,
		Payload:     payload,
	}
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal error:", err)
		safeConn.Conn.Close()
		return
	}

	if err = safeConn.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		fmt.Println("write error:", err)
		safeConn.Conn.Close()
	}
}

func (server *Server) sendDirectMessage(payload json.RawMessage, sender *User) {
	var messagePayload messagePkg.SendMessageDto

	if err := json.Unmarshal(payload, &messagePayload); err != nil {
		server.sendRawMessage(sender.SafeConn, socket.MessageTypeError, "invalid message payload")
		return
	}

	var receiver *User

	for _, user := range server.Users {
		if messagePayload.Receiver == user.Username {
			receiver = user
		}
	}

	if receiver == nil {
		server.sendRawMessage(sender.SafeConn, socket.MessageTypeError, "Username not found")
		return
	}

	if len(receiver.Password) != 0 {
		permanentData := PermanentData{
			Expect:   receiver.Password,
			Data:     messagePayload.Message,
			Receiver: receiver,
		}
		sender.PermanentData = &permanentData
		server.sendRawMessage(sender.SafeConn, socket.MessageTypePassword, "Password of user: ")
		return
	}

	server.sendActualMessage(sender, receiver, messagePayload.Message)
}

func (server *Server) handlePassword(payload json.RawMessage, sender *User) {
	var password string
	if err := json.Unmarshal(payload, &password); err != nil {
		println(err)
		server.sendRawMessage(sender.SafeConn, socket.MessageTypeError, "Unsupperted type of password")
		return
	}

	if sender.PermanentData == nil {
		server.sendRawMessage(sender.SafeConn, socket.MessageTypeError, "User did not expect password")
		return
	}

	if sender.PermanentData.Expect != password {
		server.sendRawMessage(sender.SafeConn, socket.MessageTypeError, "Wrong password")
		return
	}

	server.sendActualMessage(sender, sender.PermanentData.Receiver, sender.PermanentData.Data)

	sender.PermanentData = nil
}

func (server *Server) sendRawMessage(safeConn *socket.SafeConn, messageType socket.MessageType, message string) {
	messageJson, err := json.Marshal(strings.TrimSpace(message))
	if err != nil {
		return
	}
	server.sendMessage(safeConn, messageType, messageJson)

}

func (server *Server) sendActualMessage(sender *User, receiver *User, message string) {
	responsePayloadBytes, error := json.Marshal(messagePkg.ReceiveMessageDto{
		Sender:  sender.Username,
		Message: message,
	})

	if error != nil {
		server.sendRawMessage(sender.SafeConn, socket.MessageTypeError, "invalid message payload")
		return
	}

	server.sendMessage(receiver.SafeConn, socket.MessageTypeMessage, responsePayloadBytes)

	server.sendRawMessage(sender.SafeConn, socket.MessageTypeInfo, "Your message successfully sent")
}
