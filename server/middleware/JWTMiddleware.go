package middleware

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"net/http"
	"os"
	"strings"
)

// JWTMiddleware openssl rsa -in private.key -pubout -out public.key
type JWTMiddleware struct {
	Name   string
	Config JWTMiddlewareConfig

	publicKey *rsa.PublicKey
}

type JWTMiddlewareConfig struct {
	PublicKey string `yaml:"public_key"`
	CtxKey    string `yaml:"ctx_key"`
}

func NewJWTMiddleware(name string, cfg JWTMiddlewareConfig) interfaces.IMiddleware {
	return &JWTMiddleware{
		Name:   name,
		Config: cfg,
	}
}

func (t *JWTMiddleware) Init(app interfaces.IEngine) error {
	pubKeyData, err := os.ReadFile(t.Config.PublicKey)
	if err != nil {
		return fmt.Errorf("could not read public key: %w", err)
	}

	pubBlock, _ := pem.Decode(pubKeyData)
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		return fmt.Errorf("invalid public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return fmt.Errorf("could not parse public key: %w", err)
	}

	var ok bool
	t.publicKey, ok = pubKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an RSA public key")
	}

	return nil
}

func (t *JWTMiddleware) Stop() error {
	return nil
}

func (t *JWTMiddleware) String() string {
	return t.Name
}

func (t *JWTMiddleware) Invoke(next interfaces.RouteFunc) interfaces.RouteFunc {
	return func(c context.Context, ctx interfaces.ICtx) {
		authHeader := ctx.GetRequest().GetHeader().Get("Authorization")
		if authHeader == "" {
			ctx.GetResponse().SetStatus(http.StatusUnauthorized)
			ctx.GetResponse().SetBody([]byte("Authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.GetResponse().SetStatus(http.StatusUnauthorized)
			ctx.GetResponse().SetBody([]byte("Invalid Authorization header format"))
			return
		}

		tokenString := parts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return t.publicKey, nil
		})

		if err != nil || !token.Valid {
			ctx.GetResponse().SetStatus(http.StatusUnauthorized)
			ctx.GetResponse().SetBody([]byte("Invalid or expired token"))
			return
		}

		cNew := context.WithValue(c, t.Config.CtxKey, claims)
		next(cNew, ctx)
	}
}
