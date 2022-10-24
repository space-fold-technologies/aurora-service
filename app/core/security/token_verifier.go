package security

type TokenVerifier interface {
	VerifyToken(Token string) (*Claims, error)
}
