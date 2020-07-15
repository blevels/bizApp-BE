package authentication

import (
	"encoding/json"
	"fmt"
	"github.com/aerogo/aero"
	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
	"net/http"

	"backend/models"
)

func RequireTokenAuthentication(next aero.Handler) aero.Handler {
	return func(ctx aero.Context) error {
		if ctx.Request().Method() == "OPTIONS" {
			return next(ctx)
		}

		path := ctx.Request().Path()
		if path == "/api/1.0/login" || path ==  "/api/1.0/register" || path == "/api/1.0/logout" {
			return next(ctx)
		}

		authBackend := InitJWTAuthenticationBackend()

		keyfunc := func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			} else {
				return authBackend.PublicKey, nil
			}
		}

		token, err := request.ParseFromRequest(ctx.Request().Internal(), request.OAuth2Extractor, keyfunc)
		tokenString := ""

		authHeader := ctx.Request().Header("Authorization")
		if authHeader != "" {
			tokenString = authHeader[7:len(authHeader)]
		}

		if err == nil && token.Valid && !authBackend.IsInBlacklist(ctx.Request().Header("Authorization")) {
			requestUser := new(models.User)

			if ctx.Get("user") != "" {
				requestUser.UserName = ctx.Get("user")
			} else {
				body, _ := ctx.Request().Body().JSON()
				if body != nil {
					ctx.Session().Set("body", body)
					b, err := json.Marshal(body)
					if err != nil {
						panic(err)
					}

					err = json.Unmarshal(b, &requestUser)
					if err != nil {
						panic(err)
					}
				}
			}

			res := authBackend.ValidateUserToken(requestUser, tokenString)
			if res != true {
				return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
			}
			return next(ctx)
		} else {
			return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
		}
	}
}

func RequireAdminRole(next aero.Handler) aero.Handler {
	return func(ctx aero.Context) error {
		if ctx.Request().Method() == "OPTIONS" {
			return next(ctx)
		}

		authBackend := InitJWTAuthenticationBackend()
		requestUser := new(models.User)

		body := ctx.Session().Get("body")
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(b, &requestUser)
			if err != nil {
				panic(err)
			}
		}

		user, res := authBackend.CheckUser(requestUser)
		if res && user.Role.Role == "ADMIN" {
			return next(ctx)
		}
		return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
	}
}

func Cors(next aero.Handler) aero.Handler {
	return func(ctx aero.Context) error {
		ctx.Response().SetHeader("Access-Control-Allow-Origin","*")
		ctx.Response().SetHeader("Access-Control-Allow-Methods","POST,GET,OPTIONS,PUT,DELETE,PATCH,HEAD")
		ctx.Response().SetHeader("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		ctx.Response().SetHeader("Access-Control-Allow-Credentials","true")
		ctx.Response().SetHeader("Access-Control-Max-Age","86400")
		ctx.Response().SetHeader("Connection","Keep-Alive")
		return next(ctx)
	}
}