package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       *string            `json:"title"`
	Description *string            `json:"description"`
	Uid         string             `json:"uid"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	TimeStamp   string             `json:"time_stamp"`
}
