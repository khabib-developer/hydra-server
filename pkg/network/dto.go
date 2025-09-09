package network

import (
	"crypto/rsa"
	"encoding/json"

	"github.com/khabib-developer/hydra-server/pkg/security"
)

type Dto struct {
	MessageType MessageType `json:"messageType"`
	Data        any         `json:"data"`
}

type SignedDto struct {
	Dto
	Signature []byte `json:"signature,omitempty"`
}

func (dto *Dto) SignAndMarshalDto(priv *rsa.PrivateKey) ([]byte, error) {
	jsonDto, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	sig, err := security.Sign(priv, jsonDto)
	if err != nil {
		return nil, err
	}

	data := SignedDto{
		Dto:       *dto,
		Signature: sig,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
