package middleware

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"net/http"
	"os"
	"strings"
)

// JWTMiddleware openssl rsa -in private.key -pubout -out public.key
type JWTMiddleware struct {
	name   string
	config JWTMiddlewareConfig

	publicKey *rsa.PublicKey
}

type JWTMiddlewareConfig struct {
	PublicKey string `yaml:"public_key"`
	CtxKey    string `yaml:"ctx_key"`
}

func NewJWTMiddleware(name string) interfaces2.IMiddleware {
	return &JWTMiddleware{
		name: name,
	}
}

func (t *JWTMiddleware) Init(app interfaces2.IEngine, cfg map[string]interface{}) error {

	pubKeyData, err := os.ReadFile(t.config.PublicKey)
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

func (t *JWTMiddleware) Name() string {
	return t.name
}

func (t *JWTMiddleware) Type() string {
	return "middleware"
}

func (t *JWTMiddleware) Invoke(next interfaces2.RouteFunc) interfaces2.RouteFunc {
	return func(c context.Context, ctx interfaces2.ICtx) {
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

		cNew := context.WithValue(c, t.config.CtxKey, claims)
		next(cNew, ctx)
	}
}
