package controllers

import (
	"backend/models"
	"backend/services"
	"encoding/json"
	"github.com/aerogo/aero"
	"net/http"
)

type Response struct {
	Action  string  `json:"action"`
	Status  bool    `json:"status"`
}

func Login(ctx aero.Context) error {
	body, _ := ctx.Request().Body().JSONObject()

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		requestUser := new(models.User)
		err = json.Unmarshal(b, &requestUser)
		if err != nil {
			panic(err)
		}

		responseStatus, user := services.Login(requestUser)
		if responseStatus == http.StatusOK {
			rMap := map[string]string{
				"token": user.Token,
				"userName": user.UserName,
				"firstName": user.FirstName,
				"lastName": user.LastName,
				"email": user.Email,
				"role": user.Role.Role,
			}
			return ctx.JSON(rMap)
		}
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}

func Logout(ctx aero.Context) error {
	if ctx.Request().Method() == "OPTIONS" {
		return ctx.Error(http.StatusOK)
	}
	header := ctx.Request().Header("Authorization")
	if header != "" {
		err := services.Logout(ctx.Request().Internal())
		if err != nil {
			return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
		}
		response := &Response{
			Action: "logout",
			Status: true,
		}
		res, err := json.Marshal(response)
		if err != nil {
			return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
		}
		return ctx.JSON(string(res))
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}

func RefreshToken(ctx aero.Context) error {
	header := ctx.Request().Header("Authorization")
	body, _ := ctx.Request().Body().JSONObject()
	if body != nil && header != "" {
		// requestUser.Token = header
		b, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		requestUser := new(models.User)
		err = json.Unmarshal(b, &requestUser)
		if err != nil {
			panic(err)
		}

		responseStatus, token := services.RefreshToken(requestUser, ctx.Request().Internal())
		if responseStatus == http.StatusOK {
			return ctx.JSON(string(token))
		}
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}