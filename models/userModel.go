package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct { 
    ID              primitive.ObjectID  `bson:"_id"`
    Fullname        *string             `json:"fullname" validate:"required,min=2,max=100"`
    Email           *string             `json:"email" validate:"email,required"`
    Username        *string             `json:"username" validate:"required,min=2,max=100"`
    Password        *string             `json:"Password" validate:"required,min=6""`
    Phonenumber     *string             `json:"phone_number" validate:"required"`
    Country         *string             `json:"country" validate:"required"`
    Invitation_code *string             `json:"invitation_code"`
    Token           *string             `json:"token"`
    Refresh_token   *string             `json:"refresh_token"`
    Created_at      time.Time           `json:"created_at"`
    Updated_at      time.Time           `json:"updated_at"`
    User_id         string              `json:"user_id"`
}
