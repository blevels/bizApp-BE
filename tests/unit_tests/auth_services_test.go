package unit_tests

import (
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"testing"

	"backend/core/authentication"
	"backend/models"
	"backend/services"
	"backend/settings"
)

func init () {
	os.Setenv("GO_ENV", "tests")
	settings.Init()
}

func TestLogin(t *testing.T) {
	user := models.User{
		Username: "haku",
		Password: "testing",
	}
	response, token := services.Login(&user)
	if response != http.StatusOK {
		t.Fatalf("expected 200, got %d", response)
	}
	if token == nil {
		t.Fatalf("token is empty")
	}
}

func TestLoginIncorrectPassword(t *testing.T) {
	user := models.User{
		Username: "haku",
		Password: "Password",
	}
	response, token := services.Login(&user)
	if response != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized, got %d", response)
	}
	if token == nil {
		t.Fatalf("token is empty")
	}
}

func TestLoginIncorrectUsername(t *testing.T) {
	user := models.User{
		Username: "Username",
		Password: "testing",
	}
	response, token := services.Login(&user)
	if response != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized, got %d", response)
	}
	if token == nil {
		t.Fatalf("token is empty")
	}
}

func TestLoginEmptyCredentials(t *testing.T) {
	user := models.User{
		Username: "",
		Password: "",
	}
	response, token := services.Login(&user)
	if response != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized, got %d", response)
	}
	if token == nil {
		t.Fatalf("token is empty")
	}
}

func TestRefreshToken(t *testing.T) {
	user := models.User{
		Username: "haku",
		Password: "testing",
	}
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken(user.UUID.String())
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if token == nil {
		t.Fatalf("token is empty")
	}
}

func TestRefreshTokenInvalidToken(t *testing.T) {
	user := models.User{
		Username: "",
		Password: "",
	}

	// token := jwt.New(jwt.GetSigningMethod("RS256"))
	newToken := services.RefreshToken(&user)
	if newToken == nil {
		t.Fatalf("token is empty")
	}
}

/*
func TestLogout(t *testing.T) {
	user := models.User{
		Username: "haku",
		Password: "testing",
	}
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken(user.UUID.String())
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})

	err = authBackend.Logout(tokenString, token)
	if err != nil {
		t.Fatal(err)
	}
}
*/