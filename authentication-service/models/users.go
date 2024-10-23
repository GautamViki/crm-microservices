package models

import "time"

type User struct {
	EntityId   int    `gorm:"primaryKey;type:INT(10)" json:"entity_id"`
	FirstName  string `gorm:"type:varchar(55);not null" json:"firstName"`
	MiddleName string `gorm:"type:varchar(55);" json:"middleName"`
	LastName   string `gorm:"type:varchar(55);not null" json:"lastName"`
	Email      string `gorm:"type:varchar(55);unique;" json:"email"`
	Mobile     string `gorm:"type:varchar(55);unique;" json:"mobile"`
	Country    string `gorm:"type:varchar(55);not null" json:"country"`
	Password   string `gorm:"type:varchar(255);not null" json:"password"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
