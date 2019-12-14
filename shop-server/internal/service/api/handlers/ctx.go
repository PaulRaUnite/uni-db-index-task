package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dgrijalva/jwt-go"

	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/config"
	"github.com/PaulRaUnite/uni-db-index-task/shop-server/internal/data"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
)

func GoodQ(r *http.Request) data.GoodQ {
	return Config(r).ClonedStorage().GoodQ()
}

func UserQ(r *http.Request) data.UserQ {
	return Config(r).ClonedStorage().UserQ()
}

func CountryQ(r *http.Request) data.CountryQ {
	return Config(r).ClonedStorage().CountryQ()
}

func InvoiceQ(r *http.Request) data.InvoiceQ {
	return Config(r).ClonedStorage().InvoiceQ()
}

func InvoicePartQ(r *http.Request) data.InvoicePartQ {
	return Config(r).ClonedStorage().InvoicePartQ()
}

func Config(r *http.Request) config.Config {
	return r.Context().Value("config").(config.Config)
}

func Log(r *http.Request) *logan.Entry {
	return ape.Log(r)
}

func ComplaintsQ(r *http.Request) *mongo.Collection {
	return Config(r).Complaints()
}

func SurveysQ(r *http.Request) *mongo.Collection {
	return Config(r).Surveys()
}

type AccountType int

const (
	User    AccountType = 0
	Reviwer AccountType = 1
	Admin   AccountType = 2
)

type Claims struct {
	UserID      int         `json:"user_id,required"`
	AccountType AccountType `json:"account_type"`
	jwt.StandardClaims
}

var mySigningKey = []byte("AllYourBase")

func IssueJWT(user data.User) (string, error) {
	claims := Claims{
		UserID:      user.ID,
		AccountType: AccountType(user.AccountType),
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			Issuer:   "shop-server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func UserIDFromClaims(r *http.Request) (int, error) {
	claims, err := GetClaims(r)
	if err != nil {
		return 0, err
	}
	if claims == nil {
		return 0, errors.New("no auth")
	}
	return claims.UserID, nil
}

func GetClaims(r *http.Request) (*Claims, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, nil
	}
	token = token[7:]
	var claims Claims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	return &claims, nil
}
