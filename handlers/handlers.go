package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"encoding/json"

	"strings"

	"github.com/gerardocf9/tesis-go/router"
	"github.com/posener/h2conn"
)

//Manejadores, setea el puertoy el handler escucha la direccion de PORT
func Manejadores() {
	mx := mux.NewRouter()
	mx.HandleFunc("/echo", ServeHTTP)
	mx.HandleFunc("/general", router.VistaGeneral)
	//static resources
	mx.Handle("/assets", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	handler := cors.AllowAll().Handler(mx)

	log.Printf("Serving on " + PORT)
	srv := &http.Server{Addr: ":" + PORT, Handler: handler}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

type handler struct{}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h2conn.Accept(w, r)
	if err != nil {
		log.Printf("Failed creating connection from %s: %s", r.RemoteAddr, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	//tipo de conexion
	log.Printf("Got connection: %s", r.Proto)
	// Conn has a RemoteAddr property which helps us identify the client
	log := logger{remoteAddr: r.RemoteAddr}

	log.Printf("Joined")
	defer log.Printf("Left")

	// Create a json encoder and decoder to send json messages over the connection
	var (
		in  = json.NewDecoder(conn)
		out = json.NewEncoder(conn)
	)

	// Loop forever until the client hangs the connection, in which there will be an error
	// in the decode or encode stages.
	for {
		var msg string
		err = in.Decode(&msg)
		if err != nil {
			log.Printf("Failed decoding request: %v", err)
			return
		}

		log.Printf("Got: %q", msg)

		msg = strings.ToUpper(msg)

		err = out.Encode(msg)
		if err != nil {
			log.Printf("Failed encoding response: %v", err)
			return
		}

		log.Printf("Sent: %q", msg)
	}
}

type logger struct {
	remoteAddr string
}

func (l logger) Printf(format string, args ...interface{}) {
	log.Printf("[%s] %s", l.remoteAddr, fmt.Sprintf(format, args...))
}
