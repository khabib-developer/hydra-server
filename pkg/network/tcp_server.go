package network

import (
	"crypto/rsa"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/khabib-developer/hydra-server/pkg/security"
)

type handler func(*TcpClient)

type HydraServer struct {
	addr           string
	hardcoded_pub  *rsa.PublicKey
	hardcoded_priv *rsa.PrivateKey
	clients        map[uuid.UUID]*TcpClient
	handlers       map[MessageType]handler
}

type TcpClient struct {
	safeConn    *SafeConn
	serverID    uuid.UUID
	public_key  *rsa.PublicKey
	private_key *rsa.PrivateKey
}

func NewHydraServer(addr string, hardcoded_priv *rsa.PrivateKey, hardcoded_pub *rsa.PublicKey) *HydraServer {
	return &HydraServer{addr: addr, hardcoded_pub: hardcoded_pub, hardcoded_priv: hardcoded_priv}
}

func (s *HydraServer) ListenAndServe(handler func(conn net.Conn)) error {

	ln, err := net.Listen("tcp", s.addr)

	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Println("Hydra server listening on", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go s.AcceptConnection(conn)
	}
}

func (s *HydraServer) AcceptConnection(conn net.Conn) {
	safeConn := &SafeConn{
		Conn: conn,
	}

	err := ReceivePublicAddress(safeConn, s.hardcoded_pub)

	if err != nil {
		fmt.Println("wrong sign connection closed: ", err)
		safeConn.Close()
	}

	id := uuid.New()

	priv, pub, err := security.GenerateKeyPairs(2048)

	tcpClient := &TcpClient{
		serverID:    id,
		safeConn:    safeConn,
		public_key:  pub,
		private_key: priv,
	}

	s.clients[id] = tcpClient

	fmt.Println(id)
}

func (s *HydraServer) HandleMessages(safeConn *SafeConn, id uuid.UUID) {
	for {
		payload, err := safeConn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading messagess: ", err)
			delete(s.clients, id)
			safeConn.Close()
		}
		fmt.Println(string(payload))
	}
}
