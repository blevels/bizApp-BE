package authentication

import (
	"bufio"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"

	"backend/core/redis"
	"backend/models"
	"backend/settings"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil
var dbHelper models.DatabaseHelper

func init() {
	dbHelper = models.CreateDBHelperService()
}

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

func (backend *JWTAuthenticationBackend) GenerateToken(userUUID string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix(),
		"iat": time.Now().Unix(),
		"sub": userUUID,
	}
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}

	requestUser := new(models.User)
	requestUser.Token = tokenString
	requestUser.UUID, err = uuid.FromString(userUUID)
	if err != nil {
		panic(err)
		return "", err
	}

	res := backend.setUserToken(requestUser)
	if res {
		return tokenString, nil
	}
	return "", errors.New("Token not generated")
}

func (backend *JWTAuthenticationBackend) ValidateUserToken(user *models.User, token string) bool {
	checkUser, _ := backend.CheckUser(user)
	if checkUser != nil && checkUser.Token == token {
		return true
	}
	return false
}

func (backend *JWTAuthenticationBackend) Authenticate(user *models.User) *models.User {
	checkUser, e := backend.CheckUser(user)
	if !e {
		return nil
	}
	if checkUser != nil && checkUser.UserName == user.UserName && bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(user.Password)) == nil {
		return checkUser
	}
	return nil
}

func (backend *JWTAuthenticationBackend) Registration(user *models.User) bool {
	_, response := backend.CheckUser(user)
	if response {
		return false
	}

	userDba := models.NewUserDatabase(dbHelper)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword)
	user.Role.Role = "USER"

	if user.UUID.String() == "" || user.UUID.String() == "00000000-0000-0000-0000-000000000000" {
		user.UUID = uuid.NewV4()
	}

	err := userDba.Create(context.TODO(), user)
	if err != nil {
		return false
	}
	return true
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	redisConn := redis.Connect()
	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims.(jwt.MapClaims)["exp"]))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisConn := redis.Connect()
	redisToken, _ := redisConn.GetValue(token)

	if redisToken == nil {
		return false
	}

	return true
}

func (backend *JWTAuthenticationBackend) CheckUser(user *models.User) (*models.User, bool) {
	filter := bson.M{"username": user.UserName}
	userDba := models.NewUserDatabase(dbHelper)
	dbUser, _ := userDba.FindOne(context.TODO(), filter)
	if dbUser != nil {
		return dbUser, true
	}
	return &models.User{}, false
}

func (backend *JWTAuthenticationBackend) Profile(user *models.User) (*models.User, bool) {
	filter := bson.M{"username": user.UserName}
	userDba := models.NewUserDatabase(dbHelper)
	dbUser, _ := userDba.FindOne(context.TODO(), filter)
	if dbUser != nil {
		return dbUser, true
	}
	return &models.User{}, false
}

func (backend *JWTAuthenticationBackend) ProfileUpdate(user *models.User) (bool) {
	userDba := models.NewUserDatabase(dbHelper)
	filter := bson.M{"username": user.UserName}
	update := bson.D{
		{
			"$set",
			bson.D{
				{"firstname", user.FirstName},
				{"lastname", user.LastName},
				{"email", user.Email},
			},
		},
	}

	err := userDba.Update(context.TODO(), filter, update)
	if err != nil {
		return false
	}
	return true
}

func (backend *JWTAuthenticationBackend) GetAllUsers() *models.Users {
	filter := bson.D{}
	userDba := models.NewUserDatabase(dbHelper)
	users, _ := userDba.Find(context.TODO(), filter)
	return users
}

func (backend *JWTAuthenticationBackend) updateAuth(user *models.User) bool {
	userDba := models.NewUserDatabase(dbHelper)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword)

	if user.UUID.String() == "" || user.UUID.String() == "00000000-0000-0000-0000-000000000000" {
		user.UUID = uuid.NewV4()
	}

	// opts := options.Update().SetUpsert(true)
	filter := bson.M{"username": user.UserName}
	update := bson.D{{"$set", bson.D{{"Password", user.Password}}}}

	err := userDba.Update(context.TODO(), filter, update)
	if err != nil {
		return false
	}
	return true
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTAuthenticationBackend) setUserToken(user *models.User) bool {
	userDba := models.NewUserDatabase(dbHelper)

	if user.UUID.String() == "" || user.UUID.String() == "00000000-0000-0000-0000-000000000000" || user.Token == "" {
		return false
	}

	// opts := options.Update().SetUpsert(true)
	filter := bson.M{"uuid": user.UUID}
	update := bson.D{{"$set", bson.D{{"token", user.Token}}}}

	err := userDba.Update(context.TODO(), filter, update)
	if err != nil {
		return false
	}
	return true
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(settings.Get().PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(settings.Get().PublicKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}