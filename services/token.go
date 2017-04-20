package services

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("blind runtime!")
	}
	here := path.Dir(filename)
	// get our private key for generating jwt token
	b, err := ioutil.ReadFile(here + "/../app.rsa")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		log.Fatal(err)
	}
	b, err = ioutil.ReadFile(here + "/../app.rsa.pub")
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		log.Fatal(err)
	}

}

type Token interface {
	Make(string, string) (string, error)
}

type token struct {
	UserService User
}

func NewToken(us User) Token {
	return &token{us}
}

func (s *token) Make(u string, p string) (string, error) {

	uid := s.UserService.Authenticate(u, p)
	log.Println("auth:", uid)
	if uid == 0 {
		return "", errors.New("user failed")
	}
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 100,
		Subject:   string(uid),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwtString, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.New("token failed")
	}

	return jwtString, nil
}

func (s *token) Validate(jwtString string) bool {

	parseToken := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected method %v", token.Header["alg"])
		}
		return publicKey, nil
	}

	token, err := jwt.ParseWithClaims(jwtString, &jwt.StandardClaims{}, parseToken)
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		log.Println("claims:", claims.Subject)
		return true
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("Malformed token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			fmt.Println("token time error")
		} else {
			fmt.Println("validation cannot handle token:", err)
		}
	} else {
		fmt.Println("cannot handle token:", err)
	}
	return false
}
