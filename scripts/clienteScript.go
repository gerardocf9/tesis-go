package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gerardocf9/tesis-go/models"
	"github.com/posener/h2conn"
	"golang.org/x/net/http2"
)

func main() {

	argLength := len(os.Args[1:])

	if argLength < 3 {
		fmt.Println("Faltan argumentos")
		return
	}

	var (
		IdMotor         string
		s1              int
		s2              int
		s3              int
		s4              int
		caracteristicas string
		nDamage         int
		direccion       string
		umbInfVel       int
		umbSupVel       int
		umbInfAcel      int
		umbSupAcel      int
	)

	// flags declaration using flag package
	flag.StringVar(&direccion, "dir", "localhost:8080", "URL to connect.")
	flag.StringVar(&IdMotor, "i", "", "Specify IdMotor.")
	flag.IntVar(&nDamage, "dam", 1, "Specify nivel damage(0,...,10). Default 1.")
	flag.StringVar(&caracteristicas, "c", "Estado del motor", "Specify Caracteristicas. Default is Estado del motor")
	flag.IntVar(&s1, "s1", 0, "Specify IdSensor1 . Default 1.")
	flag.IntVar(&s2, "s2", 0, "Specify IdSensor1 . Default 2.")
	flag.IntVar(&s3, "s3", 0, "Specify IdSensor1 . Default 0.")
	flag.IntVar(&s4, "s4", 0, "Specify IdSensor1 . Default 0.")
	flag.IntVar(&umbInfVel, "uiv", 12, "Umbral inferior Velocidad . Default 12.")
	flag.IntVar(&umbSupVel, "usv", 20, "Umbral superior Velocidad . Default 20.")
	flag.IntVar(&umbInfAcel, "uia", 2, "Umbral inferior Aceleración . Default 2.")
	flag.IntVar(&umbSupAcel, "usa", 5, "Umbral superior Aceleración . Default 5.")

	flag.Parse() // after declaring flags we need to call it
	if (IdMotor == "") || (s1 == 0) || (s2 == 0) || (direccion == "") {

		fmt.Println("Faltan banderas")
		return
	}

	var ( //tipo|reservado|reservado|ubicacion|id
		ladLibre    uint64 = 0x0000000000000000
		ladCarga    uint64 = 0x0001000000000000
		chumYConex1 uint64 = 0x0002000000000000
		chumYConex2 uint64 = 0x0003000000000000
		idSensor    []uint64
		numHex      uint64
	)

	if s1 > 0 {
		numHex = ladLibre | uint64(s1)
		idSensor = append(idSensor, numHex)
	}

	if s2 > 0 {
		numHex = ladCarga | uint64(s2)
		idSensor = append(idSensor, numHex)
	}

	if s3 > 0 {
		numHex = chumYConex1 | uint64(s3)
		idSensor = append(idSensor, numHex)
	}

	if s4 > 0 {
		numHex = chumYConex2 | uint64(s4)
		idSensor = append(idSensor, numHex)
	}

	var post = models.SensorInfoGeneral{
		IdMotor:         IdMotor,
		IdSensor:        idSensor,
		Caracteristicas: caracteristicas,
		UmbInferiorVel:  int32(umbInfVel),
		UmbSuperiorVel:  int32(umbSupVel),
		UmbInferiorAcel: int32(umbInfAcel),
		UmbSuperiorAcel: int32(umbSupAcel),
	}

	ConnectServer(direccion, post, nDamage)
}

func ConnectServer(dir string, post models.SensorInfoGeneral, nivelD int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We use a client with custom http2.Transport since the server certificate is not signed by
	// an authorized CA, and this is the way to ignore certificate verification errors.
	d := &h2conn.Client{
		Client: &http.Client{
			Transport: &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		},
	}
	url := "https://" + dir + "/sensormessage"
	log.Println("0")
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
	err = in.Decode(&loginResp)
	if err != nil {
		log.Fatalf("Failed login: %v", err)
	}
	if loginResp != "OK" {
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

			if resp == "post" {
				sendPost(resp, 10, post, out, nivelD)
				log.Println("enviado post")
			}
			if resp == "exhaustiva" {
				sendExhaustive(resp, post, out, nivelD)
				log.Println("enviado exhaustiva")
			}

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
		dir := "https://tesis-fastapi-gf9gs.ondigitalocean.app/normal/" + dam + "/" + strconv.FormatInt(int64(sensor), 10)
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
	}
	post.Time = time.Now().Unix()
	// Send the message to the server
	err = out.Encode(post)
	if err != nil {
		log.Fatalf("Failed sending message: %v", err)
	}
}

func sendExhaustive(msg string, post models.SensorInfoGeneral, out *json.Encoder, nivelD int) {

	dam := int64(nivelD)
	err := out.Encode(msg)
	if err != nil {
		log.Fatalf("Failed sending message1: %v", err)
	}
	var postExhaustive = models.SensorExhaustive{
		IdMotor:         post.IdMotor,
		IdSensor:        post.IdSensor,
		Caracteristicas: post.Caracteristicas,
	}
	postExhaustive.Data = make([]models.DataExhaustive, 0, 5)

	for _, sensor := range postExhaustive.IdSensor {
		dam += int64(-2 + rand.Intn(4))
		if dam > 10 {
			dam = 10
		}
		if dam < 0 {
			dam = 0
		}

		dam2 := strconv.FormatInt(dam, 10)
		log.Println(dam2)
		dir := "https://tesis-fastapi-gf9gs.ondigitalocean.app/exhaustiva/" + dam2 + "/" + strconv.FormatInt(int64(sensor), 10)
		resp, err := http.Get(dir)
		if err != nil {
			log.Fatalln(err)
		}
		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var example models.DataExhaustive
		err = json.Unmarshal(body, &example)
		if err != nil {
			log.Fatal(err)
		}
		postExhaustive.Data = append(postExhaustive.Data, example)
	}
	postExhaustive.Time = time.Now().Unix()
	// Send the message to the server
	err = out.Encode(postExhaustive)
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
				log.Println(resp)
				message <- resp
				log.Println("Listen send")
				continue
			}
			message <- "error"
			return //No es ningun mensaje valido
		}
	}
}

func timer(message chan string, private chan int) {
	t := time.Now()
	u, _ := time.ParseDuration("5s")
	for {
		select {
		case <-private:
			log.Println("Canceled timer")
			return
		default:
			fmt.Printf(".")
			if time.Since(t) > u {
				message <- "post"
				t = time.Now()
				fmt.Printf("-")
			}
			time.Sleep(1 * time.Second)
		}
	}
}
