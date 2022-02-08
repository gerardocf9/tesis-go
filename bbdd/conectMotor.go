// Package ddbb provides ...
package bbdd

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateOne(filter, update interface{}) (result *mongo.UpdateResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("tesis")
	col := db.Collection("MotorData")
	//db.collection.find().limit(1).sort({$natural:-1})

	// A single document that match with the
	// filter will get updated.
	// update contains the filed which should get updated.
	//db.collection.find().limit(1).sort({$natural:-1})
	result, err = col.UpdateOne(ctx, filter, update)
	return
}
