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

//Manejadores, setea el puertoy el handler escucha la direccion de PORT
func Manejadores() {
	mx := mux.NewRouter()
	//mx.Handle("/sensormessage", router.Servidor)                          //clientInterconexion
	mx.HandleFunc("/sensormessage", middleware.CheckBDH(router.Servidor)) //Methods("POST")
	mx.HandleFunc("/general", router.VistaGeneral)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	handler := cors.AllowAll().Handler(mx)

	log.Printf("Serving on " + PORT)
	srv := &http.Server{Addr: ":" + PORT, Handler: handler}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
