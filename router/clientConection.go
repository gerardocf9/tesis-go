package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gerardocf9/tesis-go/models"
	"github.com/posener/h2conn"
)

type encoder interface {
	Encode(interface{}) error
}
type server struct {
	conns    map[uint]encoder
	idSensor map[uint]uint
	lock     sync.RWMutex
}

/*
type maxServer struct {
	sPost server
}
*/
var Servidor = server{
	conns:    make(map[uint]encoder),
	idSensor: make(map[uint]uint),
}

//var MaxServidor = maxServer{sPost: Servidor}

func getInfoSensor(w http.ResponseWriter, r *http.Request) {
	log.Println("sirve")

	log.Println(r.Body)
	//log.Println(c)
}
func (c server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h2conn.Accept(w, r)
	if err != nil {
		log.Printf("Failed creating http2 connection: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var (
		// in and out send and receive json messages to the client
		in, out = json.NewDecoder(conn), json.NewEncoder(conn)
		// Conn has a RemoteAddr property which helps us identify the client
		log = logger{remoteAddr: r.RemoteAddr}
	)

	// First check user login ids
	var idMotor, idSensor uint
	err = in.Decode(&idMotor)
	if err != nil {
		log.Printf("Failed reading login id: %v", err)
		return
	}
	err = in.Decode(&idSensor)
	if err != nil {
		log.Printf("Failed reading login id: %v", err)
		return
	}
	log.Printf("Trata de conectarse: %v", idMotor)

	err = c.login(idMotor, idSensor, out)
	if err != nil {
		err = out.Encode(err.Error())
		log.Printf(err.Error())
		if err != nil {
			log.Printf("Failed sending login response: %v", err)
		}
		return
	}
	//respuesta de idMotor aceptado
	err = out.Encode("OK")
	if err != nil {
		log.Printf("Failed sending login response: %v", err)
		return
	}
	// defer logout of user
	defer c.logout(idMotor)

	// Defer logout log message
	defer log.Printf("User logout: %s", idMotor)

	// wait for client to close connection
	post := models.SensorInfoGeneral{}
	for r.Context().Err() == nil {
		err = in.Decode(&post)
		if err != nil {
			log.Printf("Failed getting post: %v", err)
			return
		}
		log.Printf("Got msg: %+v", post)
		err = out.Encode("ok")
		if err != nil {
			log.Printf("failed sending response to client: %v", err)
			return
		}
	}
}

func (c *server) login(id uint, sensor uint, enc encoder) error {
	if _, ok := c.conns[id]; ok {
		return fmt.Errorf("user already exists")
	}
	c.conns[id] = enc
	c.idSensor[id] = sensor
	log.Println("login succesed")
	return nil
}

func (c *server) logout(id uint) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.conns, id)
}

type logger struct {
	remoteAddr string
}

func (l logger) Printf(format string, args ...interface{}) {
	log.Printf("[%s] %s", l.remoteAddr, fmt.Sprintf(format, args...))
}
