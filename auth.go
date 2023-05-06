package gochijwt

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const secretKey = "vryscrtkey"

// defining a bunch of types
type CustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Auth struct {
	Claims  *CustomClaims `json:"claims"`
	User    *User         `json:"user"`
	Token   string        `json:"token"`
	Expires *time.Time    `json:"-"`
}

type AuthService struct {
	key string
}

// constructor
func NewAuthService() AuthService {
	return AuthService{
		key: secretKey,
	}
}

// methods
// create a new auth struct
func (as AuthService) CreateAuth(u *User) *Auth {

	exp := time.Now().Add(2 * time.Hour).Unix()

	claims := &CustomClaims{
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "sleazy_e",
			ExpiresAt: exp,
		},
	}

	return &Auth{
		Claims: claims,
		User:   u,
	}
}

// generate a jwt
// this will get called in the login route
func (as AuthService) GenerateToken(a *Auth) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.Claims)

	tokenStr, err := token.SignedString([]byte(as.key))

	if err != nil {
		return err
	}

	a.Token = tokenStr

	return nil
}

// validate the jwt
func (as AuthService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	tkn, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(as.key), nil
	})

	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("token not valid")
	}

	return tkn, nil
}

// this probably doesn't need to be part of the authservice, but w/e
func (as AuthService) MiddlewareJWT(next http.HandlerFunc) http.HandlerFunc {
	// see this demo https://hackernoon.com/creating-a-middleware-in-golang-for-jwt-based-authentication-cx3f32z8
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("eeauth")

		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				writeJSON(w, http.StatusBadRequest, "cookie not found")
			default:
				writeJSON(w, http.StatusInternalServerError, "something went wrong")
			}
			return
		}

		token, err := as.ValidateToken(cookie.Value)

		if err != nil {
			writeJSON(w, http.StatusUnauthorized, err)
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "eeclaims", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// helpers
func ValidateUser(u *User, us []*User) error {
	var ref int
	//check that a user exists
	for i, v := range us {
		if u.Username == v.Username {
			ref = i
			break
		}
		return errors.New("User doesn't exist")
	}
	//check that passwords match
	if u.Password != us[ref].Password {
		return errors.New("Passwords don't match")
	}

	return nil
}
