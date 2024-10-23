package middlerware

import (
	"crm-service/config"
	h "crm-service/helper"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		baseUrl := os.Getenv(h.AuthBaseUrl)
		logger := config.GetLoggerInstance()
		tokenString := c.GetHeader(h.Authorization)
		tokenSlice := strings.Split(tokenString, " ")
		if len(tokenSlice) != 2 || tokenSlice[0] != h.Bearer {
			logger.Log(h.TokenHeaderEmptyError, "", h.TokenHeaderEmptyErrorCode)
			h.RespondWithError(c, http.StatusUnauthorized, h.TokenHeaderEmptyError, h.TokenHeaderEmptyErrorCode)
			c.Abort()
			return
		}

		// Send the token to the authentication service for validation
		req, err := http.NewRequest("GET", baseUrl+"/authorize", nil)
		if err != nil {
			logger.Log(h.TokenInvalidError, err.Error(), h.TokenInvalidErrorCode)
			h.RespondWithError(c, http.StatusUnauthorized, h.TokenInvalidError, h.TokenInvalidErrorCode)
			c.Abort()
			return
		}
		req.Header.Set(h.Authorization, tokenString)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Log(h.TokenInvalidError, err.Error(), h.TokenInvalidErrorCode)
			h.RespondWithError(c, http.StatusUnauthorized, h.TokenInvalidError, h.TokenInvalidErrorCode)
			c.Abort()
			return
		}
		if resp.StatusCode != http.StatusOK {
			logger.Log(h.TokenInvalidError, h.TokenIsExpiredError, h.TokenIsExpiredErrorCode)
			h.RespondWithError(c, http.StatusUnauthorized, h.TokenInvalidError, h.TokenIsExpiredErrorCode)
			c.Abort()
			return
		}
		c.Next()
	}
}
