package model

import "github.com/jinzhu/gorm"

type Users struct {
	gorm.Model
	SecureId  string `json:"secure_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}
