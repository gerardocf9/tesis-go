package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorInfoGeneral struct {
	Id              primitive.ObjectID `bson:"_id" json:"id"`
	IdMotor         uint               `bson:"IdMotor" json:"IdMotor"`
	Caracteristicas string             `bson:"caracteristicas,omitempty" json:"caracteristicas"`
	IdSensor        uint               `bson:"IdSensor" json:"IdSensor"`
	Data            DataSensor         `bson:"Data" json:"Data"`
	Time            time.Time          `bson:"Time" json:Time`
}

type DataSensor struct {
	AcelerationX float64 `bson:"AcelerationX" json:"AcelerationX"`
	AcelerationY float64 `bson:"AcelerationY" json:"AcelerationY"`
	AcelerationZ float64 `bson:"AcelerationZ" json:"AcelerationZ"`
}
