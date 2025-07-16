package middleware

import (
	"net/http"
	"time"

	"github.com/XRay-Addons/xrayman/node/internal/http/constants"
	"github.com/XRay-Addons/xrayman/node/internal/http/errors"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"go.uber.org/zap"
)

const authIssuer = "xray-node"

func Auth(jwtkey []byte, log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		jwtauth := func(w http.ResponseWriter, r *http.Request) {
			// check request authorization header
			verifiedToken, err := jwt.ParseRequest(r,
				jwt.WithKey(jwa.HS256(), jwtkey),
				jwt.WithHeaderKey(constants.AuthHeader))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// don't check token, but maybe later...
			_ = verifiedToken

			// set response authorization header
			tok, err := jwt.NewBuilder().Issuer(authIssuer).IssuedAt(time.Now()).Build()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errors.LogRequestError(log, r, err)
				return
			}
			// sign it
			sign, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256(), jwtkey))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errors.LogRequestError(log, r, err)
				return
			}
			// add this header
			w.Header().Set(constants.AuthHeader, string(sign))

			// process request
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(jwtauth)
	}
}
