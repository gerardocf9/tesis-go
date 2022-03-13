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
		s1              int
		s2              int
		s3              int
		s4              int
		caracteristicas string
		nDamage         int
		reset           string
		fInicio         int
		cantIteraciones int
	)

	// flags declaration using flag package
	flag.StringVar(&IdMotor, "i", "", "Specify IdMotor.")
	flag.IntVar(&nDamage, "d", 1, "Specify nivel damage(0,...,10). Default 1.")
	flag.StringVar(&caracteristicas, "c", "Estado del motor", "Specify Caracteristicas. Default is Estado del motor")
	flag.IntVar(&s1, "s1", 0, "Specify IdSensor1 . Default 1.")
	flag.IntVar(&s2, "s2", 0, "Specify IdSensor1 . Default 2.")
	flag.IntVar(&s3, "s3", 0, "Specify IdSensor1 . Default 0.")
	flag.IntVar(&s4, "s4", 0, "Specify IdSensor1 . Default 0.")
	flag.StringVar(&reset, "r", "false", "Reset the BBDD . Default false.")
	flag.IntVar(&fInicio, "fi", 90, "fecha Inicio cantidad de dias atras de hoy. Default 30")
	flag.IntVar(&cantIteraciones, "ci", 90, "Cantidad de iteraciones a ejecutarse.Default 30")

	flag.Parse() // after declaring flags we need to call it
	if (IdMotor == "") || (s1 == 0) || (s2 == 0) {

		fmt.Println("Faltan banderas")
		return
	}

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

	t := time.Now()
	t = t.AddDate(0, 0, -fInicio)

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
	if fInicio < cantIteraciones {
		cantIteraciones = fInicio
	}
	for i := 0; i < cantIteraciones; i++ {
		post.Data = make([]models.DataSensor, 0, 10)
		if i != 0 && (i%40) == 0 {
			nDamage += 1
			if nDamage > 10 {
				nDamage = 10
			}
			fmt.Printf("\nDamage level= %d\n", nDamage)
		}
		for _, sensor := range post.IdSensor {
			dir := "https://tesis-fastapi-gf9gs.ondigitalocean.app/normal/" + strconv.FormatInt(int64(nDamage), 10) + "/" + strconv.FormatInt(int64(sensor), 10)
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
		post.Time = t.Unix()
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
	fmt.Printf("\nDamage level= %d\n", nDamage)
}
