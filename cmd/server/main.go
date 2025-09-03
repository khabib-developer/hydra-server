package main

import (
	"fmt"
	"net/http"

	"github.com/khabib-developer/hydra-server/internal/server"
)

func main() {
	server := server.NewServer()

	http.HandleFunc("/version", server.Version)
	http.HandleFunc("/auth", server.Auth)
	http.HandleFunc("/connect", server.Connect)
	http.HandleFunc("/getActiveUsers", server.CheckAuth(server.GetActiveUsers))
	http.HandleFunc("/getChannels", server.CheckAuth(server.GetChannels))
	http.HandleFunc("/getChannelMembers", server.CheckAuth(server.GetChannelMembers))

	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
