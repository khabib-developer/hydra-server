package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Sign generates a digital signature for given data using a private key
func Sign(priv *rsa.PrivateKey, data []byte) ([]byte, error) {
	// Hash the data
	hash := sha256.Sum256(data)

	// Sign the hash
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash[:])
}

// Verify checks a digital signature using the public key
func Verify(pub *rsa.PublicKey, data []byte, signature []byte) error {
	// Hash the data
	hash := sha256.Sum256(data)

	// Verify the signature
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], signature)
}
