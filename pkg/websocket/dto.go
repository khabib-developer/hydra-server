package websocket

import "encoding/json"


type WebsocketDto struct {
	MessageType MessageType `json:"command"`
	Payload json.RawMessage `json:"payload"` 
}
