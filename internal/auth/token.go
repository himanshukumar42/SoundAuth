package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/himanshukumar42/soundauth/internal/models"
)

type JWTTokenManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTTokenManager(secret, issuer string, ttl time.Duration) *JWTTokenManager {
	return &JWTTokenManager{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

func (m *JWTTokenManager) GenerateToken(ctx context.Context, req models.GenerateTokenRequest) (*models.TokenPair, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":        req.UserID,
		"email":      req.Email,
		"tenant_id":  req.TenantID,
		"session_id": req.SessionID,
		"roles":      req.Roles,
		"scopes":     req.Scopes,
		"iss":        m.issuer,
		"iat":        now.Unix(),
		"exp":        now.Add(m.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(m.secret)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(m.ttl.Seconds()),
	}, nil
}

func (m *JWTTokenManager) VerifyToken(ctx context.Context, tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return &models.JWTClaims{
		UserID:    claims["sub"].(string),
		Email:     claims["email"].(string),
		TenantID:  claims["tenant_id"].(string),
		SessionID: claims["session"].(string),
	}, nil
}
