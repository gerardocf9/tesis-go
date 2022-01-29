package main

import (
	"log"

	"github.com/gerardocf9/tesis-go/handlers"
)

func main() {
	log.Println("comenzo")
	/*	if bbdd.CheckConnection() == 0 {
			log.Fatal("Sin conexion a la BD")
		}
	*/handlers.Manejadores()
}
