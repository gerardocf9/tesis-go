package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gerardocf9/tesis-go/bbdd"
	"github.com/gerardocf9/tesis-go/models"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {

	argLength := len(os.Args[1:])

	if argLength < 3 {
		fmt.Println("Faltan argumentos")
		return
	}

	var (
		IdMotor         string
		s1              string
		s2              string
		s3              string
		s4              string
		caracteristicas string
		nDamage         string
		reset           string
	)

	// flags declaration using flag package
	flag.StringVar(&IdMotor, "i", "", "Specify IdMotor.")
	flag.StringVar(&nDamage, "d", "1", "Specify nivel damage(0,...,10). Default 1.")
	flag.StringVar(&caracteristicas, "c", "Estado del motor", "Specify Caracteristicas. Default is Estado del motor")
	flag.StringVar(&s1, "s1", "1", "Specify IdSensor1 . Default 1.")
	flag.StringVar(&s2, "s2", "2", "Specify IdSensor1 . Default 2.")
	flag.StringVar(&s3, "s3", "0", "Specify IdSensor1 . Default 0.")
	flag.StringVar(&s4, "s4", "0", "Specify IdSensor1 . Default 0.")
	flag.StringVar(&reset, "r", "false", "Reset the BBDD . Default false.")

	flag.Parse() // after declaring flags we need to call it
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if reset == "true" {
		db := bbdd.MongoCN.Database("tesis")
		col := db.Collection("MotoresInDB")
		//reset list motors
		var arr = make([]string, 0, 10)
		_, err := col.UpdateOne(
			ctx,
			bson.D{},
			bson.D{
				{"$set", bson.D{{"IdMotor", arr}}},
			},
		)
		if err != nil {
			return
		}
		col = db.Collection("MotorData")
		_, err = col.DeleteMany(ctx, bson.M{})
		if err != nil {
			return
		}
	}
	var ( //tipo|reservado|reservado|ubicacion|id
		ladLibre    uint64 = 0x0000000000000000
		ladCarga    uint64 = 0x0001000000000000
		chumYConex1 uint64 = 0x0002000000000000
		chumYConex2 uint64 = 0x0003000000000000
		idSensor    []uint64
		numHex      uint64
	)

	num, _ := strconv.ParseInt(s1, 10, 64)
	if num > 0 {
		numHex = ladLibre | uint64(num)
		idSensor = append(idSensor, numHex)
	}

	num, _ = strconv.ParseInt(s2, 10, 64)
	if num > 0 {
		numHex = ladCarga | uint64(num)
		idSensor = append(idSensor, numHex)
	}

	num, _ = strconv.ParseInt(s3, 10, 64)
	if num > 0 {
		numHex = chumYConex1 | uint64(num)
		idSensor = append(idSensor, numHex)
	}

	num, _ = strconv.ParseInt(s4, 10, 64)
	if num > 0 {
		numHex = chumYConex2 | uint64(num)
		idSensor = append(idSensor, numHex)
	}

	t := time.Now()
	t = t.AddDate(0, -3, -4)
	var post = models.SensorInfoGeneral{
		IdMotor:         IdMotor,
		IdSensor:        idSensor,
		Caracteristicas: caracteristicas,
	}

	if bbdd.CheckConnection() == 0 {
		log.Fatal("Sin conexion a la BD")
	}
	//verificando si existe en la lista de motores con informacion
	err := bbdd.ListMotor(IdMotor)
	if err != nil {
		log.Printf("Failed listing motor %v", err)
	}

	for i := 0; i < 120; i++ {
		post.Data = make([]models.DataSensor, 0, 10)
		if i == 40 {
			aux, _ := strconv.ParseInt(nDamage, 10, 64)
			aux += 1
			if aux > 10 {
				aux = 10
			}
			nDamage = strconv.FormatInt(int64(aux), 10)
		}
		if i == 80 {
			aux, _ := strconv.ParseInt(nDamage, 10, 64)
			aux += 1
			if aux > 10 {
				aux = 10
			}
			nDamage = strconv.FormatInt(int64(aux), 10)
		}
		for _, sensor := range post.IdSensor {
			dir := "https://tesis-fastapi-gf9gs.ondigitalocean.app/normal/" + nDamage + "/" + strconv.FormatInt(int64(sensor), 10)
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

		t = t.AddDate(0, 0, 1)
		post.Time = t
		//inserting in bbdd
		_, status, err := bbdd.InsertRegistro(post)
		//error handling and ok to client
		if err != nil {
			log.Printf("Failed inserting in the bbdd: %v", err)
		}
		if !status {
			log.Printf("Cant insert Register in bbdd")
		}

		fmt.Printf("\ntoday date in format = %s", t.Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	}
	fmt.Printf("\nDamage level= %s", nDamage)
}
