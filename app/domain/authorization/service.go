package authorization

import (
	"errors"
	"time"

	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
	"github.com/space-fold-technologies/aurora-service/app/domain/users"
	"gopkg.in/square/go-jose.v2/jwt"
)

type AuthorizationService struct {
	repository      AuthorizationRepository
	passwordHandler security.PasswordHandler
	duration        time.Duration
	tokenHandler    security.TokenHandler
}

func NewService(
	repository AuthorizationRepository,
	passwordHandler security.PasswordHandler,
	tokenHandler security.TokenHandler,
	duration time.Duration) *AuthorizationService {
	instance := new(AuthorizationService)
	instance.repository = repository
	instance.passwordHandler = passwordHandler
	instance.tokenHandler = tokenHandler
	instance.duration = duration
	return instance
}

func (as *AuthorizationService) RequestSession(credentials *Credentials) (*Session, error) {
	if user, err := as.repository.Fetch(credentials.GetEmail()); err != nil {
		return nil, err
	} else if match, err := as.passwordHandler.CompareHashToPassword(credentials.GetPassword(), user.PasswordHash); err != nil {
		return nil, err
	} else if !match {
		return nil, errors.New("credentials mismatch")
	} else {
		return as.createSession(user)
	}
}

func (as *AuthorizationService) createSession(user *users.User) (*Session, error) {
	if token, err := as.createToken(user); err != nil {
		return nil, err
	} else {
		return &Session{Token: token}, nil
	}
}

func (as *AuthorizationService) createToken(user *users.User) (string, error) {
	if location, err := time.LoadLocation("UTC"); err != nil {
		return "", err
	} else {
		currentTime := time.Now().In(location)
		expiryTime := currentTime.Add(time.Minute * as.duration)
		claims := &security.Claims{
			Issuer:      "aurora-service",
			ID:          user.Identifier.String(),
			Subject:     user.Identifier.String(),
			Expiry:      jwt.NewNumericDate(expiryTime),
			IssuedAt:    jwt.NewNumericDate(currentTime),
			Permissions: as.transformPermissions(user.Permissions),
			Teams:       as.transformTeams(user.Teams),
		}
		return as.tokenHandler.CreateToken(claims)
	}
}

func (as *AuthorizationService) transformTeams(list []*teams.Team) []string {
	entries := make([]string, 0)
	for _, entry := range list {
		entries = append(entries, entry.Name)
	}
	return entries
}

func (as *AuthorizationService) transformPermissions(list []users.Permission) []string {
	entries := make([]string, 0)
	for _, entry := range list {
		entries = append(entries, entry.Name)
	}
	return entries
}
