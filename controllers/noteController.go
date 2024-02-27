package controllers

import (
	"context"
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
		}

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

		result, err := noteCollection.Find(ctx, bson.D{{Key: "uid", Value: userId}})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error occuered while fetching notes"})
			return
		}

		var notes []models.Note

		if err = result.All(ctx, &notes); err != nil {
			log.Fatal(err)
		}

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, notes)

	}
}
