package module

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"os"
	"time"
)

// JWT openssl genrsa -out private.key 2048
type JWT struct {
	Cfg        JWTConfig
	Name       string
	privateKey *rsa.PrivateKey
}

type JWTConfig struct {
	PrivateKey string        `yaml:"public_key"`
	Expire     time.Duration `yaml:"expire"`
}

func NewJWT(name string, cfg JWTConfig) interfaces.IModules {
	return &JWT{
		Name: name,
		Cfg:  cfg,
	}
}

func (t *JWT) Init(_ interfaces.IEngine) error {
	privKeyData, err := os.ReadFile(t.Cfg.PrivateKey)
	if err != nil {
		return fmt.Errorf("could not read private key: %w", err)
	}

	privBlock, _ := pem.Decode(privKeyData)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("invalid private key")
	}

	t.privateKey, err = x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return fmt.Errorf("could not parse private key: %w", err)
	}

	return nil
}

func (t *JWT) Stop() error {
	return nil
}

func (t *JWT) String() string {
	return t.Name
}

func (t *JWT) GenerateToken(data jwt.MapClaims) (string, error) {
	if _, ok := data["exp"]; !ok {
		data["exp"] = time.Now().Add(t.Cfg.Expire).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, data)

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return tokenString, nil
}
