package bbdd

import (
	"context"
	"time"

	"github.com/gerardocf9/tesis-go/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckExist(idMotor string) (models.SensorInfoGeneral, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoCN.Database("tesis")
	col := db.Collection("MotorData")

	condicion := bson.M{"idMotor": idMotor}
	var resultado models.SensorInfoGeneral

	err := col.FindOne(ctx, condicion).Decode(&resultado)
	ID := resultado.IdMotor

	if err != nil {
		return resultado, false, ID
	}
	return resultado, true, ID
}
