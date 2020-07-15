package services

import (
	"backend/core/authentication"
	"backend/models"
	"net/http"
)

func Register(requestUser *models.User) (bool) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	if authBackend.Registration(requestUser) {
		return true
	}
	return false
}

func Profile(requestUser *models.User) (int, *models.User) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	profile, u := authBackend.Profile(requestUser)

	if u {
		profile.Password = ""
		return http.StatusOK, profile
	}

	return http.StatusUnauthorized, profile
}

func ProfileUpdate(requestUser *models.User) (bool) {
	authBackend := authentication.InitJWTAuthenticationBackend()

	if authBackend.ProfileUpdate(requestUser) {
		return true
	}
	return false
}

func Settings(requestUser *models.User) (int, *models.User) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	profile, u := authBackend.Profile(requestUser)

	if u {
		profile.Password = ""
		return http.StatusOK, profile
	}

	return http.StatusUnauthorized, profile
}

func SettingsUpdate(requestUser *models.User) (bool) {
	authBackend := authentication.InitJWTAuthenticationBackend()

	if authBackend.ProfileUpdate(requestUser) {
		return true
	}
	return false
}

func AllUsers() *models.Users  {
	authBackend := authentication.InitJWTAuthenticationBackend()
	users := authBackend.GetAllUsers()
	return users
}