package tokenservice

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenService handles JWT token generation and validation
type jwtTokenService struct {
	accessSecret    string
	refreshSecret   string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewTokenService creates a new TokenService instance
func NewTokenService(accessSecret, refreshSecret string, accessTokenTTL, refreshTokenTTL time.Duration) *jwtTokenService {
	return &jwtTokenService{
		accessSecret:    accessSecret,
		refreshSecret:   refreshSecret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// TokenPair contains both access and refresh tokens
type TokenPair struct {
	Access  Token `json:"access_token"`
	Refresh Token `json:"refresh_token"`
}

// Token represents a token with its TTL
type Token struct {
	Token string
	TTL   time.Duration
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateTokenPair generates both access and refresh tokens
func (ts *jwtTokenService) GenerateTokenPair(userID uint) (*TokenPair, error) {
	accessToken, err := ts.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := ts.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		Access:  Token{Token: accessToken, TTL: ts.accessTokenTTL},
		Refresh: Token{Token: refreshToken, TTL: ts.refreshTokenTTL},
	}, nil
}

// GenerateAccessToken creates a new JWT access token
func (ts *jwtTokenService) GenerateAccessToken(userID uint) (string, error) {
	if ts.accessSecret == "" {
		return "", errors.New("signing key is empty")
	}

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ts.accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ts.accessSecret))
}

// GenerateRefreshToken creates a secure random refresh token
func (ts *jwtTokenService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ValidateToken validates the JWT token and returns claims
func (ts *jwtTokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(ts.accessSecret), nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			}
		}
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
