package bbdd

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoCN is the conection to bbdd
var MongoCN = ConectBD()

//cambiar esta variable para redireccionar a otra BBDD mongo
var clientOptions = options.Client().ApplyURI("mongodb+srv://gacf9:tesis123@cluster0.aohmy.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

//ConectBD Permite conectar a la bbdd
func ConectBD() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}
	log.Println("Conexion exitosa con la BD")
	return client
}

//CheckConnection pin to bbdd
func CheckConnection() int {
	err := MongoCN.Ping(context.TODO(), nil)
	if err != nil {
		return 0
	}
	return 1
}

func Close() {
	defer func() {
		if err := MongoCN.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
}
