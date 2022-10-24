package security

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	jwt "gopkg.in/square/go-jose.v2/jwt"
)

const (
	RS_256                 = "RS256"
	RS_384                 = "RS384"
	RS_512                 = "RS512"
	NO_SIGNING_KEY_PRESENT = "no Signing Key Present"
	MALFORMED_TOKEN_ERROR  = "session Token is malformed"
	TOKEN_EXPIRED          = "session has expired"
	INVALID_TOKEN_ERROR    = "token is not valid or possibly not signed from this service"
)

var (
	ErrKeyMustBePEMEncoded = errors.New("invalid key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotRSAPrivateKey    = errors.New("key is not a valid RSA private key")
	ErrNotRSAPublicKey     = errors.New("key is not a valid RSA public key")
)

type RSATokenVerifier struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	log        *logrus.Logger
}

func New(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) TokenVerifier {
	instance := new(RSATokenVerifier)
	instance.log = logging.GetInstance()
	instance.publicKey = publicKey
	instance.privateKey = privateKey
	return instance
}

func (ckv *RSATokenVerifier) VerifyToken(Token string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseEncrypted(Token)
	if err != nil {
		ckv.log.Printf(err.Error())
		return claims, errors.New(MALFORMED_TOKEN_ERROR)
	}
	if err := tok.Claims(ckv.privateKey, &claims); err != nil {
		ckv.log.Printf(err.Error())
		return claims, errors.New(MALFORMED_TOKEN_ERROR)
	}
	if claims.Expiry.Time().Before(time.Now()) {
		return claims, errors.New(TOKEN_EXPIRED)
	}
	return claims, nil
}
