package security

type TokenHandler interface {
	VerifyToken(Token string) (*Claims, error)
	CreateToken(Claims *Claims) (string, error)
}
