package websocket

import (
	"fmt"
	"slices"
)

type MessageType string

const (
	MessageTypeError      MessageType = "error"
	MessageTypePassword   MessageType = "password"
	MessageTypeInfo       MessageType = "info"
	MessageTypeMessage    MessageType = "message"
	MessageTypeBroadcast  MessageType = "broadcast"
	MessageTypeClose      MessageType = "close"
	MessageTypeJoin       MessageType = "join"
	MessageTypeCreate     MessageType = "create"
	MessageTypeDestroy    MessageType = "destroy"
	MessageTypeCancel     MessageType = "cancel"
	MessageTypeFile       MessageType = "file"
	MessageTypeFileChunk  MessageType = "file_chunk"
)


var AllMessageTypes = []MessageType{
	MessageTypeError,
	MessageTypeInfo,
	MessageTypeMessage,
	MessageTypePassword,
	MessageTypeBroadcast,
	MessageTypeClose,
	MessageTypeJoin,
	MessageTypeCreate,
	MessageTypeDestroy,
	MessageTypeFile,
	MessageTypeFileChunk,
}

func (c MessageType) IsValid() bool {
	return slices.Contains(AllMessageTypes, c)
}


func (w WebsocketDto) Validate() error {
	if !w.MessageType.IsValid() {
		return fmt.Errorf("invalid command: %s", w.MessageType)
	}
	return nil
}
