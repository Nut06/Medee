package authadapter

import (
	authport "backend/internal/port/auth"
	utils "backend/internal/utils"
	"context"
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
		accessTTL = utils.FifteenMin
	}
	if refreshTTL == 0 {
		refreshTTL = utils.SevenDays
	}
	return &JWTService{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *JWTService) GetSecret() []byte {
	return s.secret
}

func (s *JWTService) GenerateAccess(ctx context.Context, userID uuid.UUID) (string, error) {
	expiresAt := time.Now().Add(s.accessTTL).Unix()
	
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": expiresAt,
		"iat": time.Now().Unix(),
		"typ": "access",
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	return signed, nil
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


func (s *JWTService) ParseToken(tokenString string) (*jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return s.secret, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, jwt.ErrTokenInvalidId
}
