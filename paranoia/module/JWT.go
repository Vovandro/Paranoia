package module

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	interfaces2 "gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/go_utils/decode"
	"os"
	"time"
)

// JWT openssl genrsa -out private.key 2048
type JWT struct {
	config     JWTConfig
	name       string
	privateKey *rsa.PrivateKey
}

type JWTConfig struct {
	PrivateKey string        `yaml:"public_key"`
	Expire     time.Duration `yaml:"expire"`
}

func NewJWT(name string) interfaces2.IModules {
	return &JWT{
		name: name,
	}
}

func (t *JWT) Init(_ interfaces2.IEngine, cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.PrivateKey == "" {
		return errors.New("missing private key")
	}

	privKeyData, err := os.ReadFile(t.config.PrivateKey)
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

func (t *JWT) Name() string {
	return t.name
}

func (t *JWT) Type() string {
	return "module"
}

func (t *JWT) GenerateToken(data jwt.MapClaims) (string, error) {
	if _, ok := data["exp"]; !ok {
		data["exp"] = time.Now().Add(t.config.Expire).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, data)

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return tokenString, nil
}
