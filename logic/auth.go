package logic

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/golobby/container"
	"go-jwt-server/dal"
	"go-jwt-server/models"
	"go-jwt-server/types"
	"go-jwt-server/view_models"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const SALT = "SALT_HERE"
const SECRET = "SECRET_HERE"
const ExpiresInMinutes = 5

func Login(w http.ResponseWriter, r *http.Request) {
	var db *types.DatabaseT
	err := container.Make(&db)

	if err != nil {
		panic("Failed to resolve db context")
	}

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var loginViewModel view_models.LoginViewModel
	err = json.NewDecoder(r.Body).Decode(&loginViewModel)
	if err != nil {
		http.Error(w, "Failed decoding JSON body", 400)
		return
	}

	validate := validator.New()
	err = validate.Struct(loginViewModel)
	if err != nil {
		http.Error(w, "Register view model is not valid", 400)
		return
	}

	table := map[string]interface{}{
		"username": loginViewModel.Username,
		"password": HashPassword(loginViewModel.Password),
	}

	user, err := dal.QueryUser(db, &table)
	if err != nil {
		http.Error(w, "Failed querying user in database", 400)
		return
	} else {
		log.Printf("Successfully queried a user with is %d", user.Id)
	}

	tokenString, err := GenerateRandomString()
	if err != nil {
		http.Error(w, "Generation of token failed", 400)
		return
	}

	table = map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"rand":  tokenString,
		"id":    user.Id,
	}

	expiration := time.Now().Add(time.Minute * ExpiresInMinutes)
	token := CreateToken(&table, expiration)

	session := models.Session{
		Expires: expiration,
		UserId:  user.Id,
		Token:   token,
	}
	dal.AddSession(db, &session)

	_, err = w.Write([]byte(token))
	if err != nil {
		panic("Failed retuning result")
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var db *types.DatabaseT
	err := container.Make(&db)

	tokenString := ResolveBearerToken(r)
	session, err := dal.QuerySession(db, &map[string]interface{}{
		"token": tokenString,
	})

	if err != nil {
		http.Error(w, "Failed decoding find the token", 401)
		return
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("Success"), nil
	})

	tokenUserId := uint(token.Claims.(jwt.MapClaims)["id"].(float64))

	if token == nil || token.Claims == nil || session.UserId != tokenUserId {
		http.Error(w, "Unauthorized", 401)
		return
	}

	_, err = dal.DeleteSessions(db, session.UserId)

	if err != nil {
		http.Error(w, "Failed decoding JSON body", 400)
		return
	}

	_, err = w.Write([]byte("Successfully logged out"))
	if err != nil {
		panic("Failed retuning result")
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	var db *types.DatabaseT
	err := container.Make(&db)

	if err != nil {
		panic("Failed to resolve db context")
	}

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var registerViewModel view_models.RegisterViewModel
	err = json.NewDecoder(r.Body).Decode(&registerViewModel)
	if err != nil {
		http.Error(w, "Failed decoding JSON body", 400)
		return
	}

	validate := validator.New()
	err = validate.Struct(registerViewModel)
	if err != nil {
		http.Error(w, "Register view model is not valid", 400)
		return
	}

	user := models.User{
		Name:     registerViewModel.Name,
		Realm:    registerViewModel.Realm,
		Email:    registerViewModel.Email,
		Username: registerViewModel.Username,
		Password: HashPassword(registerViewModel.Password),
	}
	_, err = dal.AddUser(db, &user)
	if err != nil {
		http.Error(w, "Failed adding user to database", 400)
		return
	} else {
		log.Printf("Successfully added a user with is %d", user.Id)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)

	if err != nil {
		http.Error(w, "Failed encoding JSON result", 400)
		return
	}
}

func AuthorizeMiddleware(inner http.Handler) http.Handler {
	log.Print("AuthorizeMiddleware: called")
	mw := func(w http.ResponseWriter, r *http.Request) {
		var db *types.DatabaseT
		err := container.Make(&db)

		tokenString := ResolveBearerToken(r)
		session, err := dal.QuerySession(db, &map[string]interface{}{
			"token": tokenString,
		})

		if err != nil {
			http.Error(w, "Failed decoding find the token", 401)
			return
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("Success"), nil
		})

		tokenUserId := uint(token.Claims.(jwt.MapClaims)["id"].(float64))

		if token == nil || token.Claims == nil || session.UserId != tokenUserId {
			http.Error(w, "Unauthorized", 401)
			return
		}

		inner.ServeHTTP(w, r)
	}
	return http.HandlerFunc(mw)
}

func CreateToken(table *map[string]interface{}, expiration time.Time) string {
	props := jwt.MapClaims{
		"iss": "auth-app",
		"sub": "medium",
		"aud": "any",
		"exp": expiration.Unix(),
	}

	for key, value := range *table {
		props[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, props)
	jwtToken, err := token.SignedString([]byte(SECRET))

	if err != nil {
		panic("Failed creating a JWT token")
	}

	return jwtToken
}

func HashPassword(plaintextPassword string) string {
	h := sha256.New()
	h.Write([]byte(plaintextPassword + SALT))

	return b64.StdEncoding.EncodeToString(h.Sum(nil))
}

func GenerateRandomString() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ResolveBearerToken(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	return reqToken
}
