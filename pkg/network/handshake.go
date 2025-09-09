package network

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/khabib-developer/hydra-server/pkg/security"
)

func SendPublicAddress(safeConn *SafeConn, public_key *rsa.PublicKey, hardcoded_private_key *rsa.PrivateKey) error {

	EncodedPublicKey, err := security.EncodePublicKeyToPEM(public_key)

	if err != nil {
		return err
	}

	JsonPublicKey, err := json.Marshal(EncodedPublicKey)

	if err != nil {
		return err
	}

	data := Dto{
		MessageType: PUBLIC_KEY_EXCHANGE,
		Data:        JsonPublicKey,
	}

	err = SendMessageWithSign(safeConn, data, hardcoded_private_key)

	return err
}

func ReceivePublicAddress(safeConn *SafeConn, hardcoded_public_key *rsa.PublicKey) error {
	payload, err := safeConn.ReadMessage()
	if err != nil {
		return err
	}

	data := SignedDto{}
	err = json.Unmarshal(payload, &data)

	if err != nil {
		return err
	}

	jsonDto, err := json.Marshal(data.Dto)

	if err != nil {
		return err
	}

	err = security.Verify(hardcoded_public_key, jsonDto, data.Signature)

	return err
}
