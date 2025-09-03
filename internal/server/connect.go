package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	socket "github.com/khabib-developer/hydra-server/pkg/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

func (server *Server) Connect(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	connID := r.Header.Get("connID")

	safeConn := &socket.SafeConn{
		Conn:  ws,
		Mutex: sync.Mutex{},
	}

	if connID == "" {
		server.closeConnection(safeConn, "connID is not exist")
	}

	err = server.add(connID, safeConn)

	if err != nil {
		server.closeConnection(safeConn, err.Error())
	}

	go server.listen(connID)
}

func (server *Server) closeConnection(safeConn *socket.SafeConn, msg string) {
	server.sendMessage(safeConn, socket.MessageTypeClose, []byte(msg))
	safeConn.Conn.Close()
}

func (server *Server) add(connID string, safeConn *socket.SafeConn) error {
	if server.Users == nil {
		server.Users = make(map[string]*User)
	}

	u, ok := server.Users[connID]
	if !ok {
		return fmt.Errorf("user with connID %q not found", connID)
	}

	// attach connection to the copy and save it back
	u.SafeConn = safeConn
	server.Users[connID] = u

	return nil
}
