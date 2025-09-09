package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func GenerateKeyPairs(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}

// =========================
// Encode helpers
// =========================

// EncodePrivateKeyToPEM converts *rsa.PrivateKey to a PEM string
func EncodePrivateKeyToPEM(priv *rsa.PrivateKey) string {
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	return string(privPEM)
}

// EncodePublicKeyToPEM converts *rsa.PublicKey to a PEM string
func EncodePublicKeyToPEM(pub *rsa.PublicKey) (string, error) {
	pubBytes := x509.MarshalPKCS1PublicKey(pub)
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	})
	return string(pubPEM), nil
}

// =========================
// Decode helpers
// =========================

// DecodePrivateKeyFromPEM parses a PEM string into *rsa.PrivateKey
func DecodePrivateKeyFromPEM(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key PEM")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// DecodePublicKeyFromPEM parses a PEM string into *rsa.PublicKey
func DecodePublicKeyFromPEM(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("invalid public key PEM")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}
