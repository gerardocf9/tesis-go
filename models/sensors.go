package models

import "time"

type Post struct {
	Motor string
	Data  DataSensor
	Time  time.Time
}


type DataSensor struct {
	AcelerationX  float64
	AcelerationY  float64
	AcelerationZ float64
}

type Motor struct {
	Id        string
	Caract    string
	NumSensor int
}
