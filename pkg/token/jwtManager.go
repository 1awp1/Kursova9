package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwksPath = "./.well-known"
)

var (
	ErrExpired        = "token is expired"
	ErrClaims         = "error parsing claims"
	ErrEmptyKey       = "empty key"
	ErrTokenIsInvalid = "token is invalid"
)

type TokenManager interface {
	NewJWT(input AuthInfo) (string, error)
	Parse(accessToken string) (result AuthInfo, err error)
	RefreshToken() (string, error)
	ValidateRefreshToken(tokenString string) (bool, error)
}

type Manager struct {
	signingMethod jwt.SigningMethod
	ttl           time.Duration
	privateKey    *rsa.PrivateKey
	publickKey    *rsa.PublicKey
}

func NewManager(ttl time.Duration, privateKeyPath, publicKeyPath string) (*Manager, error) {
	m := &Manager{signingMethod: jwt.SigningMethodRS256, ttl: ttl}
	err := m.loadKeys(privateKeyPath, publicKeyPath)
	return m, err
}

func (m *Manager) loadKeys(privateKeyPath, publicKeyPath string) error {
	privKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("cannot read private key file: %s", err.Error())
	}
	block, _ := pem.Decode(privKeyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return fmt.Errorf("invalid private key format")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("cannot parse private key: %s", err.Error())
	}
	m.privateKey = privateKey.(*rsa.PrivateKey)

	pubKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("cannot read public key file: %s", err.Error())
	}
	block, _ = pem.Decode(pubKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return fmt.Errorf("invalid public key format")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("cannot parse public key: %s", err.Error())
	}
	m.publickKey = publicKey.(*rsa.PublicKey)

	return nil
}

func (m *Manager) NewJWT(input AuthInfo) (string, error) {
	if m.privateKey == nil {
		return "", fmt.Errorf("no private key provided")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &tokenClaims{
		input.UserID,
		input.Login,
		input.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString(m.privateKey)
}

func (m *Manager) Parse(accessToken string) (result AuthInfo, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
		}
		return m.publickKey, nil
	})
	if err != nil {
		return result, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return result, errors.New(ErrTokenIsInvalid)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return result, errors.New(ErrClaims)
	}

	result.UserID = claims.UserID
	result.Login = claims.Login
	result.Role = claims.Role

	return result, nil
}

func (m *Manager) RefreshToken() (string, error) {
	token := jwt.New(m.signingMethod)

	tokenString, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) ValidateRefreshToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("token is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != m.signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.publickKey, nil
	})

	if err != nil {
		return false, err
	}

	if token.Valid {
		return true, nil
	}

	return false, errors.New("invalid token")
}
