package domain

import (
	"database/sql"
)

type User struct {
	ID                  int64        `json:"id"`
	Name                string       `json:"name"`
	Username            string       `json:"username"`
	Email               string       `json:"email" validate:"required"`
	Password            string       `json:"password" validate:"required"`
	Address             string       `json:"address"`
	Tokenhash           string       `json:"tokenhash"`
	Isverified          bool         `json:"isverified"`
	OauthID             []string     `json:"oauth_id"`
	Role                string       `json:"role"`
	MailVerfyCode       string       `json:"mail_verfy_code"`
	MailVerfyExpire     sql.NullTime `json:"mail_verfy_expire"`
	PasswordVerfyCode   string       `json:"password_verfy_code"`
	PasswordVerfyExpire sql.NullTime `json:"password_verfy_expire"`
}

type CreateUserParams struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
