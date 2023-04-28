package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Telephone    string `json:"telephone"`
	ProfileImage string `json:"profile_image"`
	Address      string `json:"address"`
	City         string `json:"city"`
	Province     string `json:"province"`
	Country      string `json:"country"`
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
