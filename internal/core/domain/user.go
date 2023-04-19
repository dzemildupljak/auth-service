package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// Id uuid.UUID `json:"id" gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()" bson:"-"`
	Id         uuid.UUID `json:"id" validate:"nonzero" gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()" bson:"_id,omitempty"`
	Name       string    `json:"name" gorm:"not null" bson:"name"`
	Username   string    `json:"username" validate:"nonzero" gorm:"not null" bson:"username,omitempty"`
	Email      string    `json:"email" validate:"nonzero" gorm:"unique;not null; index" bson:"email,omitempty"`
	Password   string    `json:"-" validate:"nonzero" gorm:"not null" bson:"password,omitempty"`
	Address    string    `json:"address" bson:"address"`
	Isverified bool      `json:"isverified" gorm:"default:false" bson:"isverified"`
	Tokenhash  []byte    `json:"-" validate:"nonzero" gorm:"not null" bson:"tokenhash"`
	CreatedAt  time.Time `json:"-" gorm:"not null" bson:"createdat"`
	UpdatedAt  time.Time `json:"-" gorm:"not null" bson:"updatedat"`
	DeletedAt  time.Time `json:"-" gorm:"index; not null" bson:"deletedat"`
	Role       string    `json:"role" validate:"nonzero" gorm:"not null" bson:"role,omitempty"`
	// MailVerifyCode       string       `json:"mail_verify_code"`
	// MailVerifyExpire     sql.NullTime `json:"mail_verify_expire"`
	// PasswordVerifyCode   string       `json:"password_verify_code"`
	// PasswordVerifyExpire sql.NullTime `json:"password_verify_expire"`
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

type UserMiddleware struct {
	Id         string `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Isverified bool   `json:"isverified"`
}
