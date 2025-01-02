package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func newToken(userId string, typeToken *ModifierToken) (*Claims, string, error) {
	expirationTime := time.Now().Add(typeToken.Expires)
	jwtExpirationTime := jwt.NewNumericDate(expirationTime)
	claims := &Claims{
		UUID: uuid.New(),
		GUID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwtExpirationTime,
			Issuer:    "bank account holder",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(typeToken.Key)
	if err != nil {
		return nil, "", err
	}

	return claims, tokenString, nil
}

func GeneratePairToken(userId string) (*http.Cookie, *http.Cookie, error) {
	accessClaims, accessToken, err := newToken(userId, AccessToken)
	if err != nil {
		slog.Error("error maintain access_token", slog.String("err", err.Error()))
		return nil, nil, err
	}

	refreshClaims, refreshToken, err := newToken(userId, RefreshToken)
	if err != nil {
		slog.Error("error maintain refresh_token", slog.String("err", err.Error()))
		return nil, nil, err
	}

	// "bcrypt: password length exceeds 72 bytes"
	// refreshTokenCript, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

	theRedis[userId] = PairToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &http.Cookie{
			Name:    AccessToken.Name,
			Path:    "/",
			Value:   accessToken,
			Expires: accessClaims.ExpiresAt.Time,
		}, &http.Cookie{
			Name:    RefreshToken.Name,
			Path:    "/",
			Value:   refreshToken,
			Expires: refreshClaims.ExpiresAt.Time,
		}, err
}

func ValidateToken(tokenArrived string) (string, error) {
	claims, err := parseToken(tokenArrived)
	if err != nil {
		return "", fmt.Errorf("parseToken: %w", err)
	}

	pair := theRedis[claims.GUID]
	if pair.AccessToken != tokenArrived && pair.RefreshToken != tokenArrived {
		return "", fmt.Errorf("reddis: token not found")
	}

	return claims.GUID, nil
}

func MaintainToken(refreshtoken string, accesstoken string) (*http.Cookie, error) {
	logger := slog.With(
		slog.String("konponent", "jwt.MaintainToken"),
	)

	refreshclaims, err := parseToken(refreshtoken)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	accessclaims, err := parseToken(accesstoken)
	switch {
	case err != nil && err.Error() == "token has invalid claims: token is expired":
		logger.Error(err.Error())
	case err != nil:
		logger.Error(err.Error())
		return nil, err
	}

	fmt.Println(time.Until(accessclaims.ExpiresAt.Time))

	almostExpired := time.Until(accessclaims.ExpiresAt.Time)
	// нет смысла обновлять раньше чем за 5 минут до истечения
	if almostExpired > 5*time.Minute {
		err := errors.New("too little time has passed since the token was created")
		logger.Error(err.Error(), slog.String("token hasn't expired", almostExpired.String()))
		return nil, err
	}

	var userID string
	if refreshclaims.GUID != accessclaims.GUID {
		err := errors.New("tokens are not linked")
		logger.Error(err.Error())
		return nil, err
	} else {
		userID = refreshclaims.GUID
	}

	pair := theRedis[userID]

	if pair.RefreshToken != refreshtoken || pair.AccessToken != accesstoken {
		return nil, errors.New("token not found")
	}

	newclaims, newtoken, err := newToken(userID, AccessToken)
	if err != nil {
		return nil, fmt.Errorf("maintain token: %w", err)
	}
	pair.AccessToken = newtoken
	theRedis[userID] = pair

	return &http.Cookie{
		Name:    AccessToken.Name,
		Path:    "/",
		Value:   newtoken,
		Expires: newclaims.ExpiresAt.Time,
	}, err
}

func parseToken(tokenArrived string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenArrived,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(accessKey), nil
		})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return claims, err
	}
}
