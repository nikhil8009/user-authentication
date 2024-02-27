package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Note struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       *string            `json:"title"`
	Description *string            `json:"description"`
	Uid         string             `json:"uid"`
}
