// Package ddbb provides ...
package bbdd

import (
	"context"
	"time"

	"github.com/gerardocf9/tesis-go/models"
	"go.mongodb.org/mongo-driver/bson"
)

func ListMotor(idMotor uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("tesis")
	col := db.Collection("MotoresInDB")
	result, err := getList()

	if err != nil {
		return err
	}

	for _, idAct := range result.IdMotor {
		if idAct == idMotor {
			return nil
		}
	}

	result.IdMotor = append(result.IdMotor, idMotor)
	//db.collection.find().limit(1).sort({$natural:-1})

	// A single document that match with the
	// filter will get updated.
	// update contains the filed which should get updated.
	//db.collection.find().limit(1).sort({$natural:-1})
	_, err = col.UpdateOne(ctx, bson.D{}, result)
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
