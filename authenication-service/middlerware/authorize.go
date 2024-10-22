package middlerware

import (
	"authentication-service/config"
	h "authentication-service/helper"
	authhelper "authentication-service/helper/authHelper"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := config.ConnectDB()
		redis := config.ConnectRedis()
		logger := config.GetLoggerInstance()
		tokenString := r.Header.Get(h.Authorization)
		tokenSlice := strings.Split(tokenString, " ")
		if len(tokenSlice) != 2 || tokenSlice[0] != h.Bearer {
			logger.Log(h.TokenInvalidError, "", h.TokenInvalidErrorCode)
			h.RespondWithError(w, http.StatusUnauthorized, h.TokenInvalidError, h.TokenInvalidErrorCode)
			return
		}

		token, _, err := new(jwt.Parser).ParseUnverified(tokenSlice[1], jwt.MapClaims{})
		if err != nil {
			logger.Log("Error while decoding token:", err.Error(), h.GeneralDecline)
			h.RespondWithError(w, http.StatusUnauthorized, h.UnableToProcessError, h.GeneralDecline)
			return
		}
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Log("Unable to fetch token claims", "", h.GeneralDecline)
			h.RespondWithError(w, http.StatusUnauthorized, h.UnableToProcessError, h.GeneralDecline)
			return
		}
		jwkKey := h.GenerateJwkRedisKey(mapClaims["aud"])
		publicKey, _ := redis.Get(context.Background(), jwkKey).Result()
		// Validate jwt token
		if err := authhelper.ValidateToken(tokenSlice[1], publicKey, db); err != nil {
			logger.Log(h.TokenInvalidError, err.Error(), h.TokenInvalidErrorCode)
			h.RespondWithError(w, http.StatusUnauthorized, h.TokenInvalidError, h.TokenInvalidErrorCode)
			return
		}
		next.ServeHTTP(w, r)
	})
}
