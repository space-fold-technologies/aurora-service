package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

type CryptHashHandler struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewHashHandler(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) HashHandler {
	instance := new(CryptHashHandler)
	instance.publicKey = publicKey
	instance.privateKey = privateKey
	return instance
}

func (chh *CryptHashHandler) Encrypt(content []byte) ([]byte, error) {
	hash := sha512.New()
	if ciphered, err := rsa.EncryptOAEP(hash, rand.Reader, chh.publicKey, content, nil); err != nil {
		return []byte{}, err
	} else {
		return ciphered, nil
	}
}

func (chh *CryptHashHandler) Decrypt(content []byte) ([]byte, error) {
	hash := sha512.New()
	if plain, err := rsa.DecryptOAEP(hash, rand.Reader, chh.privateKey, content, nil); err != nil {
		return []byte{}, err
	} else {
		return plain, nil
	}
}
