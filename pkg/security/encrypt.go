package security

import (
	"crypto/rand"
	"crypto/rsa"
)

// Encrypt encrypts data with a public key
func Encrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// Decrypt decrypts data with a private key
func Decrypt(priv *rsa.PrivateKey, cipher []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipher)
}
