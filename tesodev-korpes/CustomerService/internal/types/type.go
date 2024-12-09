package types

import (
	"time"
)

type CustomerRequestModel struct {
	FirstName string    `bson:"first_name" json:"first_name" validate:"required"`
	LastName  string    `bson:"last_name" json:"last_name" validate:"required"`
	Age       int       `bson:"age" json:"age" `
	Email     string    `bson:"email" json:"email"`
	Username  string    `bson:"username" json:"username" validate:"required"`
	Password  string    `bson:"password" json:"password" validate:"required"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type CustomerResponseModel struct {
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Username  string `bson:"username" json:"username"`
	Email     string `bson:"email" json:"email"`
}
type CustomerLoginResponseModel struct {
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Username  string `bson:"username" json:"username"`
	Email     string `bson:"email" json:"email"`
	Token     string `bson:"token" json:"token"`
	Id        string `bson:"_id" json:"id"`
	Password  string `bson:"password" json:"password"`
}

type CustomerUpdateModel struct {
	FirstName string    `bson:"first_name" json:"first_name"`
	LastName  string    `bson:"last_name" json:"last_name"`
	Email     string    `bson:"email" json:"email"`
	Username  string    `bson:"username" json:"username"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CustomerLogin struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
