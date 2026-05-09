package jwtval

import (
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/golang-jwt/jwt/v5"
)

const defaultTTL = 72 * time.Hour
const defaultSubject = "xrayman subject"

const bearerTokenType = "Bearer"

type genOpts struct {
	ttl     time.Duration
	subject string
}

type genOpt = func(o *genOpts)

func WithTTL(ttl time.Duration) genOpt {
	return func(o *genOpts) {
		o.ttl = ttl
	}
}

func WithSubject(subject string) genOpt {
	return func(o *genOpts) {
		o.subject = subject
	}
}

func GenerateToken(sec []byte, iss string, options ...genOpt) (string, error) {
	cfg := genOpts{
		ttl:     defaultTTL,
		subject: defaultSubject,
	}
	for _, o := range options {
		o(&cfg)
	}

	now := time.Now()
	exp := now.Add(cfg.ttl)

	claims := jwt.RegisteredClaims{
		Issuer:    iss,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(exp),
		Subject:   cfg.subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(sec)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateToken(tok string, sec []byte, iss string) error {
	secret := []byte(sec)
	// Parse the token
	token, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	// check parsing
	if err != nil {
		return err
	}
	// check method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return xerr.New("invalid singing method")
	}
	// check claims
	if claimIss, err := token.Claims.GetIssuer(); err != nil || claimIss != iss {
		return xerr.Newf("invalid issuer: %s", claimIss)
	}
	if exp, err := token.Claims.GetExpirationTime(); err != nil || exp.Before(time.Now()) {
		return xerr.Newf("invalid exp time: %v", exp.Time)
	}

	return nil
}
