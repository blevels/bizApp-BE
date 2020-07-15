package services

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
	"net/http"

	"backend/api/parameters"
	"backend/core/authentication"
	"backend/models"
)

func Login(requestUser *models.User) (int, *models.User) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	authUser := authBackend.Authenticate(requestUser)
	if authUser != nil {
		token, err := authBackend.GenerateToken(authUser.UUID.String())
		if err != nil {
			return http.StatusInternalServerError, new(models.User)
		} else {
			authUser.Token = parameters.TokenAuthentication{token}.Token
			return http.StatusOK, authUser
		}
	}
	return http.StatusUnauthorized, new(models.User)
}

func RefreshToken(requestUser *models.User, req *http.Request) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()

	keyfunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	}

	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, keyfunc)
	tokenString := req.Header.Get("Authorization")

	if err == nil && token.Valid && !authBackend.IsInBlacklist(tokenString) {
		authUser, _ := authBackend.CheckUser(requestUser)
		if authUser != nil {
			//if authUser != nil && tokenString == authUser.Token{
			token, err := authBackend.GenerateToken(authUser.UUID.String())
			if err != nil {
				panic(err)
			}
			response, err := json.Marshal(parameters.TokenAuthentication{token})
			if err != nil {
				panic(err)
			}

			return http.StatusOK, response
		}
	}
	return http.StatusUnauthorized, []byte("")
}

func Logout(req *http.Request) error {
	authBackend := authentication.InitJWTAuthenticationBackend()

	keyfunc := func(*jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	}

	tokenRequest, err := request.ParseFromRequest(req, request.OAuth2Extractor, keyfunc)
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authBackend.Logout(tokenString, tokenRequest)
}