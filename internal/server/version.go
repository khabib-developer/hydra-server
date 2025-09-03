package server

import (
	"net/http"

	"github.com/khabib-developer/hydra-server/pkg/version"
)

func (server *Server) Version(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(version.Version))
}
