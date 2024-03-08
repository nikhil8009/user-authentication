package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"user-athentication/database"
	"user-athentication/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var noteCollection *mongo.Collection = database.OpenCollection(database.Client, "notes")

func AddNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("uid")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var note models.Note

		note.ID = primitive.NewObjectID()
		str, ok := userId.(string)
		if ok {
			note.Uid = str
		}

		if err := c.BindJSON(&note); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			cancel()
			return
		}

		note.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		note.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		insertionNumber, err := noteCollection.InsertOne(ctx, note)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error while creating note"})
			return
		}

		c.JSON(http.StatusOK, insertionNumber)

	}
}

func GetNotes() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, _ := c.Get("uid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		match := bson.D{{Key: "$match", Value: bson.D{{Key: "uid", Value: userId}}}}
		sort := bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}}

		result, err := noteCollection.Aggregate(ctx, mongo.Pipeline{match, sort})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		var notes []models.Note = make([]models.Note, 0)

		if err = result.All(ctx, &notes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		}

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": notes})

	}
}

func UpdateNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		noteId := c.Param("id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var note *models.Note

		if err := c.BindJSON(&note); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			defer cancel()
			return
		}

		objectId, error := primitive.ObjectIDFromHex(noteId)
		if error != nil {
			log.Println(error.Error())
			defer cancel()
			return
		}

		jsonData := bson.M{}
		if note.Title != nil {
			jsonData["title"] = note.Title
		}
		if note.Description != nil {
			jsonData["description"] = note.Description
		}

		err := noteCollection.FindOneAndUpdate(ctx, bson.M{"_id": objectId}, bson.M{"$set": jsonData})
		defer cancel()
		fmt.Println(err.Err())
		if err.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Err()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Note updated successfully"})
	}
}

type List struct {
	Ids []primitive.ObjectID `bson:"ids" binding:"required"`
}

func DeleteNote() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		data := new(List)

		if error := c.BindJSON(&data); error != nil {
			fmt.Println("error occured")
		}
		fmt.Println(data)
		result, err := noteCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": data.Ids}})
		defer cancel()

		if result.DeletedCount <= 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "No records found to delete"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Note deleted successfully"})
	}

}
