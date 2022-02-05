package router

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func VistaExhaustiva(w http.ResponseWriter, r *http.Request) {
	//obtenemos id
	aux := r.URL.Query().Get("idMotor")
	if len(aux) < 1 {
		http.Error(w, "Debe enviar el Parametro del ID", http.StatusBadRequest)
		return
	}

	sidMotor, err := strconv.ParseInt(aux, 10, 64)
	if err != nil {
		http.Error(w, "No se pudo obtener el ID", http.StatusBadRequest)
		return
	}
	idMotor := uint64(sidMotor)
	//Chequeamos exista el sensor a los motores conectado
	conn, ok := Servidor.connections[idMotor]
	if !ok {
		http.Error(w, "Cant connect with the motor", http.StatusBadRequest)
		return
	}
	canalInterno <- conn //stopping normal server, asking mem

	resp := <-dataExhaustive //receiving the data from the connection

	w.Header().Set("context-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
