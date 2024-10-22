package dto

import (
	httpresponse "crm-backend/helper/httpResponse"
	"crm-backend/models"
)

type CustomersResponse struct {
	httpresponse.Response
	Total     int               `json:"total"`
	Customers []models.Customer `json:"customers"`
}
type CustomerResponse struct {
	httpresponse.Response
	Customer models.Customer `json:"customer"`
}

type CustomerRequest struct {
	FirstName string `gorm:"size:100"`
	LastName  string `gorm:"size:100"`
	Company   string `gorm:"size:100"`
	Address   string `gorm:"size:255"`
	City      string `gorm:"size:100"`
	County    string `gorm:"size:100"`
	Postal    string `gorm:"size:20"`
	Phone     string `gorm:"size:20"`
	Email     string `gorm:"size:100"`
	Web       string `gorm:"size:100"`
}
