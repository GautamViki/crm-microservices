package middlerware

import (
	"authentication-service/config"
	h "authentication-service/helper"
	authhelper "authentication-service/helper/authHelper"
	"net/http"
	"strings"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := config.ConnectDB()
		logger := config.GetLoggerInstance()
		tokenString := r.Header.Get(h.Authorization)
		tokenSlice := strings.Split(tokenString, " ")
		if len(tokenSlice) != 2 || tokenSlice[0] != h.Bearer {
			logger.Log(h.TokenIsIncorrectFormatError, "", h.TokenIsIncorrectFormatErrorCode)
			h.RespondWithError(w, http.StatusUnauthorized, h.TokenIsIncorrectFormatError, h.TokenIsIncorrectFormatErrorCode)
			return
		}
		if err := authhelper.ValidateToken(tokenSlice[1], db); err != nil {
			logger.Log(h.TokenInvalidError, err.Error(), h.TokenInvalidErrorCode)
			h.RespondWithError(w, http.StatusUnauthorized, h.TokenInvalidError, h.TokenInvalidErrorCode)
			return
		}
		next.ServeHTTP(w, r)
	})
}
