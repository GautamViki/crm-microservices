package repo

import (
	h "authentication-service/helper"
	"authentication-service/internals/dto"
	"authentication-service/models"
	"time"

	"gorm.io/gorm"
)

type mysqlUserRepo struct{}

func NewUserRepo() *mysqlUserRepo {
	return &mysqlUserRepo{}
}

func (u *mysqlUserRepo) CreateUser(req dto.UserRequest, db *gorm.DB) (models.User, error) {
	password, err := h.GenerateBcryptHash(req.Password)
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		Email:      req.Email,
		Mobile:     req.Mobile,
		Country:    req.Country,
		Password:   password,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *mysqlUserRepo) GetUsers(db *gorm.DB) ([]models.User, error) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *mysqlUserRepo) GetUserByUserIdentifier(identifier map[string]string, db *gorm.DB) (models.User, error) {
	var user models.User
	if email, ok := identifier[h.Email]; ok && !h.IsEmptyString(email) {
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			return models.User{}, err
		}
	} else if mobile, ok := identifier[h.Mobile]; ok && !h.IsEmptyString(mobile) {
		if err := db.Where("email = ?", mobile).First(&user).Error; err != nil {
			return models.User{}, err
		}
	}
	return user, nil
}

func (u *mysqlUserRepo) GetUserById(entityId int, db *gorm.DB) (models.User, error) {
	var user models.User
	if err := db.Where("entity_id = ?", entityId).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *mysqlUserRepo) UpdateUser(entityId int, post dto.UserRequest, db *gorm.DB) (models.User, error) {
	user := models.User{
		FirstName: post.FirstName,
		LastName:  post.LastName,
		Mobile:    post.Mobile,
		Email:     post.Email,
		Country:   post.Country,
		UpdatedAt: time.Now(),
	}
	if err := db.Model(&models.User{}).Where("entity_id = ?", entityId).Updates(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *mysqlUserRepo) DeleteUser(entityId int, db *gorm.DB) (models.User, error) {
	var user models.User
	if err := db.Delete(&user).Where("entity_id = ?", entityId).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
