package unit_tests

import (
	jwt "github.com/dgrijalva/jwt-go"
	"os"
	"testing"

	"backend/core/authentication"
	"backend/core/redis"
	"backend/models"
	"backend/settings"
)

func init () {
	os.Setenv("GO_ENV", "tests")
	settings.Init()
}

func TestInitJWTAuthenticationBackend(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	if authBackend == nil {
		t.Fatalf("token is empty")
	}
	if authBackend.PublicKey == nil {
		t.Fatalf("token is empty")
	}
}

func TestGenerateToken(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken("1234")
	if err != nil {
		t.Fatal(err)
	}
	if tokenString == "" {
		t.Fatalf("tokenString is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if token.Valid != true {
		t.Fatalf("expected true, got '%T'", token.Valid)
	}
}

func TestAuthenticate(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	user := &models.User{
		Username: "haku",
		Password: "testing",
	}
	auth := authBackend.Authenticate(user)
	if auth != true {
		t.Fatalf("expected true, got '%T'", auth)
	}
}

func TestAuthenticateIncorrectPass(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	user := models.User{
		Username: "haku",
		Password: "Password",
	}

	auth := authBackend.Authenticate(&user)
	if auth != false {
		t.Fatalf("expected false, got '%T'", auth)
	}
}

func TestAuthenticateIncorrectUsername(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	user := &models.User{
		Username: "Haku",
		Password: "testing",
	}
	auth := authBackend.Authenticate(user)
	if auth != false {
		t.Fatalf("expected false, got '%T'", auth)
	}
}

func TestLogout(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken("1234")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	err = authBackend.Logout(tokenString, token)
	if err != nil {
		t.Fatal(err)
	}

	redisConn := redis.Connect()
	redisValue, err := redisConn.GetValue(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	if redisValue == nil {
		t.Fatalf("expected a value, got '%s'", redisValue)
	}
}

func TestIsInBlacklist(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenString, err := authBackend.GenerateToken("1234")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	err = authBackend.Logout(tokenString, token)
	if err != nil {
		t.Fatal(err)
	}

	bl := authBackend.IsInBlacklist(tokenString)
	if bl != true {
		t.Fatalf("expected true, got '%T'", bl)
	}
}

func TestIsNotInBlacklist(t *testing.T) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	bl := authBackend.IsInBlacklist("1234")
	if bl != false {
		t.Fatalf("expected false, got '%T'", bl)
	}
}
