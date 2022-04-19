// Package ddbb provides ...
package bbdd

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gerardocf9/tesis-go/models"
	"go.mongodb.org/mongo-driver/bson"
)

//ListMotor devuelve una lista de motores para saber de cuales se tiene información en la BBDD
func ListMotor(idMotor string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("tesis")
	col := db.Collection("MotoresInDB")

	if CheckConnection() == 0 {
		return errors.New("Can't connect with the DB")
	}

	result, err := getList()

	if err != nil {
		return err
	}

	for _, idAct := range result.IdMotor {
		if idAct == idMotor {
			log.Println("ya existe en la BBDD")
			return nil
		}
	}

	result.IdMotor = append(result.IdMotor, idMotor)
	// ultimo elemento: db.collection.find().limit(1).sort({$natural:-1})

	// A single document that match with the
	// filter will get updated.
	// update contains the filed which should get updated.

	_, err = col.UpdateOne(
		ctx,
		bson.D{},
		bson.D{
			{"$set", bson.D{{"IdMotor", result.IdMotor}}},
		},
	)

	log.Println("agregado a la BBDD")
	return err
}

func getList() (models.MotorConnected, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("tesis")
	col := db.Collection("MotoresInDB")
	var result models.MotorConnected

	err := col.FindOne(ctx, bson.D{}).Decode(&result)
	return result, err
}
