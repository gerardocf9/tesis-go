# tesis-go
Repositorio backend tesis GO

Este sistema contiene un microservicio de sensorica, punto de entrada main.go
y un cliente de sensorica, punto de entrada cliente.go.
Ademas, contiene 2 scripts, LlenarBBDD.go y clienteScript.go.

# Estructura de carpetas

## main.go
Punto de entrada del servidor.

## cliente.go
Punto de entrada del cliente, interfaz grafica y llamado a la
función de interconexion en "client/interconexion.go"


## BBDD
Conectividad y acceso a la base de datos.

## client
Funcionamiento e interconexión del cliente con el servidor.

## handlers
Manejador de las rutas, endpoints, de la aplicación.
3 endpoints
	mx.HandleFunc("/sensormessage", middleware.CheckBDH(router.Servidor)) //Methods("POST")
servidor de comunicación continua a traves de un middleware.

	mx.HandleFunc("/exhaustive", router.VistaExhaustiva)
Información exhaustiva, solicitud de gran cantidad de muestras
para la FFT en el servicio WEB.

    mx.HandleFunc("/general", router.VistaGeneral)
Despliegue de un template (HTML), utilizado para comprobación grafica
del servicio una vez montado en DigitalOcean.

## middleware
middleware encargado de comprobar una conección correcta con el servicio
de BBDD (mongoDB Atlas)

## models
modelos y estructuras de datos definidos.

## router
implementación de los endpoints.

    clientConection.go
Servidor de sensorica.

    vistaExhaustiva.go
API para solicitar multiples mediciones del sensor y poder
generar una FFT.

    vistaGeneral.go
Despliegue de un template basico.

## scripts
scripts adicionales creados para automatizar el proceso.

## static/template
template desplegado en vistaGeneral.go

## vendor
Dependencias utilizadas en el proyecto, para evitar problemas
de versiones, funcionalidades descontinuadas o librerias eliminadas
y facilitar la compilación (solo se debe tener instalado Go).

## server.*
public y private key para conexión via HTTPS.
















