package models

import (
	"time"
)

type SensorExhaustive struct {
	//	Id              primitive.ObjectID `bson:"_id" json:"id"`
	IdMotor         string           `bson:"IdMotor" json:"IdMotor"`
	Caracteristicas string           `bson:"caracteristicas,omitempty" json:"caracteristicas"`
	IdSensor        []uint64         `bson:"IdSensor" json:"IdSensor"`
	Data            []DataExhaustive `bson:"Data" json:"Data"`
	Time            time.Time        `bson:"Time" json:"Time"`
}
type DataExhaustive struct {
	IdSensorData uint64    `bson:"IdSensorData" json:"IdSensorData"`
	Data         []float64 `bson:"Data" json:"Data"`
}
