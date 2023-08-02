package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type SignUpResponse struct {
	Id            primitive.ObjectID `bson:"_id"`
	Fullname      string             `json:"fullname"`
	Username      string             `json:"username"`
	Email         string             `json:"email"`
	Token         string             `json:"token"`
	Refresh_token string             `json:"refresh_token"`
}

type LoginResponse struct {
	Id            primitive.ObjectID `bson:"_id"`
	Email         string             `json:"email"`
	Token         string             `json:"token"`
	Refresh_token string             `json:"refresh_token"`
}