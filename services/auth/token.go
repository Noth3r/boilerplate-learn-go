package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type tokenservice struct{}

func NewToken() *tokenservice {
	return &tokenservice{}
}

type TokenInterface interface {
	CreateToken(string, bool) (*TokenDetails, error)
	ExtractTokenMetadata(*http.Request) (*AccessDetails, error)
}

var _ TokenInterface = &tokenservice{}

func (t *tokenservice) CreateToken(userId string, admin bool) (*TokenDetails, error) {
	td := &TokenDetails{}

	td.AtExpires = time.Now().Add(time.Minute * 1).Unix()
	td.TokenUuid = uuid.NewString()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + userId

	var err error

	// token := jwt.New(jwt.SigningMethodHS256)
	// claims := token.Claims.(jwt.MapClaims)
	// claims["access_uuid"] = td.TokenUuid
	// claims["user_id"] = userId
	// claims["exp"] = td.AtExpires
	// t, err := token.SignedString([]byte("secret"))

	atClaims := jwt.MapClaims{}
	atClaims["admin"] = admin
	atClaims["access_uuid"] = td.TokenUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte("secret"))

	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte("secret"))

	if err != nil {
		return nil, err
	}

	return td, nil
}

func TokenValid(r *http.Request) (jwt.Claims, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.Claims)
	if !ok && !token.Valid {
		return nil, err
	}

	return claims, nil
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	return token, err
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func extract(token *jwt.Token) (*AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		userId, userOk := claims["user_id"].(string)
		if ok == false || userOk == false {
			return nil, errors.New("Unathorized")
		} else {
			return &AccessDetails{
				TokenUuid: accessUuid,
				UserId:    userId,
			}, nil
		}
	}
	return nil, errors.New("Something went wrong")
}

func (t *tokenservice) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	acc, err := extract(token)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
