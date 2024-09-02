package services

import (
	"context"
	"errors"
	"motorbike-rental-backend/internal/app/user-and-auth/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAuthService interface {
	GenerateTokenPair(userID int64, refreshTokenID uuid.UUID, role float64) (models.AuthTokenPair, error)
	ParseRefreshToken(refreshToken string) (refreshTokenID uuid.UUID, userID int64, role float64, err error)
	GetAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID) (models.AuthRefreshToken, error)
	CreateAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID, userID int64, role float64) error
	UpdateAuthRefreshTokenExpires(ctx context.Context, refreshToken models.AuthRefreshToken) error
	DeleteAuthRefreshToken(ctx context.Context, userID int64) error
}

type AuthService struct {
	DB                     *gorm.DB
	jwtSecret              string
	accessTokenExpireTime  time.Duration
	refreshTokenExpireTime time.Duration
}

func NewAuthService(db *gorm.DB, jwtSecret string, accessTokenExpireTime, refreshTokenExpireTime time.Duration) IAuthService {
	return &AuthService{
		DB:                     db,
		jwtSecret:              jwtSecret,
		accessTokenExpireTime:  accessTokenExpireTime,
		refreshTokenExpireTime: refreshTokenExpireTime,
	}
}

func (s AuthService) ParseRefreshToken(refreshToken string) (refreshTokenID uuid.UUID, userID int64, role float64, err error) {
	// Parse the refresh token.
	refreshClaims := refreshTokenClaims{}
	claims, err := jwt.ParseWithClaims(refreshToken, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return
	}
	if !claims.Valid {
		err = errors.New("invalid token")
		return
	}

	// Verify that the refresh token is not expired
	now := time.Now()
	rtokenClaims := claims.Claims.(*refreshTokenClaims)
	if rtokenClaims.VerifyExpiresAt(now, false) == false {
		err = errors.New("token expired")
		return
	}

	return rtokenClaims.ID, rtokenClaims.UserID, rtokenClaims.Role, nil
}

type accessTokenClaims struct {
	jwt.RegisteredClaims
	ID     uuid.UUID `json:"id"`
	UserID int64     `json:"uid"`
	Role   float64   `json:"role"`
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
	ID     uuid.UUID `json:"id"`
	UserID int64     `json:"uid"`
	Role   float64   `json:"role"`
}

func (s AuthService) GenerateTokenPair(userID int64, refreshTokenID uuid.UUID, role float64) (models.AuthTokenPair, error) {
	var err error
	var m models.AuthTokenPair
	now := time.Now()

	accessClaims := accessTokenClaims{
		ID:     uuid.New(),
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpireTime)),
		},
	}

	m.AccessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return m, err
	}

	refreshClaims := refreshTokenClaims{
		ID:     refreshTokenID,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenExpireTime)),
		},
	}

	m.RefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return m, err
	}

	return m, nil
}

func (s AuthService) GetAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID) (models.AuthRefreshToken, error) {
	var refreshToken models.AuthRefreshToken
	err := s.DB.WithContext(ctx).Where("token_id = ?", refreshTokenID).First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return refreshToken, errors.New("refresh token not found")
		}
		return refreshToken, errors.New("failed to retrieve auth refresh token: " + err.Error())
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return refreshToken, errors.New("refresh token expired")
	}
	return refreshToken, nil
}

func (s AuthService) CreateAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID, userID int64, role float64) error {
	refreshToken := models.AuthRefreshToken{
		TokenID:   refreshTokenID,
		UserID:    userID,
		Role:      role,
		ExpiresAt: time.Now().Add(s.refreshTokenExpireTime),
	}
	err := s.DB.WithContext(ctx).Create(&refreshToken).Error
	if err != nil {
		return errors.New("failed to create auth refresh token: " + err.Error())
	}
	return nil
}

func (s AuthService) UpdateAuthRefreshTokenExpires(ctx context.Context, refreshToken models.AuthRefreshToken) error {
	err := s.DB.WithContext(ctx).Model(&models.AuthRefreshToken{}).
		Where("token_id = ?", refreshToken.TokenID).
		Update("expires_at", time.Now().Add(s.refreshTokenExpireTime)).Error
	if err != nil {
		return errors.New("failed to update auth refresh token expiration: " + err.Error())
	}
	return nil
}
func (s AuthService) DeleteAuthRefreshToken(ctx context.Context, userID int64) error {
	err := s.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.AuthRefreshToken{}).Error
	if err != nil {
		return errors.New("failed to delete auth refresh token: " + err.Error())
	}
	return nil
}
