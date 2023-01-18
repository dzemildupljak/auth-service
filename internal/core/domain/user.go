package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id         int64  `json:"id"`
	Name       string `json:"name" gorm:"not null"`
	Username   string `json:"username" gorm:"not null"`
	Email      string `json:"email" validate:"required" gorm:"unique;not null;index"`
	Password   string `json:"-" validate:"required" gorm:"not null"`
	Address    string `json:"address"`
	Isverified bool   `json:"isverified" gorm:"default:false"`
	Tokenhash  []byte `json:"tokenhash" gorm:"not null"`
	// Role                 string       `json:"role"`
	// MailVerifyCode       string       `json:"mail_verify_code"`
	// MailVerifyExpire     sql.NullTime `json:"mail_verify_expire"`
	// PasswordVerifyCode   string       `json:"password_verify_code"`
	// PasswordVerifyExpire sql.NullTime `json:"password_verify_expire"`
}

// // SAME STRUCT WITH DIFFERENT TAGS USER!!!
// type TbUser struct {
// 	gorm.Model
// 	Name       string `gorm:"not null"`
// 	Username   string `gorm:"not null"`
// 	Email      string `gorm:"unique;not null;index"`
// 	Password   string `gorm:"not null"`
// 	Address    string
// 	Isverified bool `gorm:"default:false"`
// 	// Tokenhash            string       `json:"tokenhash"`
// 	// Role                 string       `json:"role"`
// 	// MailVerifyCode       string       `json:"mail_verify_code"`
// 	// MailVerifyExpire     sql.NullTime `json:"mail_verify_expire"`
// 	// PasswordVerifyCode   string       `json:"password_verify_code"`
// 	// PasswordVerifyExpire sql.NullTime `json:"password_verify_expire"`
// }

// func (TbUser) TableName() string {
// 	return "user"
// }

type CreateUserParams struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type SignupUserParams struct {
	Name      string `json:"name" validate:"nonzero"`
	Username  string `json:"username" validate:"nonzero"`
	Email     string `json:"email" validate:"nonzero"`
	Password  string `json:"password" validate:"nonzero"`
	RPassword string `json:"repeat_password" validate:"nonzero"`
	Address   string `json:"address"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}
