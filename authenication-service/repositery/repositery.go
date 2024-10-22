package repositery

import (
	"authentication-service/internals/dto"
	"authentication-service/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(req dto.UserRequest, db *gorm.DB) (models.User, error)
	GetUsers(db *gorm.DB) ([]models.User, error)
	GetUserByUserIdentifier(identifier map[string]string, db *gorm.DB) (models.User, error)
	GetUserById(entityId int, db *gorm.DB) (models.User, error)
	UpdateUser(entityId int, req dto.UserRequest, db *gorm.DB) (models.User, error)
	DeleteUser(entityId int, db *gorm.DB) (models.User, error)
}
