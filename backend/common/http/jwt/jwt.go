package jwtval

import (
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/golang-jwt/jwt/v5"
)

const defaultTTL = 72 * time.Hour
const defaultSubject = ""
const defaultIssuer = ""

const bearerTokenType = "Bearer"

type options struct {
	ttl     time.Duration
	subject string
	issuer  string
}

type option = func(o *options)

func WithTTL(ttl time.Duration) option {
	return func(o *options) {
		o.ttl = ttl
	}
}

func WithSubject(subject string) option {
	return func(o *options) {
		o.subject = subject
	}
}

func WithIssuer(issuer string) option {
	return func(o *options) {
		o.issuer = issuer
	}
}

func GenerateToken(sec []byte, opts ...option) (string, error) {
	o := options{
		ttl:     defaultTTL,
		subject: defaultSubject,
		issuer:  defaultIssuer,
	}
	for _, opt := range opts {
		opt(&o)
	}

	now := time.Now()
	exp := now.Add(o.ttl)

	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    o.issuer,
		Subject:   o.subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(sec)
	if err != nil {
		return "", err
	}

	return signed, nil
}

type checks struct {
	subject *string
	issuer  *string
}

type check = func(c *checks)

func WithSubjectCheck(subject *string) check {
	return func(c *checks) {
		c.subject = subject
	}
}

func WithIssuerCheck(issuer *string) check {
	return func(c *checks) {
		c.issuer = issuer
	}
}
func ValidateToken(tok string, sec []byte, chks ...check) error {
	c := checks{}
	for _, chk := range chks {
		chk(&c)
	}

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
	// check expiration time
	if exp, err := token.Claims.GetExpirationTime(); err != nil || exp.Before(time.Now()) {
		return xerr.Newf("invalid exp time: %v", exp.Time)
	}
	// check claims
	if c.issuer != nil {
		if claimIss, err := token.Claims.GetIssuer(); err != nil || claimIss != *c.issuer {
			err := xerr.Newf("invalid issuer: %s, required issuer: %s", claimIss, *c.issuer)
			return err
		}
	}
	if c.subject != nil {
		if claimSubj, err := token.Claims.GetSubject(); err != nil || claimSubj != *c.subject {
			err := xerr.Newf("invalid subject: %s, required subject: %s", claimSubj, *c.subject)
			return err
		}
	}

	return nil
}
