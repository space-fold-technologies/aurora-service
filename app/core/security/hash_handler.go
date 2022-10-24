package security

type HashHandler interface {
	Encrypt(content []byte) ([]byte, error)
	Decrypt(content []byte) ([]byte, error)
}
