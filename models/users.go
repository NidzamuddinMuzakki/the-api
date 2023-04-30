package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `json:"username" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	FirstName    string `json:"firstname" validate:"required"`
	LastName     string `json:"lastname" validate:"required"`
	Telephone    string `json:"telephone" validate:"required,regexp=(\\+62)[0-9]+$"`
	ProfileImage string `json:"profile_image" `
	Address      string `json:"address" validate:"required"`
	City         string `json:"city" validate:"required"`
	Province     string `json:"province" validate:"required"`
	Country      string `json:"country" validate:"required"`
}

type Profile struct {
	Username     string `json:"username"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Telephone    string `json:"telephone"`
	ProfileImage string `json:"profile_image"`
	Address      string `json:"address"`
	City         string `json:"city"`
	Province     string `json:"province"`
	Country      string `json:"country"`
}
