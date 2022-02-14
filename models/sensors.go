package models

import (
	"time"
)

type SensorInfoGeneral struct {
	//	Id              primitive.ObjectID `bson:"_id" json:"id"`
	IdMotor         string       `bson:"IdMotor" json:"IdMotor"`
	Caracteristicas string       `bson:"caracteristicas,omitempty" json:"caracteristicas"`
	IdSensor        []uint64     `bson:"IdSensor" json:"IdSensor"`
	Data            []DataSensor `bson:"Data" json:"Data"`
	Time            time.Time    `bson:"Time" json:"Time"`
}

type DataSensor struct {
	IdSensorData uint64  `bson:"IdSensorData" json:"IdSensorData"`
	AcelerationX float64 `bson:"AcelerationX" json:"AcelerationX"`
	AcelerationY float64 `bson:"AcelerationY" json:"AcelerationY"`
	AcelerationZ float64 `bson:"AcelerationZ" json:"AcelerationZ"`
}
