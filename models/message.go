package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Box struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	OwnerID     primitive.ObjectID `bson:"owner_id" json:"owner_id"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Message struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	BoxID       primitive.ObjectID   `bson:"box_id" json:"box_id"`
	SenderID    primitive.ObjectID   `bson:"sender_id" json:"sender_id,omitempty"`
	Content     string               `bson:"content" json:"content"`
	IsAnonymous bool                 `bson:"is_anonymous" json:"is_anonymous"`
	LikeCount   int                  `bson:"like_count" json:"like_count"`
	LikedBy     []primitive.ObjectID `bson:"liked_by" json:"liked_by,omitempty"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
}

type MessageResponse struct {
	ID          primitive.ObjectID `json:"id"`
	BoxID       primitive.ObjectID `json:"box_id"`
	Content     string             `json:"content"`
	SenderEmail string             `json:"sender_email,omitempty"`
	IsAnonymous bool               `json:"is_anonymous"`
	LikeCount   int                `json:"like_count"`
	IsLiked     bool               `json:"is_liked"`
	CreatedAt   time.Time          `json:"created_at"`
}
