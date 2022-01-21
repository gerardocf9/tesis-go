package router
/*

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/posener/h2conn"
)

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h2conn.Accept(w, r)
	if err != nil {
		log.Printf("Failed creating connection from %s: %s", r.RemoteAddr, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

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
	go func() {
		for {
			err = out.Encode("Mensaje asincrono desde la go routine...")
			if err != nil {
				log.Printf("Failed encoding response: %v", err)
				return
			}

			log.Printf("Sent: GOROUTINE message ")

			time.Sleep(10 * time.Second)
		}
	}()
}
*/
