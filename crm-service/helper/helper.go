package helper

import (
	httpresponse "crm-service/helper/httpResponse"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

// // RespondWithJSON sends a JSON response
// func RespondWithJSON(c *gin.Context, res interface{}, statusCode int) {
// 	c.JSON(statusCode, res)
// }

// // RespondWithError sends an error response
//
//	func RespondWithError(c *gin.Context, statusCode int, res httpresponse.Response) {
//		RespondWithJSON(c, res, statusCode)
//	}
func RespondWithJSON(c *gin.Context, res interface{}, statusCode int) {
	// response, _ := json.Marshal(res)
	c.Header("Content-Type", "application/json")
	c.JSON(statusCode, res)
}

// Respond Error with JSON
// Used for responding error messages in JSON content with w and code
func RespondWithError(c *gin.Context, code int, msg string, responseCode string) {
	res := httpresponse.PrepareResponse(responseCode, msg)
	fmt.Println(msg, code)
	RespondWithJSON(c, res, code)
}

func ValidateEmail(email string) bool {
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func ValidatePhoneNumberWithCountry(phone string) bool {
	const phoneRegex = `^[+]{1}(?:[0-9\-\(\)\/\.]\s?){6, 15}[0-9]{1}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}
