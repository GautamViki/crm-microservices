package handler

import (
	"authentication-service/config"
	h "authentication-service/helper"
	authhelper "authentication-service/helper/authHelper"
	httpresponse "authentication-service/helper/httpResponse"
	"authentication-service/internals/dto"
	"authentication-service/models"
	"authentication-service/repositery"
	"authentication-service/repositery/repo"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type userHandler struct {
	userRepo repositery.UserRepo
}

func NewUserHandler() *userHandler {
	return &userHandler{
		userRepo: repo.NewUserRepo(),
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.UserRequest true "User data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /users/register [post]
func (u *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRequest
	db := config.ConnectDB()
	logger := config.GetLoggerInstance()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log(h.BodyDecodeError, err.Error(), h.BodyDecodeErrorCode)
		h.RespondWithError(w, http.StatusBadRequest, h.BodyDecodeError, h.BodyDecodeErrorCode)
		return
	}
	user, err := u.userRepo.CreateUser(req, db)
	if err != nil {
		logger.Log(h.UserCreateError, err.Error(), h.UserCreateErrorCode)
		h.RespondWithError(w, http.StatusInternalServerError, h.UnableToProcessError, h.UserCreateErrorCode)
		return
	}
	logger.Log(h.UserCreateSuccess, "", h.APISuccessCode)
	userResp := createUserResponse(user, h.UserCreateSuccess)
	h.RespondWithJSON(w, userResp, http.StatusCreated)
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve all users from the database
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.UsersResponse
// @Failure 500 {object} httpresponse.Response
// @Router /users [get]
func (u *userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	logger := config.GetLoggerInstance()
	users, err := u.userRepo.GetUsers(db)
	if err != nil {
		logger.Log(h.UsersFetchError, err.Error(), h.UsersFetchErrorCode)
		h.RespondWithError(w, http.StatusInternalServerError, h.UnableToProcessError, h.UsersFetchErrorCode)
		return
	}
	logger.Log(h.UsersFetchSuccess, "", h.UsersFetchSuccessCode)
	usersResp := createUsersResponse(users, h.UsersFetchSuccess)
	h.RespondWithJSON(w, usersResp, http.StatusOK)
}

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Retrieve a user by email from the database
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param email path string true "Email of the user"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /users/{email} [get]
func (u *userHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	logger := config.GetLoggerInstance()
	email := chi.URLParam(r, "email")
	identifier := map[string]string{
		h.Email: email,
	}
	user, err := u.userRepo.GetUserByUserIdentifier(identifier, db)
	if err != nil {
		logger.Log(h.UserFetchByEmailError, err.Error(), h.UserFetchByEmailErrorCode)
		h.RespondWithError(w, http.StatusInternalServerError, h.UnableToProcessError, h.UserFetchByEmailErrorCode)
		return
	}
	logger.Log(h.UserFetchSuccess, "", h.APISuccessCode)
	userResp := createUserResponse(user, h.UserFetchSuccess)
	h.RespondWithJSON(w, userResp, http.StatusCreated)
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Retrieve a user by ID from the database
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "ID of the user"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /users/{id} [get]
func (u *userHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	logger := config.GetLoggerInstance()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log(h.UserIdInvalidError, err.Error(), h.UserIdInvalidErrorCode)
		h.RespondWithError(w, http.StatusBadRequest, h.UserIdInvalidError, h.UserIdInvalidErrorCode)
		return
	}
	user, err := u.userRepo.GetUserById(id, db)
	if err != nil {
		logger.Log(h.UserFetchByIdError, err.Error(), h.UserFetchByIdErrorCode)
		h.RespondWithError(w, http.StatusInternalServerError, h.UnableToProcessError, h.UserFetchByIdErrorCode)
		return
	}
	logger.Log(h.UserFetchSuccess, "", h.APISuccessCode)
	userResp := createUserResponse(user, h.UserFetchSuccess)
	h.RespondWithJSON(w, userResp, http.StatusCreated)
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return a token
// @Tags Users
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequest true "Login credentials"
// @Success 202 {object} dto.TokenResponse
// @Failure 400 {object} httpresponse.Response
// @Failure 500 {object} httpresponse.Response
// @Router /login [post]
func (u *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	redis := config.ConnectRedis()
	logger := config.GetLoggerInstance()
	loginRequest := dto.LoginRequest{}
	// Decode the incoming JSON request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logger.Log(h.BodyDecodeError, err.Error(), h.BodyDecodeErrorCode)
		h.RespondWithError(w, http.StatusBadRequest, h.BodyDecodeError, h.BodyDecodeErrorCode)
		return
	}
	if h.IsEmptyString(loginRequest.Email) && h.IsEmptyString(loginRequest.Mobile) {
		logger.Log(h.EmailAndMobileShouldNotEmptyError, "", h.EmailAndMobileShouldNotEmptyErrorCode)
		h.RespondWithError(w, http.StatusBadRequest, h.EmailAndMobileShouldNotEmptyError, h.EmailAndMobileShouldNotEmptyError)
		return
	}

	identifier := map[string]string{
		h.Email:  loginRequest.Email,
		h.Mobile: loginRequest.Mobile,
	}
	user, err := u.userRepo.GetUserByUserIdentifier(identifier, db)
	if err != nil {
		logger.Log(h.LoginRequestFailedError, err.Error(), h.LoginRequestFailedErrorCode)
		h.RespondWithError(w, http.StatusInternalServerError, h.UnableToProcessError, h.LoginRequestFailedErrorCode)
		return
	}
	if user.EntityId == 0 {
		logger.Log(h.InvalidUserCredientialError, "", h.InvalidUserCredientialErrorCode)
		h.RespondWithError(w, http.StatusNotFound, h.InvalidUserCredientialError, h.InvalidUserCredientialErrorCode)
		return
	}
	if err := h.CompareBcryptHash(loginRequest.Password, user.Password); err != nil {
		logger.Log(h.LoginPasswordIncorrectError, err.Error(), h.LoginPasswordIncorrectErrorCode)
		h.RespondWithError(w, http.StatusBadRequest, h.LoginPasswordIncorrectError, h.LoginPasswordIncorrectErrorCode)
		return
	}

	tokenResp, err := authhelper.GenerateToken(user, redis)
	if err != nil {
		logger.Log(h.TokenGenerationError, err.Error(), h.TokenGenerationErrorCode)
		h.RespondWithError(w, http.StatusBadRequest, h.TokenGenerationError, h.TokenGenerationErrorCode)
		return
	}
	h.RespondWithJSON(w, tokenResp, http.StatusAccepted)
}

// Authentication godoc
// @Summary Validate user token
// @Description Validate the user's token
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} httpresponse.Response
// @Failure 401 {object} httpresponse.Response
// @Router /authorize [get]
func (u *userHandler) Authentication(w http.ResponseWriter, r *http.Request) {
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
	logger.Log(h.AuthorizationSuccess, "", h.APISuccessCode)
	res := httpresponse.PrepareResponse(h.APISuccessCode, h.AuthorizationSuccess)
	h.RespondWithJSON(w, res, http.StatusOK)
}

func createUserDto(user models.User) dto.User {
	return dto.User{
		EntityId:   user.EntityId,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Email:      user.Email,
		Mobile:     user.MiddleName,
		Country:    user.Country,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func createUsersResponse(users []models.User, message string) dto.UsersResponse {
	userDto := []dto.User{}
	for _, user := range users {
		userDto = append(userDto, createUserDto(user))
	}
	return dto.UsersResponse{
		Response: httpresponse.PrepareResponse(h.APISuccessCode, message),
		Total:    len(users),
		Users:    userDto,
	}
}

func createUserResponse(user models.User, message string) dto.UserResponse {
	return dto.UserResponse{
		Response: httpresponse.PrepareResponse(h.APISuccessCode, message),
		User:     createUserDto(user),
	}
}
