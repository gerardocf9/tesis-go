package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gerardocf9/tesis-go/models"
	"github.com/posener/h2conn"
	"golang.org/x/net/http2"
)

const url = "https://localhost:8080/sensormessage"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go catchSignal(cancel)

	// We use a client with custom http2.Transport since the server certificate is not signed by
	// an authorized CA, and this is the way to ignore certificate verification errors.
	d := &h2conn.Client{
		Client: &http.Client{
			Transport: &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		},
	}

	conn, resp, err := d.Connect(ctx, url)
	if err != nil {
		log.Fatalf("Initiate conn: %s", err)
	}
	defer conn.Close()

	// Check server status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad status code: %d", resp.StatusCode)
	}

	var (
		// in and out send and receive json messages to the server
		in  = json.NewDecoder(conn)
		out = json.NewEncoder(conn)
		//loginRespuesta
		loginResp string
	)
	defer log.Println("Exited")

	var sidMotor, sidSensor int
	//login values
	log.Println("ingrese el id del motor: ")
	_, err = fmt.Scanf("%d", &sidMotor)
	if err != nil {
		log.Fatalf("invalid input: %v", err)
	}
	log.Println("ingrese el id del sensor: ")
	_, err = fmt.Scanf("%d", &sidSensor)

	if err != nil {
		log.Fatalf("invalid input: %v", err)
	}
	//convert to uint64
	idMotor := uint64(sidMotor)
	idSensor := []uint64{uint64(sidSensor)}
	//login request
	err = out.Encode(idMotor)
	if err != nil {
		log.Fatalf("Failed Encoding idMotor: %v", err)
	}
	//send 2 messages
	err = out.Encode(idSensor)
	if err != nil {
		log.Fatalf("Failed Encoding idSensor: %v", err)
	}
	//chequeando
	log.Println("Check idmotor")
	err = in.Decode(&loginResp)
	if err != nil {
		log.Fatalf("Failed login: %v", err)
	}
	if loginResp == "OK" {
		log.Println("login succeded")
	} else {
		log.Fatalf("Failed login 2: %v", err)
	}

	post := models.SensorInfoGeneral{
		IdMotor:         idMotor,
		IdSensor:        idSensor,
		Caracteristicas: "potencia ... ubicacion en planta ...",
	}

	var (
		x float64
		y float64
		z float64
	)
	min := float64(3)
	max := float64(10)

	go func() {
		for {
			//numeros aleatorios para la data, me da un numero entre 0 y 1
			x = min + rand.Float64()*(max-min)
			y = min + rand.Float64()*(max-min)
			z = min + rand.Float64()*(max-min)

			post.Time = time.Now()
			post.Data = models.DataSensor{
				AcelerationX: x,
				AcelerationY: y,
				AcelerationZ: z,
			}
			// Send the message to the server
			err = out.Encode(post)
			if err != nil {
				log.Fatalf("Failed sending message: %v", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	// Loop until user terminates
	for ctx.Err() == nil {

		// Receive the response from the server
		var resp string
		err = in.Decode(&resp)
		if err != nil {
			log.Fatalf("Failed receiving message: %v", err)
		}
		fmt.Printf("Got response %q\n", resp)
	}
}

func catchSignal(cancel context.CancelFunc) {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("Cancelling due to interrupt")
	cancel()
}
