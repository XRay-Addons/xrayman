package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/XRay-Addons/xrayman/node/internal/http/httperr"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"go.uber.org/zap"
)

const authIssuer = "xray-node"

func Auth(jwtkey []byte, log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		jwtauth := func(w http.ResponseWriter, r *http.Request) {
			var err error
			defer func() { httperr.Write(r.Context(), err, w, log) }()

			// check auth header
			if err = checkAuth(r, jwtkey); err != nil {
				err = httperr.New(httperr.ErrAuth, err)
				return
			}

			// sign auth header
			if err = signAuth(w, jwtkey); err != nil {
				return
			}

			// process request
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(jwtauth)
	}
}

func checkAuth(r *http.Request, key []byte) error {
	// check request authorization header
	verifiedToken, err := jwt.ParseRequest(r,
		jwt.WithKey(jwa.HS256(), key),
		jwt.WithHeaderKey(constants.AuthHeader))
	if err != nil {
		return fmt.Errorf("check auth: %w", err)
	}
	// don't check token, but maybe later...
	_ = verifiedToken

	return nil
}

func signAuth(w http.ResponseWriter, key []byte) error {
	// set response authorization header
	tok, err := jwt.NewBuilder().Issuer(authIssuer).IssuedAt(time.Now()).Build()
	if err != nil {
		return fmt.Errorf("sign auth: %w", err)
	}
	// sign it
	sign, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256(), key))
	if err != nil {
		return fmt.Errorf("sign auth: %w", err)
	}
	// add this header
	w.Header().Set(constants.AuthHeader, string(sign))
	return nil
}
