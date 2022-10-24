package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2PasswordHandler struct {
	EncryptionParameters *EncryptionParameters
}

func NewPasswordHandler(EncryptionParameters *EncryptionParameters) PasswordHandler {
	instance := new(Argon2PasswordHandler)
	instance.EncryptionParameters = EncryptionParameters
	return instance
}

func (cp *Argon2PasswordHandler) CreateHash(Password string) (string, error) {
	if len(Password) == 0 {
		return "", errors.New(PASSWORD_IS_EMPTY)
	}
	// Generate a cryptographically secure random salt.
	salt, err := cp.generateRandomBytes(cp.EncryptionParameters.SaltLength)
	if err != nil {
		return "", err
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey(
		[]byte(Password),
		salt, cp.EncryptionParameters.Iterations,
		cp.EncryptionParameters.Memory,
		cp.EncryptionParameters.Parallelism,
		cp.EncryptionParameters.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		cp.EncryptionParameters.Memory,
		cp.EncryptionParameters.Iterations,
		cp.EncryptionParameters.Parallelism,
		b64Salt,
		b64Hash)

	return encodedHash, nil
}

func (cp *Argon2PasswordHandler) generateRandomBytes(n uint32) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (cp *Argon2PasswordHandler) CompareHashToPassword(Password string, PasswordHashToMatch string) (bool, error) {
	if len(Password) == 0 {
		return false, errors.New(PASSWORD_IS_EMPTY)
	}

	if len(PasswordHashToMatch) == 0 {
		return false, errors.New(PASSWORD_HASH_IS_EMPTY)
	}
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	details, err := cp.decodeHash(PasswordHashToMatch)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(Password), details.Salt, details.Iterations, details.Memory, details.Parallelism, details.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(details.Hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (cp *Argon2PasswordHandler) decodeHash(encodedHash string) (DecodedDetails, error) {
	details := DecodedDetails{}
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return details, errors.New(INVALID_HASH)
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return details, err
	}
	if version != argon2.Version {
		return details, errors.New(INCOMPATIBLE_VERSION)
	}

	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &details.Memory, &details.Iterations, &details.Parallelism)
	if err != nil {
		return details, err
	}

	details.Salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return details, err
	}

	details.SaltLength = uint32(len(details.Salt))

	details.Hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return details, err
	}
	details.KeyLength = uint32(len(details.Hash))

	return details, nil
}
