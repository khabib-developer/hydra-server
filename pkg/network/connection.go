package network

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
)

type SafeConn struct {
	net.Conn
	mutex sync.Mutex
}

func (c *SafeConn) WriteMessage(data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	length := uint32(len(data))
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, length)

	_, err := c.Conn.Write(append(header, data...))
	return err
}

func (c *SafeConn) ReadMessage() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	header := make([]byte, 4)
	if _, err := io.ReadFull(c.Conn, header); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(header)

	payload := make([]byte, length)
	if _, err := io.ReadFull(c.Conn, payload); err != nil {
		return nil, err
	}

	return payload, nil
}
