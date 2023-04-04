package jwt

import (
	"context"
	"fmt"
	"net/http"
	jwtlib "osoc-dialog/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"osoc-dialog/pkg/log"
)

type claimsContextKey struct{}

const XUserIDKey = "x-user-id"

func New(opts ...Option) gin.HandlerFunc {
	o := &options{
		logger: log.NewDiscardLogger(),
	}
	for _, opt := range opts {
		opt(o)
	}

	var keyFunc jwt.Keyfunc

	if len(o.hmacSecret) == 0 {
		o.logger.Error().Msg("JWT auth: no HMAC secret provided")
	} else {
		keyFunc = hmacKeyFunc(o.hmacSecret)
	}

	return func(c *gin.Context) {
		if len(o.hmacSecret) == 0 {
			o.logger.Error().Msg("JWT auth: HMAC secret is not defined")
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			c.Abort()
			return
		}

		scheme, tokenString := parseAuthHeader(c.GetHeader("Authorization"))
		if scheme != "Bearer" {
			c.String(http.StatusUnauthorized, "unexpected authorization header format")
			c.Abort()
			return
		}
		if tokenString == "" {
			c.String(http.StatusUnauthorized, "JWT is missing")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwtlib.CustomClaims{}, keyFunc)
		if err != nil {
			if _, ok := err.(*jwt.ValidationError); ok {
				c.String(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}

			o.logger.Err(err).Msg("JWT auth: could not parse token string")
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			c.Abort()
			return
		}

		if !token.Valid {
			c.String(http.StatusUnauthorized, "JWT is invalid")
			c.Abort()
			return
		}

		userID, err := jwtlib.ExtractIDFromToken(tokenString, string(o.hmacSecret))
		if err != nil {
			o.logger.Err(err).Msg("JWT auth: could not extract user id from token")
			c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			c.Abort()
			return
		}
		c.Set(XUserIDKey, userID)

		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), claimsContextKey{}, token.Claims),
		)

		c.Next()
	}
}

func FromContext(ctx context.Context) (*jwtlib.CustomClaims, bool) {
	v, ok := ctx.Value(claimsContextKey{}).(*jwtlib.CustomClaims)
	return v, ok
}

func parseAuthHeader(s string) (scheme, token string) {
	chunks := strings.Split(strings.Trim(s, " "), " ")
	if len(chunks) == 2 {
		scheme = chunks[0]
		token = chunks[1]
	}

	return
}

func hmacKeyFunc(hmacSecret []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSecret, nil
	}
}
