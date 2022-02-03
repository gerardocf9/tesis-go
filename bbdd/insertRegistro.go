package bbdd

import (
	"context"
	"log"
	"time"

	"github.com/gerardocf9/tesis-go/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertRegistro(u models.SensorInfoGeneral) (*mongo.InsertOneResult, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//slecting BBDD and collection
	db := MongoCN.Database("tesis")
	col := db.Collection("MotorData")

	//	u.Password, _ = EncriptarPassword(u.Password)

	result, err := col.InsertOne(ctx, u)
	if err != nil {
		return result, false, err
	}

	objID := result.InsertedID
	log.Println("El id de la impresion es ")
	log.Println(objID)
	return result, true, nil
}

func InsertManyRegistro(u []interface{}) (*mongo.InsertManyResult, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//slecting BBDD and collection
	db := MongoCN.Database("tesis")
	col := db.Collection("MotorData")

	result, err := col.InsertMany(ctx, u)
	if err != nil {
		return result, false, err
	}
	return result, true, nil
	/*
	   fmt.Println("Result of InsertMany")

	   // print the insertion ids of the multiple
	   // documents, if they are inserted.
	   for id := range insertManyResult.InsertedIDs {
	       fmt.Println(id)
	   }
	*/
}
