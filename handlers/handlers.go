package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/gerardocf9/tesis-go/middleware"
	"github.com/gerardocf9/tesis-go/router"
)

//Manejadores, setea el puerto y el handler escucha la direccion de PORT
func Manejadores() {
	mx := mux.NewRouter()
	//endpoints:
	mx.HandleFunc("/sensormessage", middleware.CheckBDH(router.Servidor)) //clientInterconexion
	mx.HandleFunc("/exhaustive", router.VistaExhaustiva)                  //API data exhaustiva
	mx.HandleFunc("/general", router.VistaGeneral)                        //despliego HTML basico

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}
	//permite request desde todos los sitios
	handler := cors.AllowAll().Handler(mx)

	log.Printf("Serving on " + PORT)
	srv := &http.Server{Addr: ":" + PORT, Handler: handler}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
