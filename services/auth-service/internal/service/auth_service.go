package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/example/microservices/auth-service/internal/model"
	"github.com/example/microservices/auth-service/internal/repository"
	appErrors "github.com/example/microservices/shared/pkg/errors"
)

type AuthService struct {
	repo                  *repository.UserRepository
	jwtSecret             string
	tokenExpiration        int
	refreshTokenExpiration int
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string, tokenExpiration int, refreshTokenExpiration int) *AuthService {
	return &AuthService{
		repo:                  repo,
		jwtSecret:             jwtSecret,
		tokenExpiration:        tokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}
}

func (s *AuthService) SignUp(ctx context.Context, req *model.SignUpRequest) (*model.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, appErrors.NewValidationError("Email and password are required", nil)
	}

	exists, err := s.repo.UserExists(ctx, req.Email)
	if err != nil {
		return nil, appErrors.NewInternalError("Failed to check user existence", err)
	}
	if exists {
		return nil, appErrors.NewConflictError("User with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErrors.NewInternalError("Failed to hash password", err)
	}

	user := &model.User{
		ID:    uuid.New().String(),
		Email: req.Email,
		Password: string(hashedPassword),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, appErrors.NewInternalError("Failed to create user", err)
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) SignIn(ctx context.Context, req *model.SignInRequest) (*model.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, appErrors.NewValidationError("Email and password are required", nil)
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, appErrors.NewAuthError("Invalid email or password", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, appErrors.NewAuthError("Invalid email or password", err)
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) generateAuthResponse(user *model.User) (*model.AuthResponse, error) {
	accessToken, err := s.generateToken(user.ID, s.tokenExpiration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, s.refreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.tokenExpiration,
		User:         *user,
	}, nil
}

func (s *AuthService) generateToken(userID string, expirationSeconds int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(expirationSeconds) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", appErrors.NewInternalError("Failed to generate token", err)
	}

	return tokenString, nil
}

func (s *AuthService) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", appErrors.NewAuthError("Invalid or expired token", err)
	}

	claims := token.Claims.(*jwt.MapClaims)
	userID, ok := (*claims)["sub"].(string)
	if !ok {
		return "", appErrors.NewAuthError("Invalid token claims", nil)
	}

	return userID, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	userID, err := s.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, appErrors.NewAuthError("User not found", err)
	}

	return s.generateAuthResponse(user)
}
