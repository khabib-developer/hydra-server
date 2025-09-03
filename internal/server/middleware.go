package server

import "net/http"

func (server *Server) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connID := r.Header.Get("connID")
		if connID == "" {
			http.Error(w, "missing connID header", http.StatusUnauthorized)
			return
		}

		// check if connID exists in Users map
		if _, ok := server.Users[connID]; !ok {
			http.Error(w, "invalid connID", http.StatusUnauthorized)
			return
		}

		// continue to the next handler
		next.ServeHTTP(w, r)
	}
}

