package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/khabib-developer/hydra-server/pkg/user"
	"github.com/khabib-developer/hydra-server/pkg/websocket"
)

type User struct {
	ID            string
	Username      string
	Password      string
	PermanentData *PermanentData
	JoinedChannel *Channel
	PermanentFile *PermanentFileData
	SafeConn      *websocket.SafeConn
}

type PermanentData struct {
	Expect   string
	Data     string
	Receiver *User
}

func (server *Server) GetActiveUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]user.UserDTO, 0, len(server.Users))
	for _, userItem := range server.Users {
		users = append(users, user.UserDTO{
			Username: userItem.Username,
			Private:  len(userItem.Password) > 0,
		})
	}

	msg, err := json.Marshal(users)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
