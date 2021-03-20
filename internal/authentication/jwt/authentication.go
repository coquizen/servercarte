package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/authentication"
	"github.com/CaninoDev/gastro/server/internal/config"
	"github.com/CaninoDev/gastro/server/internal/model"
)

type JWT struct {
	secretKey        []byte
	algorithm        jwt.SigningMethod
	expirationPeriod time.Duration
	minSecretLength  int
}

type Settings struct {
	Algorithm        string
	ExpirationPeriod int
	MinKeyLength     int
	SecretKey        string
}

type jwtWrappedClaims struct {
	jwt.StandardClaims
	authentication.CustomClaims
}

func (j jwtWrappedClaims) Valid() error {
	if j.AccountID == uuid.Nil {
		return errors.New("empty claims")
	}
	if j.Username == "" {
		return errors.New("no username specified")
	}
	if j.Expiry == 0 {
		return errors.New("no expiration period specified")
	}
	return nil
}


// minSecretLen is the fallback secret key-length in case configuration did not declare such a length
var minSecretLen = 32

// Bind binds the supplied settings to an instance of jwt.
func Bind(algorithm jwt.SigningMethod, secKey []byte, period time.Duration) *JWT {
	return &JWT{algorithm: algorithm, expirationPeriod: period, minSecretLength: minSecretLen, secretKey: secKey}
}

// New returns a configured instance of JWT
func New(cfg config.Authentication) (*JWT, error) {
	if cfg.MinKeyLength > 0 {
		minSecretLen = cfg.MinKeyLength
	}
	// if len(cfg.SecretKey) < minSecretLen {
	// 	return &JWT{}, fmt.Errorf("jwt secret length is %v, which is less than required %v", len(cfg.SecretKey),
	// 		minSecretLen)
	// }
	signingMethod := jwt.GetSigningMethod(cfg.Algorithm)
	if signingMethod == nil {
		return &JWT{}, fmt.Errorf("invalid algorithm; given %v without SigningMethod interface implemented",
			cfg.Algorithm)
	}

	secKey := []byte(cfg.SecretKey)
	period := time.Duration(cfg.ExpirationPeriod) * time.Minute
	return Bind(signingMethod, secKey, period), nil
}

// GenerateToken generates a token with cstClaims encoded
func (a *JWT) GenerateToken(ctx context.Context, acct *model.Account) (string, error) {
	now := time.Now().Unix()
	exp := time.Now().Add(a.expirationPeriod).Unix()
	cstClaims := authentication.CustomClaims{
		AccountID: acct.ID,
		Username:  acct.Username,
		Role:      acct.Role,
		Expiry:    exp,
	}
	//wrappedCstClaims := &jwtTokenizer(claims)
	// add in standard claims to take advantage of built in methods.
	wrappedClaims := &jwtWrappedClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
			IssuedAt: now,
		},
		CustomClaims: cstClaims,
	}
	token := jwt.NewWithClaims(a.algorithm, wrappedClaims)
	signedToken, err := token.SignedString(a.secretKey)
	if err != nil {
		return signedToken, err
	}
	return signedToken, err
}

// verifyToken verifies the token is the appropriate format and can be decoded
func (a *JWT) verifyToken(req *http.Request) (*jwt.Token, error) {
	tokenString := a.ExtractToken(req)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if token.Method != a.algorithm {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
		return a.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// TokenValid checks the token validity
func (a *JWT) TokenValid(req *http.Request) error {
	token, err := a.verifyToken(req)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}



// parseToken reads the header string and parses the token as encoded by GenerateToken
func (a *JWT) parseToken(req *http.Request) (*jwt.Token, error) {
	tokenString := a.ExtractToken(req)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if a.algorithm != token.Method {
			return nil, errors.New("could not decode token; please re-authenticate")
		}
		return a.secretKey, nil
	})
}
// ExtractToken extracts the token string from the request header
func (a *JWT) ExtractToken(req *http.Request) string {
	authorizatioHeader := req.Header.Get("Authorization")
	if authorizatioHeader != "" {
		bearerToken := strings.Split(authorizatioHeader, " ")
		if len(bearerToken) == 2 {
			return bearerToken[1]
		}
	}
	return ""
}

// ExtractTokenClaims extracts the claims as encoded in the token
func (a *JWT) ExtractTokenClaims(req *http.Request) (authentication.CustomClaims, error) {
	token, err := a.verifyToken(req)
	if err != nil {
		return authentication.CustomClaims{}, err
	}
	claims, ok := token.Claims.(jwtWrappedClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		return authentication.CustomClaims(claims.CustomClaims), nil
	}
	return authentication.CustomClaims{}, err
}

