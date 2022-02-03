package middleware

import (
	"net/http"

	"github.com/gerardocf9/tesis-go/bbdd"
)

func CheckBDH(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if bbdd.CheckConnection() == 0 {
			http.Error(w, "Conexion perdida con la Base de Datos", 500)
			return
		}
		next.ServeHTTP(w, r)
	}
}

/*CheckBD is a middleware to know bbdd status*/
func CheckBD(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if bbdd.CheckConnection() == 0 {
			http.Error(w, "Conexion perdida con la Base de Datos", 500)
			return
		}
		next.ServeHTTP(w, r)
	}
}
