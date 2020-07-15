package controllers

import (
	"encoding/json"
	"github.com/aerogo/aero"
	"net/http"

	"backend/helpers"
	"backend/models"
	"backend/services"
)

type UserDataResponse struct {
	Action  string  	 `json:"action"`
	Data  	*models.User `json:"data"`
}

func Register(ctx aero.Context) error {
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

		responseStatus := services.Register(requestUser)
		// rMap := map[string]bool{"register": responseStatus}
		response := &Response{
			Action: "register",
			Status: responseStatus,
		}
		res, err := json.Marshal(response)
		return ctx.JSON(res)
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}

func Profile(ctx aero.Context) error {
	switch rType := ctx.Request().Method(); rType {
	case "GET":
		requestUser := new(models.User)
		requestUser.UserName = ctx.Get("user")

		responseStatus, data := services.Profile(requestUser)
		response := &UserDataResponse{
			Action: "profile",
			Data: data,
		}

		res, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		if responseStatus == http.StatusOK {
			return ctx.JSON(string(res))
		}
	case "POST":
		// body, _ := ctx.Request().Body().JSONObject()
		body := ctx.Session().Get("body")
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

			responseStatus := services.ProfileUpdate(requestUser)
			response := &Response{
				Action: "profile",
				Status: responseStatus,
			}
			res, err := json.Marshal(response)
			return ctx.JSON(string(res))
		}
	}

	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}

func Settings(ctx aero.Context) error {
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

		responseStatus := services.Register(requestUser)
		// rMap := map[string]bool{"register": responseStatus}
		response := &Response{
			Action: "register",
			Status: responseStatus,
		}
		res, err := json.Marshal(response)
		return ctx.JSON(res)
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}

func Tasks(ctx aero.Context) error {
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

		responseStatus := services.Register(requestUser)
		// rMap := map[string]bool{"register": responseStatus}
		response := &Response{
			Action: "register",
			Status: responseStatus,
		}
		res, err := json.Marshal(response)
		return ctx.JSON(res)
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}

func Team(ctx aero.Context) error {
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

		responseStatus := services.Register(requestUser)
		// rMap := map[string]bool{"register": responseStatus}
		response := &Response{
			Action: "register",
			Status: responseStatus,
		}
		res, err := json.Marshal(response)
		return ctx.JSON(res)
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}

func Calendar(ctx aero.Context) error {
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

		responseStatus := services.Register(requestUser)
		// rMap := map[string]bool{"register": responseStatus}
		response := &Response{
			Action: "register",
			Status: responseStatus,
		}
		res, err := json.Marshal(response)
		return ctx.JSON(res)
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}

func Leads(ctx aero.Context) error {
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

		responseStatus := services.Register(requestUser)
		// rMap := map[string]bool{"register": responseStatus}
		response := &Response{
			Action: "register",
			Status: responseStatus,
		}
		res, err := json.Marshal(response)
		return ctx.JSON(res)
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}

func Users(ctx aero.Context) error {
	switch rType := ctx.Request().Method(); rType {
	case "POST":
		response := services.AllUsers()

		if response != nil {
			res := helpers.UserSliceToMap(response)
			return ctx.JSON(res)
		}
		return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
	case "OPTIONS":
		return ctx.Error(http.StatusOK)
	}
	return ctx.Error(http.StatusUnauthorized, "Not authorized", nil)
}
