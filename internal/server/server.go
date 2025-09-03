package server

import (
	"encoding/json"

	"github.com/khabib-developer/hydra-server/pkg/websocket"
)

func NewServer() *Server {
	s := &Server{
		Users: make(map[string]*User),
	}
	s.handlers = map[websocket.MessageType]func(json.RawMessage, *User){
		websocket.MessageTypeMessage:   s.sendDirectMessage,
		websocket.MessageTypePassword:  s.handlePassword,
		websocket.MessageTypeJoin:      s.joinChannel,
		websocket.MessageTypeCreate:    s.createChannel,
		websocket.MessageTypeBroadcast: s.broadcastMessage,
		websocket.MessageTypeDestroy:   s.destroyChannel,
		websocket.MessageTypeFile:      s.TransferFileMetadata,
		websocket.MessageTypeFileChunk: s.TransferFileChunk,
		websocket.MessageTypeCancel:    s.CancelFileTransfer,
	}
	return s
}

type Server struct {
	Users    map[string]*User
	Channels []*Channel
	handlers map[websocket.MessageType]func(json.RawMessage, *User)
}
