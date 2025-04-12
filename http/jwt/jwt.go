package jwt

import (
	"errors"
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/http/session"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

// TODO: move to config
var (
	secret          = "test"
	errInvalidToken = errors.New("invalid token")
)

func GenerateJwt(ses *session.Session, ttl time.Duration, isRefresh bool) (string, error) {
	claims := ses.Claims()
	if claims == nil {
		claims = make(map[string]any)
	}
	for k, v := range claims {
		if v == nil {
			delete(claims, k)
		}
	}

	claims["user_id"] = ses.UserID()
	claims["email"] = ses.Email()
	claims["account_role"] = ses.AccountRole().String()
	claims["is_refresh"] = isRefresh
	claims["exp"] = time.Now().UTC().Add(ttl).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJwt(tokenString string) (*session.Session, bool, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimPrefix(tokenString, "bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, false, errInvalidToken
	}

	return getSession(claims)
}

func getSession(claims jwt.MapClaims) (*session.Session, bool, error) {
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, false, errInvalidToken
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, false, errInvalidToken
	}
	role, ok := claims["account_role"].(string)
	if !ok {
		return nil, false, errInvalidToken
	}
	accountRole, err := enum.AccountRoleString(role)
	if err != nil {
		return nil, false, errInvalidToken
	}
	isRefresh, ok := claims["is_refresh"].(bool)
	if !ok {
		return nil, false, errInvalidToken
	}
	ses := session.NewSession(userID, email, accountRole, claims)
	return ses, isRefresh, nil
}
