package models

import "time"

type Auth struct {
	PasswordHash []byte
}

type AuthParams struct {
	Password string
}

type AuthResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   time.Duration
}
