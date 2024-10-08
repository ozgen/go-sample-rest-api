package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"go-sample-rest-api/config"
	"go-sample-rest-api/logging"
	"go-sample-rest-api/types"
	"go-sample-rest-api/utils"
	"net/http"
	"strconv"
	"time"
)

type contextKey string

const UserKey contextKey = "userID"

type jwtValidatorFunc func(string) (*jwt.Token, error)

var validateJWT jwtValidatorFunc = func(tokenString string) (*jwt.Token, error) {
	return validateJWTDefault(tokenString)
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log := logging.GetLogger()
		tokenString := utils.GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to validate token")

			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Error("Invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to convert userID to int")
			permissionDenied(w)
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to get user by ID")
			permissionDenied(w)
			return
		}

		// Add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func validateJWTDefault(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
