package dto

import (
	httpresponse "authentication-service/helper/httpResponse"
	"time"
)

type UserRequest struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Country    string `json:"country"`
	Password   string `json:"password"`
}

type User struct {
	EntityId   int       `json:"entityId"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Email      string    `json:"email"`
	Mobile     string    `json:"mobile"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type UsersResponse struct {
	httpresponse.Response
	Total int    `json:"total"`
	Users []User `json:"users"`
}
type UserResponse struct {
	httpresponse.Response
	User User `json:"user"`
}
