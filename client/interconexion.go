package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/gerardocf9/tesis-go/models"
	"github.com/posener/h2conn"
	"golang.org/x/net/http2"
)

const url = "https://localhost:8080/sensormessage"

func ConnectServer(ch chan int, post models.SensorInfoGeneral, logp binding.String, nivelD int) {
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

	//chanels for internal comunication
	var message = make(chan string)
	var private = make(chan int)
	go timer(message, private)
	go listen(message, in, private)
	// Loop until user terminates
	for {

		select {
		case resp := <-message:
			// Receive the response from the server
			if resp == "error" {
				log.Fatalf("Llego un error")
			}
			if resp == "ok" {
				logp.Set("Got response " + resp)
			}
			if resp == "post" {
				sendPost(resp, 10, post, out, nivelD)
				logp.Set("sending post ")
			}
			if resp == "exhaustiva" {
				sendPost(resp, 1000, post, out, nivelD)
				logp.Set("sending exhaustive")
			}

		case <-ch:
			log.Println("Discconecting")
			private <- 1
			err = out.Encode("disconnect")
			if err != nil {
				log.Fatalf("Failed sending message: %v", err)
			}
			log.Println("disconnect")
			ch <- 1
			log.Println("Returning")
			return
		}
	}

}

func sendPost(msg string, nData int, post models.SensorInfoGeneral, out *json.Encoder, nivelD int) {

	dam := strconv.FormatInt(int64(nivelD), 10)

	err := out.Encode(msg)
	if err != nil {
		log.Fatalf("Failed sending message1: %v", err)
	}
	post.Data = make([]models.DataSensor, 0, nData)

	for _, sensor := range post.IdSensor {
		dir := "https://tesis-fastapi-gf9gs.ondigitalocean.app/" + dam + "/" + strconv.FormatInt(int64(sensor), 10)
		resp, err := http.Get(dir)
		if err != nil {
			log.Fatalln(err)
		}
		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var examples []models.DataSensor
		err = json.Unmarshal(body, &examples)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range examples {
			post.Data = append(post.Data, v)
		}
		log.Printf(" %+v", post)
	}
	post.Time = time.Now()
	// Send the message to the server
	err = out.Encode(post)
	if err != nil {
		log.Fatalf("Failed sending message: %v", err)
	}
}

func listen(message chan string, in *json.Decoder, private chan int) {
	for {
		select {
		case <-private:
			log.Println("Listen canceled")
			return
		default:
			var resp string
			err := in.Decode(&resp)
			if err != nil {
				log.Println("Failed receiving message: %v", err)
			}
			if (resp == "post") || (resp == "ok") || (resp == "exhaustiva") {
				message <- resp
				continue
			}
			message <- "error"
			return //No es ningun mensaje valido
		}
	}
}

func timer(message chan string, private chan int) {
	for {
		select {
		case <-private:
			log.Println("Canceled timer")
			return
		default:
			message <- "post"
			time.Sleep(15 * time.Second)
		}
	}
}
