package middleware

import (
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/XRay-Addons/xrayman/node/internal/http/errproc"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"go.uber.org/zap"
)

const authIssuer = "xray-node"

func Auth(jwtkey []byte, log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		jwtauth := func(w http.ResponseWriter, r *http.Request) {
			// check auth header
			if err := checkAuth(r, jwtkey); err != nil {
				errproc.ResponseErr(w, http.StatusUnauthorized, "invalid jwt")
				errproc.LogRequestErr(r.Context(), log, err)
				return
			}

			// sign auth header
			err := signAuth(w, jwtkey)
			if err != nil {
				errproc.ResponseErr(w, http.StatusInternalServerError, "")
				errproc.LogResponseErr(r.Context(), log, err)
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
		return err
	}
	// don't check token, but maybe later...
	_ = verifiedToken

	return nil
}

func signAuth(w http.ResponseWriter, key []byte) error {
	// set response authorization header
	tok, err := jwt.NewBuilder().Issuer(authIssuer).IssuedAt(time.Now()).Build()
	if err != nil {
		return err
	}
	// sign it
	sign, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256(), key))
	if err != nil {
		return err
	}
	// add this header
	w.Header().Set(constants.AuthHeader, string(sign))
	return nil
}
