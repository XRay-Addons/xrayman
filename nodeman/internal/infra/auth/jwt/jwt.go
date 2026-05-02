package jwt

import (
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/http/security"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"github.com/XRay-Addons/xrayman/nodeman/internal/service/auth"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

var _ (auth.JWT) = (*JWT)(nil)
var _ (security.JWT) = (*JWT)(nil)

const defaultTTL = 72 * time.Hour
const defaultIssuer = "xrayman-nodeman"
const adminSubject = "admin"
const bearerTokenType = "Bearer"

func New(secret string) (*JWT, error) {
	if secret == "" {
		return nil, errdefs.NewNilArg("secret")
	}
	return &JWT{
		secret: []byte(secret),
		ttl:    defaultTTL,
		issuer: defaultIssuer,
	}, nil
}

func (j *JWT) GenerateToken() (models.AuthResult, error) {
	now := time.Now()
	exp := now.Add(j.ttl)

	claims := jwt.RegisteredClaims{
		Issuer:    j.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(exp),
		Subject:   adminSubject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(j.secret)
	if err != nil {
		return models.AuthResult{}, err
	}

	return models.AuthResult{
		AccessToken: signed,
		TokenType:   bearerTokenType,
		ExpiresIn:   j.ttl,
	}, nil
}

func (j *JWT) ValidateToken(tokenString string) error {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	// check parsing
	if err != nil {
		return err
	}
	// check method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return errdefs.NewAccessDenied()
	}
	// check claims
	if iss, err := token.Claims.GetIssuer(); err != nil || iss != j.issuer {
		return errdefs.NewAccessDenied()
	}
	if exp, err := token.Claims.GetExpirationTime(); err != nil || exp.Before(time.Now()) {
		return errdefs.NewAccessDenied()
	}

	return nil
}
