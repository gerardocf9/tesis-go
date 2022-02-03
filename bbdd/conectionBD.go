package bbdd

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoCN is the conection to bbdd
var MongoCN = ConectBD()
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
	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := MongoCN.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
}
