package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/gerardocf9/tesis-go/models"
	"github.com/posener/h2conn"
	"golang.org/x/net/http2"
)

const url = "https://localhost:8080/sensormessage"

func ConnectServer(ch chan int, post models.SensorInfoGeneral, logp binding.String) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	defer logp.Set("Exit")

	//login request
	err = out.Encode(post.IdMotor)
	if err != nil {
		log.Fatalf("Failed Encoding idMotor: %v", err)
	}
	//send 2 messages
	err = out.Encode(post.IdSensor)
	if err != nil {
		log.Fatalf("Failed Encoding idSensor: %v", err)
	}
	//chequeando
	log.Println("Check idmotor")
	logp.Set("Check idmotor")
	err = in.Decode(&loginResp)
	if err != nil {
		log.Fatalf("Failed login: %v", err)
	}
	if loginResp == "OK" {
		logp.Set("login succeded")
	} else {
		log.Fatalf("Failed login 2: %v", err)
	}
	// Loop until user terminates
	for {

		select {
		case <-ch:
			log.Println("Discconecting")
			err = out.Encode("disconnect")
			if err != nil {
				log.Fatalf("Failed sending message: %v", err)
			}
			log.Println("disconnected from server")
			ch <- 1
			return

		default:
			// Receive the response from the server
			var resp string
			err = in.Decode(&resp)
			if err != nil {
				log.Fatalf("Failed receiving message: %v", err)
			}
			if resp == "ok" {
				logp.Set("Got response " + resp)
				continue
			}
			if resp == "post" {
				sendPost("post", 10, false, post, out)
				log.Println("Sended post")
				continue
			}
			if resp == "VistaExhaustiva" {
				sendPost("VistaExhaustiva", 1000, false, post, out)
				log.Println("Sended exhaustive")
			}
		}
	}
}

func sendPost(msg string, nData int, repetir bool, post models.SensorInfoGeneral, out *json.Encoder) {
	var (
		x float64
		y float64
		z float64
	)
	min := float64(3)
	max := float64(10)

	err := out.Encode(msg)
	if err != nil {
		log.Fatalf("Failed sending message1: %v", err)
	}
	post.Data = make([]models.DataSensor, 0, nData)
	for i := 0; i < nData; i++ {
		//numeros aleatorios para la data, me da un numero entre 0 y 1
		x = min + rand.Float64()*(max-min)
		y = min + rand.Float64()*(max-min)
		z = min + rand.Float64()*(max-min)

		post.Data = append(post.Data, models.DataSensor{
			AcelerationX: x,
			AcelerationY: y,
			AcelerationZ: z,
		})
	}
	post.Time = time.Now()
	// Send the message to the server
	err = out.Encode(post)
	if err != nil {
		log.Fatalf("Failed sending message: %v", err)
	}
}
