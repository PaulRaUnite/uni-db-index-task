package handlers

import (
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

func CustomerQ(r *http.Request) data.CustomerQ {
	return Config(r).ClonedStorage().CustomerQ()
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

type Claims struct {
	UserID int `json:"user_id,required"`
	jwt.StandardClaims
}

var mySigningKey = []byte("AllYourBase")

func IssueJWT(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
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
	token := r.Header.Get("Authorization")[7:]
	var claims Claims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return mySigningKey, nil
	})
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
