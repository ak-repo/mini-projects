package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email         string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash  string `gorm:"size:255;not null" json:"-"`
	Username      string `gorm:"size:255;not null"`
	Role          string `gorm:"size:100;not null"`
	Status        string `gorm:"size:50;default:active"`
	EmailVerified bool
}

type Address struct {
	gorm.Model
	Phone       string `gorm:"size:20;not null"`
	UserID      uint   `gorm:"not null;index"`
	AddressLine string `gorm:"type:text;not null"`
	City        string `gorm:"size:100;not null"`
	State       string `gorm:"size:100;not null"`
	PostalCode  string `gorm:"size:20;not null"`
	Country     string `gorm:"size:100;not null"`
}
