package network

import (
	"crypto/rsa"
	"encoding/json"
)

func SendMessage(safeConn *SafeConn, data Dto) error {
	payload, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = safeConn.WriteMessage(payload)
	return err
}

func SendMessageWithSign(safeConn *SafeConn, data Dto, priv *rsa.PrivateKey) error {
	signedPayload, err := data.SignAndMarshalDto(priv)
	if err != nil {
		return err
	}

	err = safeConn.WriteMessage(signedPayload)
	return err
}
