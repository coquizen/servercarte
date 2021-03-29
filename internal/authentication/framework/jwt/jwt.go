package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CaninoDev/gastro/server/api/authentication"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/config"
)

// JWT is an authentication adapter
type adapter struct {
	secretKey        []byte
	algorithm        jwt.SigningMethod
	expirationPeriod time.Duration
	minSecretLength  int
}

// claims is a local alias for authentication.CustomClaims. To facilitate satisfaction
// of a contract as specified in jwt-go
type claims authentication.CustomClaims

// Valid satisfies the contract to make CustomClaims a jwt.Claims object
// (See: https://github.com/dgrijalva/jwt-go/blob/dc14462fd58732591c7fa58cc8496d6824316a82/claims.go#L9 )
func (j claims) Valid() error {
	if j.AccountID == uuid.Nil {
		return errors.New("empty claims")
	}
	if j.Expiry == 0 {
		return errors.New("no expiration period specified")
	}
	if isTokenExpired(j.Expiry) {
		return authentication.ErrExpiredToken
	}
	return nil
}

func isTokenExpired(tokenExpiry int64) bool {
	exp := time.Unix(tokenExpiry, 0).UTC()
	currTime := time.Now().UTC()
	return currTime.After(exp)
}
//
// func isValidRole(role account.AccessLevel) bool {
// 	switch role {
// 	case account.Admin, account.Employee, account.Guest:
// 		return true
// 	default:
// 		return false
// 	}
// }

// minSecretLen is the fallback secret key-length in case configuration did not declare such a length
var minSecretLen = 32

// New returns a configured instance.
func New(cfg config.Authentication) (*adapter, error) {
	if cfg.MinKeyLength > 0 {
		minSecretLen = cfg.MinKeyLength
	}

	signingMethod := jwt.GetSigningMethod(cfg.Algorithm)
	if signingMethod == nil {
		return &adapter{}, fmt.Errorf("invalid algorithm; given %v without SigningMethod interface implemented",
			cfg.Algorithm)
	}
	secKey := []byte(cfg.SecretKey)
	period := time.Duration(cfg.ExpirationPeriod) * time.Minute
	return &adapter{algorithm: signingMethod, expirationPeriod: period, minSecretLength: minSecretLen, secretKey: secKey}, nil
}

// GenerateToken generates a token with claims encoded
func (s *adapter) GenerateToken(_ context.Context, accountID uuid.UUID) (string, error) {
	exp := time.Now().Add(s.expirationPeriod).Unix()
	fmt.Printf("accountID: %v", accountID)
	cstClaims := &claims{
		AccountID: accountID,
		Expiry:    exp,
	}

	token := jwt.NewWithClaims(s.algorithm, cstClaims)
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return signedToken, err
	}
	return signedToken, err
}

// // parseToken reads the header string and parses the token as encoded by GenerateToken
// func (s *adapter) parseToken(tokenString string) (*jwt.Token, error) {
// 	return jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
// 		if s.algorithm != token.Method {
// 			return nil, errors.New("could not decode token; please re-authenticate")
// 		}
// 		return s.secretKey, nil
// 	})
// }

// ExtractRawsToken extracts the token string from the request header
func (s *adapter) ExtractRawToken(req *http.Request) (string, error) {
	authorizationHeader := req.Header.Get("Authorization")
	if authorizationHeader != "" {
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) == 2 {
			return bearerToken[1], nil
		}
	}
	return "", nil
}

// ExtractClaims extracts the claims as encoded in the token
func (s *adapter) ExtractClaims(req *http.Request) (uuid.UUID, error) {
	tokenString, err := s.ExtractRawToken(req)
	if err != nil {
		return uuid.Nil, err
	}
	token, err := s.verifyToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	if err := s.TokenValid(tokenString); err != nil {
		return uuid.Nil, err
	}
	accountID := token.Claims.(*claims).AccountID
	// the token claims should conform to MapClaims
	return accountID, nil
}

// TokenValid checks the token validity
func (s *adapter) TokenValid(tokenString string) error {
	token, err := s.verifyToken(tokenString)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(*claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// verifyToken verifies the token is the appropriate format and can be decoded
func (s *adapter) verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if token.Method != s.algorithm {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}