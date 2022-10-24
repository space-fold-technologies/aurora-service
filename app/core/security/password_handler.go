package security

const (
	PASSWORD_IS_EMPTY      = "NO PASSWORD HAS BEEN SUPPLIED"
	PASSWORD_HASH_IS_EMPTY = "NO PASSWORD HASH SUPPLIED"
	INVALID_HASH           = "the encoded hash is not in the correct format"
	INCOMPATIBLE_VERSION   = "incompatible version of argon2"
)

type PasswordHandler interface {
	CreateHash(Password string) (string, error)
	CompareHashToPassword(Password string, PasswordHashToMatch string) (bool, error)
}
