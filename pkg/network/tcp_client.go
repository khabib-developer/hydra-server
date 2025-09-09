package network

import (
	"crypto/rsa"
	"net"

	"github.com/khabib-developer/hydra-server/pkg/security"
)

type HydraClient struct {
	addr        string
	safeConn    *SafeConn
	public_key  *rsa.PublicKey
	private_key *rsa.PrivateKey
}

func NewHydraClient(addr string) *HydraClient {
	priv, pub, _ := security.GenerateKeyPairs(2048)
	return &HydraClient{addr: addr, public_key: pub, private_key: priv}
}

func (c *HydraClient) Connect(hardcoded_private_key string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.safeConn = &SafeConn{
		Conn: conn,
	}

	priv_key, err := security.DecodePrivateKeyFromPEM(hardcoded_private_key)

	if err != nil {
		return err
	}

	err = SendPublicAddress(c.safeConn, c.public_key, priv_key)

	return err

}
