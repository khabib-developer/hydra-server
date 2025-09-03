package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SafeConn struct {
	Conn *websocket.Conn
	Mutex sync.Mutex
}