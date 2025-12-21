package authconfig

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/LuukBlankenstijn/gewish/internal/logging"
	"github.com/LuukBlankenstijn/gewish/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func getSecretkey() []byte {
	return []byte(utils.EnvOrPanic("JWT_SECRET"))
}

type JWTClaims struct {
	LoggedIn         bool  `json:"logged_in"`
	SessionCreatedAt int64 `json:"session_created_at"`
	jwt.RegisteredClaims
}

type BasicStrategy struct {
	baseStrategy
	PasswordHash []byte `json:"password"`
}

func (BasicStrategy) Type() AuthStrategyType { return AuthStrategyBasic }

func NewBasicStrategy(password string) AuthStrategy {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &BasicStrategy{
		PasswordHash: hash,
	}
}

func (s *BasicStrategy) Authenticate(ctx context.Context, req connect.AnyRequest) (context.Context, error) {
	header := req.Header().Get("Cookie")
	cookies, err := http.ParseCookie(header)
	if err != nil || len(cookies) == 0 {
		return ctx, errors.New("failed to parse cookie from header")
	}
	tokenString := cookies[0].Value

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecretkey(), nil
	})

	if err != nil || !token.Valid {
		return ctx, errors.New("invalid token")
	}
	ctx = s.setContext(ctx, Admin, nil)
	return ctx, nil
}

func (s *BasicStrategy) Login(password string) (string, error) {
	if err := bcrypt.CompareHashAndPassword(s.PasswordHash, []byte(password)); err != nil {
		return "", err
	}

	sessionData := JWTClaims{
		LoggedIn:         true,
		SessionCreatedAt: time.Now().Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "gewi.sh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sessionData)
	tokenString, err := token.SignedString(getSecretkey())
	if err != nil {
		logging.Logger().Warn("failed to sign token", slog.Any("err", err))
		return "", errors.New("failed to generate token")
	}
	return tokenString, nil
}
