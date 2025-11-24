package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID         string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Address    string `json:"address"`
	Email      string `gorm:"unique" json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password,omitempty"`
	IsActive   bool   `json:"is_active"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  int64  `json:"created_at"`
}
