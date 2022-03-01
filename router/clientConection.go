package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gerardocf9/tesis-go/bbdd"
	"github.com/gerardocf9/tesis-go/models"
	"github.com/posener/h2conn"
)

type encoder interface {
	Encode(interface{}) error
}
type server struct {
	conns       map[string]encoder
	connections map[string]string
	idSensor    map[string][]uint64
	lock        sync.RWMutex
}

/*
type maxServer struct {
	sPost server
}
*/
var Servidor = server{
	connections: make(map[string]string),
	conns:       make(map[string]encoder),
	idSensor:    make(map[string][]uint64),
}

var canalInterno = make(chan string)

var dataExhaustive = make(chan models.SensorExhaustive)

//var MaxServidor = maxServer{sPost: Servidor}

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
	var idMotor string
	var idSensor []uint64

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

	err = c.login(idMotor, idSensor, out, r.RemoteAddr)
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

	//verificando si existe en la lista de motores con informacion
	err = bbdd.ListMotor(idMotor)
	if err != nil {
		log.Printf("Failed listing motor %v", err)

	}

	// wait for client to close connection
	post := models.SensorInfoGeneral{}
	postExhaustive := models.SensorExhaustive{}

	var option string
	for r.Context().Err() == nil {
		select {
		case conEsperada := <-canalInterno:
			if r.RemoteAddr != conEsperada {
				continue
			}
			err = out.Encode("exhaustiva")
			if err != nil {
				log.Printf("failed sending response to client: %v", err)
				return
			}
			//type of message
			err = in.Decode(&option)
			if err != nil {
				log.Printf("Failed getting post: %v", err)
				return
			}

			//incoming data is not a post, reboot client
			if option != "exhaustiva" {
				return
			}

			//Recibing post
			err = in.Decode(&postExhaustive)
			if err != nil {
				log.Printf("Failed getting post: %v", err)
				return
			}

			log.Printf("enviado")
			dataExhaustive <- postExhaustive

		default:

			//type of message
			err = in.Decode(&option)
			if err != nil {
				log.Printf("Failed getting post: %v", err)
				return
			}

			//incoming data is not a post, reboot client
			if option != "post" {
				return
			}

			//Recibing post
			err = in.Decode(&post)
			if err != nil {
				log.Printf("Failed getting post: %v", err)
				return
			}

			//inserting in bbdd
			_, status, err := bbdd.InsertRegistro(post)
			//error handling and ok to client
			if err != nil {
				log.Printf("Failed inserting in the bbdd: %v", err)
			}
			if !status {
				log.Printf("Cant insert Register in bbdd")
			}

		}
	}
}

func (c *server) login(id string, sensor []uint64, enc encoder, s string) error {
	if _, ok := c.conns[id]; ok {
		return fmt.Errorf("user already exists")
	}
	c.conns[id] = enc
	c.connections[id] = s
	c.idSensor[id] = sensor
	log.Println("login succesed")
	return nil
}

func (c *server) logout(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.connections, id)
	delete(c.idSensor, id)
	delete(c.conns, id)
	log.Println("loggedout")
}

type logger struct {
	remoteAddr string
}

func (l logger) Printf(format string, args ...interface{}) {
	log.Printf("[%s] %s", l.remoteAddr, fmt.Sprintf(format, args...))
}
