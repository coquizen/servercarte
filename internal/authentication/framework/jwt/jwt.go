package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CaninoDev/gastro/server/authentication"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/config"
)

var (
	ErrMalformedToken = errors.New("token could not be parsed")
	ErrNonExistentToken = errors.New("no token found in Authorization header")
)

var NullCustomClaims = authentication.CustomClaims{}


// adapter is an authentication adapter
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
	if j.Username == "" {
		return errors.New("empty username")
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
func (s *adapter) GenerateToken(_ context.Context, accountID uuid.UUID, username string, accessLevel int) (string, error) {

	exp := time.Now().Add(s.expirationPeriod).Unix()
	fmt.Printf("accountID: %v", accountID)
	cstClaims := &claims{
		AccountID: accountID,
		Username: username,
		Role: accessLevel,
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
// 			return nil, errors.CreateEmployeeAccount("could not decode token; please re-authenticate")
// 		}
// 		return s.secretKey, nil
// 	})
// }

// ExtractToken extracts the token string from the request header
func (s *adapter) ExtractToken(req *http.Request) (string, error) {
	authorizationHeader := req.Header.Get("Authorization")
	if authorizationHeader != "" {
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) == 2 {
			return bearerToken[1], nil
		}
	}
	return "", ErrNonExistentToken
}

// ParseTokenClaims extracts the claims as encoded in the token
func (s *adapter) ParseTokenClaims(tokenString string) (authentication.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != s.algorithm {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	if err != nil {
		return NullCustomClaims, ErrMalformedToken
	}
	c := token.Claims.(*claims)
	if err := c.Valid(); err != nil {
		return NullCustomClaims, err
	}
	return authentication.CustomClaims{
		AccountID: c.AccountID,
		Username: c.Username,
		Role: c.Role,
		Expiry: c.Expiry}, nil
}