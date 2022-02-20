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
	Aceleracion  float64 `bson:"Aceleracion" json:"Aceleracion"`
	VelocidadX   float64 `bson:"VelocidadX" json:"VelocidadX"`
	VelocidadY   float64 `bson:"VelocidadY" json:"VelocidadY"`
	VelocidadZ   float64 `bson:"VelocidadZ" json:"VelocidadZ"`
}
