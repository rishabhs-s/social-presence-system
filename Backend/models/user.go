package models


import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
	Friends  Array  `gorm:"type:jsonb"`
	Parties  []Party `gorm:"many2many:user_parties"`
}

type OnlineUser struct {
	gorm.Model
	Name string `gorm:"column:username"` // Optional: customize column name
  }