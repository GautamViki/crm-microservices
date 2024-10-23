package helper

const (
	BodyDecodeErrorCode                   = "1000"
	BodyDecodeError                       = "Error occurred in decodig request body."
	UnableToProcessError                  = "Unable to process request."
	APISuccessCode                        = "0"
	EmailAndMobileShouldNotEmptyErrorCode = "1001"
	EmailAndMobileShouldNotEmptyError     = "Please provide email or mobile number."
	InvalidUserCredientialErrorCode       = "1002"
	InvalidUserCredientialError           = "Email or mobile number is invalid."
	LoginRequestFailedErrorCode           = "1003"
	LoginRequestFailedError               = "Error occurred while login."
	LoginPasswordIncorrectErrorCode       = "1004"
	LoginPasswordIncorrectError           = "Login password is invalid."
	TokenGenerationErrorCode              = "1005"
	TokenGenerationError                  = "Failed to generate token."
	TokenInvalidErrorCode                 = "1006"
	TokenInvalidError                     = "Invalid token."
	TokenIsExpiredErrorCode               = "1007"
	TokenIsExpiredError                   = "Token is expired."
	GeneralDecline                        = "1008"
	AuthorizationSuccess                  = "User authorization successfully."
	TokenIsIncorrectFormatErrorCode       = "1009"
	TokenIsIncorrectFormatError           = "Token formate is incorrect."
)

const (
	UserCreateErrorCode       = "2000"
	UserCreateSuccessCode     = "2001"
	UsersFetchErrorCode       = "2002"
	UsersFetchSuccessCode     = "2003"
	UserFetchByEmailErrorCode = "2004"
	UserIdInvalidErrorCode    = "2005"
	UserFetchByIdErrorCode    = "2006"
	UserCreateError           = "Error occurred while creating user."
	UserCreateSuccess         = "User creation successfully."
	UsersFetchError           = "Error occurred while fetching user list."
	UsersFetchSuccess         = "User list fetched successfully."
	UserFetchByEmailError     = "Error occurred while fetching user by email."
	UserFetchSuccess          = "User fetched successfully."
	UserIdInvalidError        = "Invalid user id."
	UserFetchByIdError        = "Error occurred while fetching user by id."
)

const (
	Email         = "email"
	Mobile        = "mobile"
	Bearer        = "Bearer"
	Authorization = "Authorization"
)
