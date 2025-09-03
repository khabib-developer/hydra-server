package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/khabib-developer/hydra-server/pkg/auth"
)

func (server *Server) Auth(w http.ResponseWriter, r *http.Request) {
	req := auth.AuthDto{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if candidate := server.usernameExists(req.Username); candidate {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	connID := uuid.New().String()

	server.Users[connID] = &User{
		ID:       connID,
		Username: req.Username,
		Password: req.Password,
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(connID))
}
func (server *Server) usernameExists(username string) bool {
	for _, user := range server.Users {
		if strings.EqualFold(user.Username, username) {
			return true
		}
	}
	return false

}

