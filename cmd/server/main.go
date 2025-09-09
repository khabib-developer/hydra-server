package main

import (
	"fmt"

	"github.com/khabib-developer/hydra-server/internal/server"
	"github.com/khabib-developer/hydra-server/pkg/network"
)

func main() {
	server, err := network.NewHydraServer("localhost:9087", server.RootPrivKey, server.RootPublicKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = server.ListenAndServe()
	fmt.Println("server successfully ran")
	if err != nil {
		fmt.Println(err)
	}
}
