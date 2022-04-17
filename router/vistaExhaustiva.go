package router

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func VistaExhaustiva(w http.ResponseWriter, r *http.Request) {
	//obtenemos id
	idMotor := r.URL.Query().Get("idMotor")
	if len(idMotor) < 1 {
		http.Error(w, "Debe enviar el Parametro del ID", http.StatusBadRequest)
		return
	}

	//Chequeamos exista el sensor a los motores conectado
	conn, ok := Servidor.connections[idMotor]
	if !ok {
		http.Error(w, "Cant connect with the motor", http.StatusBadRequest)
		return
	}
	fmt.Println("Solicitada Vista exhaustiva")
	canalInterno <- conn //stopping normal server, asking mem

	resp := <-dataExhaustive //receiving the data from the connection

	w.Header().Set("context-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
	fmt.Println("envio web de  Vista exhaustiva")
}
