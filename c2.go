package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gerardocf9/tesis-go/bbdd"
	"github.com/gerardocf9/tesis-go/models"
	"github.com/posener/h2conn"
)

type decoder interface {
	Decode(interface{}) error
}
type encoder interface {
	Encode(interface{}) error
}
type server struct {
	//connOut     map[uint64]encoder //out,int
	//connIn      map[uint64]decoder
	connections map[uint64]string
	idSensor    map[uint64][]uint64
	lock        sync.RWMutex
}

/*
type maxServer struct {
	sPost server
}
*/
var Servidor = server{
	//connOut:    make(map[uint64]encoder),
	//connIn:     make(map[uint64]decoder),
	connections: make(map[uint64]string),
	idSensor:    make(map[uint64][]uint64),
}

var canalInterno = make(chan string)

var dataExhaustive = make(chan models.SensorInfoGeneral)

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
	)

	// First check user login ids
	var idMotor uint64
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

	err = c.login(idMotor, idSensor, r.RemoteAddr)
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
	defer log.Printf("User logout: %d", idMotor)

	// wait for client to close connection
	post := models.SensorInfoGeneral{}

	var cPost = make(chan int)
	//escuchar en el servidor
	go pedirInfo(cPost, out)
	var option string
	for r.Context().Err() == nil {

		select {
		case conEsperada := <-canalInterno:
			if r.RemoteAddr != conEsperada {
				continue
			}
			cPost <- 1
			log.Printf("Enviando data especial")
			out.Encode("VistaExhaustiva")

			err = in.Decode(&option)
			if err != nil {
				log.Printf("Failed getting vista exhaustiva: %v", err)
				return
			}
			//incoming data is not a post, reboot client
			if option != "VistaExhaustiva" {
				return
			}
			//Recibing post
			err = in.Decode(&post)
			if err != nil {
				log.Printf("Failed decoding VistaExhaustiva: %v", err)
				return
			}
			log.Println("enviado")
			dataExhaustive <- post
			go pedirInfo(cPost, out)

		default:
			//type of message
			err = in.Decode(&option)
			if err != nil {
				log.Printf("Failed getting post: %v", err)
				cPost <- 1
				return
			}

			//incoming data is not a post, reboot client
			if option != "post" {
				log.Printf("Disconecting that is not a post")
				cPost <- 1
				return
			}
			//Recibing post
			err = in.Decode(&post)
			if err != nil {
				log.Printf("Failed getting post: %v", err)
				cPost <- 1
				return
			}
			log.Printf("Got msg: %+v \n\n", post)

			//inserting in bbdd
			_, status, err := bbdd.InsertRegistro(post)
			//error handling and ok to client
			if err != nil {
				log.Printf("Failed inserting in the bbdd: %v", err)
			}
			if !status {
				log.Printf("Cant insert Register in bbdd")
			}

			//response
			err = out.Encode("ok")
			if err != nil {
				log.Printf("failed sending response to client: %v", err)
				cPost <- 1
				return
			}
			//fin del default
		} //fin del select
	} //fin del for
} //fin func

func (c *server) login(id uint64, sensor []uint64, conn string) error {
	if _, ok := c.connections[id]; ok {
		return fmt.Errorf("user already exists")
	}
	//	c.connOut[id] = enc
	//	c.connIn[id] = dec
	c.connections[id] = conn
	c.idSensor[id] = sensor
	log.Println("login succesed")
	return nil
}

func (c *server) logout(id uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.connections, id)
	delete(c.idSensor, id)
	//	delete(c.connOut, id)
	log.Println("loggedout")
}

//peticion continua de posts con func anonima
func pedirInfo(cPost chan int, out *json.Encoder) {
	for {
		select {
		case <-cPost:
			return
		default:

			log.Println("Sending message to get a post")
			err := out.Encode("post")
			if err != nil {
				log.Fatal("Error codificando post")
			}

			time.Sleep(15 * time.Second)
			//fin
		}
	}
}

/*
type logger struct {
	remoteAddr string
}

func (l logger) Printf(format string, args ...interface{}) {
	log.Printf("[%s] %s", l.remoteAddr, fmt.Sprintf(format, args...))
}
*/
