package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IchwanDwiNursid/go_restfullapi/config"
	"github.com/IchwanDwiNursid/go_restfullapi/types"
	"github.com/IchwanDwiNursid/go_restfullapi/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userId"

func CreateJwt(secret []byte, useId int) (string,error) {

	expiration := time.Second * time.Duration(config.Envs.JwtExpirationInSeconds)
	// This is Like Payload 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"userId" : strconv.Itoa(useId),
		"expiredAt" : time.Now().Add(expiration).Unix(),
	})

	tokenString , err := token.SignedString(secret)
	if err != nil {
		return "",err
	}

	return tokenString , nil 
}

func WithJWTAuth(handlerFunc http.HandlerFunc , store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter,r *http.Request){
		// get the token from user request
		tokenString := getTokenFromRequest(r)
		//validate JWT
		token,err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permisionDenied(w)
			return 
		}

		if !token.Valid{
			log.Printf("failed to validate token: %v", err)
			permisionDenied(w)
			return 
		}
		// if is we need to fetch the UserId from DB (id from token)
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)
		
		userId,err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userId to int: %v", err)
			permisionDenied(w)
			return
		}
		
		u,err := store.GetUserById(userId)
		
		if err != nil {
			log.Printf("failed to get user by id ")
			permisionDenied(w)
			return
		}
		
		// set contex userid 
		ctx := r.Context()
		ctx = context.WithValue(ctx ,UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w,r)
	}
}


func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != ""{
		return tokenAuth
	}

	return ""
}

func validateToken(t string)(*jwt.Token,error){
	return jwt.Parse(t,func(t *jwt.Token) (interface{} , error){
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil , fmt.Errorf("unexpected signin method: %v",t.Header["alg"])
		}
	

	return []byte(config.Envs.JWTSecret),nil
})}

func permisionDenied(w http.ResponseWriter){
	utils.WriteError(w,http.StatusForbidden,fmt.Errorf("permission denied"))
}

func GetUserIdFromContext(ctx context.Context) int {
	userId , ok := ctx.Value(UserKey).(int)

	if !ok {
		return -1
	}

	return userId

}