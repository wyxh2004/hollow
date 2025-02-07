package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string             `bson:"email" json:"email" binding:"required,email"`
	Password  string             `bson:"password" json:"password,omitempty" binding:"required,min=6"`
	Avatar    string             `bson:"avatar" json:"avatar"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserResponse struct {
	ID     primitive.ObjectID `json:"id"`
	Email  string             `json:"email"`
	Avatar string             `json:"avatar"`
}
