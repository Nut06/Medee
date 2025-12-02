package authadapter

import (
	authport "backend/internal/port/auth"
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTService(secret string, accessTTL, refreshTTL time.Duration) *JWTService {
	if accessTTL == 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL == 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &JWTService{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *JWTService) GenerateAccess(ctx context.Context, userID uuid.UUID) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.accessTTL)
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
		"typ": "access",
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

func (s *JWTService) GenerateRefresh(ctx context.Context, userID uuid.UUID) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.refreshTTL)
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
		"typ": "refresh",
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

func (s *JWTService) DecodeAccessToken(ctx context.Context, tokenString string) (uuid.UUID, error) {
	return s.decodeToken(tokenString, "access")
}

func (s *JWTService) DecodeRefreshToken(ctx context.Context, tokenString string) (uuid.UUID, error) {
	return s.decodeToken(tokenString, "refresh")
}

func (s *JWTService) decodeToken(tokenString, expectedType string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secret, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if typ, ok := claims["typ"].(string); !ok || typ != expectedType {
			return uuid.Nil, jwt.ErrTokenInvalidClaims
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return uuid.Nil, jwt.ErrTokenInvalidClaims
		}

		return uuid.Parse(sub)
	}

	return uuid.Nil, jwt.ErrTokenInvalidId
}

var _ authport.TokenService = (*JWTService)(nil)

// Helper for Middleware
func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	// Note: In a real app, you should inject the secret or use the service instance.
	// For simplicity here, we might need to read env again or make this a method of JWTService if we can access the instance.
	// Since Middleware is static, let's read env for now or better, make a global/singleton verifier.

	// Better approach: Let's assume we use the same secret from env
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, jwt.ErrTokenInvalidId
}
