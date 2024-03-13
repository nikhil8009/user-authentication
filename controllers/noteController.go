package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"user-auth/database"
	"user-auth/models"

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

		var notes []models.Note

		if err := c.BindJSON(&notes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			cancel()
			return
		}

		insertArray := make([]interface{}, len(notes))

		for i := range notes {

			notes[i].ID = primitive.NewObjectID()
			str, ok := userId.(string)
			if ok {
				notes[i].Uid = str
			}

			notes[i].Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			notes[i].Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

			insertArray[i] = notes[i]
		}

		_, err := noteCollection.InsertMany(ctx, insertArray)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error while creating note"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Notes inserted successfully"})

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

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var notes []models.Note

		if err := c.BindJSON(&notes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			defer cancel()
			return
		}

		var wm []mongo.WriteModel
		for i := 0; i < len(notes); i++ {
			jsonData := bson.M{"title": notes[i].Title, "description": notes[i].Description}

			wm = append(wm, mongo.NewUpdateOneModel().SetFilter(bson.M{"timestamp": notes[i].TimeStamp}).SetUpdate(bson.M{"$set": jsonData}))
		}

		_, err := noteCollection.BulkWrite(ctx, wm)

		defer cancel()
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

	}
}

type List struct {
	Timestamps []string `json:"timestamps" binding:"required"`
}

func DeleteNote() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		data := new(List)

		if error := c.BindJSON(&data); error != nil {
			fmt.Println("error occured")
		}

		result, err := noteCollection.DeleteMany(ctx, bson.M{"timestamp": bson.M{"$in": data.Timestamps}})
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
