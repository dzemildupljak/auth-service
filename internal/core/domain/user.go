package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// Id uuid.UUID `json:"id" gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()" bson:"-"`
	Id         uuid.UUID `json:"id" gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()" bson:"_id,omitempty"`
	Name       string    `json:"name" gorm:"not null" bson:"name"`
	Username   string    `json:"username" gorm:"not null" bson:"username,omitempty"`
	Email      string    `json:"email" validate:"required" gorm:"unique;not null; index" bson:"email,omitempty"`
	Password   string    `json:"-" validate:"required" gorm:"not null" bson:"password,omitempty"`
	Address    string    `json:"address" bson:"address"`
	Isverified bool      `json:"isverified" gorm:"default:false" bson:"isverified"`
	Tokenhash  []byte    `json:"tokenhash" gorm:"not null" bson:"tokenhash"`
	CreatedAt  time.Time `json:"-" gorm:"not null" bson:"createdat"`
	UpdatedAt  time.Time `json:"-" gorm:"not null" bson:"updatedat"`
	DeletedAt  time.Time `json:"-" gorm:"index; not null" bson:"deletedat"`
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
