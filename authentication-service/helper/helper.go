package helper

import (
	httpresponse "authentication-service/helper/httpResponse"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RespondWithJSON(w http.ResponseWriter, res interface{}, statusCode int) {
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// Respond Error with JSON
// Used for responding error messages in JSON content with w and code
func RespondWithError(w http.ResponseWriter, code int, msg string, responseCode string) {
	res := httpresponse.PrepareResponse(responseCode, msg)
	RespondWithJSON(w, res, code)
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

func GenerateBcryptHash(stringToHash string) (string, error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(strings.ToLower(stringToHash)), bcrypt.DefaultCost)
	return string(hashedByte), err
}
func CompareBcryptHash(stringToHash string, hashedString string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(strings.ToLower(stringToHash)))
	return err
}

func IsEmptyString(req string) bool {
	return strings.TrimSpace(req) == ""
}

func GenerateUuidV4() string {
	id := uuid.NewString()
	return id
}

func GenerateJwkRedisKey(aud interface{}) string {
	return fmt.Sprintf("user_jwk:%v", aud)
}
