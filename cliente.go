package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Motor")

	//motor
	id_log := binding.NewString()
	id_log.Set("2451")
	id := widget.NewEntryWithData(id_log)

	pot_log := binding.NewString()
	pot_log.Set("1")
	pot := widget.NewEntryWithData(pot_log)

	subNivel3 := container.New(layout.NewFormLayout(), widget.NewLabel("ID: "), id)
	subNivel4 := container.New(layout.NewFormLayout(), widget.NewLabel("Potencia(HP): "), pot)

	nivel2 := container.New(layout.NewGridLayout(2), subNivel3, subNivel4)

	dmg := binding.NewFloat()
	dmg.Set(0)
	damageLevel := widget.NewSliderWithData(0, 10.0, dmg)

	subSubNivel := container.New(layout.NewGridWrapLayout(fyne.NewSize(110, 30)), widget.NewLabel("Damage level: "), widget.NewLabelWithData(binding.FloatToStringWithFormat(dmg, "%0.0f%%")))
	nivel3 := container.New(layout.NewFormLayout(), subSubNivel, damageLevel)

	info_log := binding.NewString()
	info_log.Set("Datos del motor")
	info := widget.NewEntryWithData(info_log)
	nivel4 := container.New(layout.NewFormLayout(), widget.NewLabel("Informacion: "), info)
	//sensores
	s1 := binding.NewString()
	s1.Set("1")
	sen1 := widget.NewEntryWithData(s1)
	subS1 := container.New(layout.NewFormLayout(), widget.NewLabel("s1: "), sen1)
	subNivels1 := container.New(layout.NewGridLayout(2), subS1, widget.NewLabel("Lado Libre"))

	s2 := binding.NewString()
	s2.Set("2")
	sen2 := widget.NewEntryWithData(s2)
	subS2 := container.New(layout.NewFormLayout(), widget.NewLabel("s2: "), sen2)
	subNivels2 := container.New(layout.NewGridLayout(2), subS2, widget.NewLabel("Lado Carga"))

	s3 := binding.NewString()
	s3.Set("3")
	sen3 := widget.NewEntryWithData(s3)
	subS3 := container.New(layout.NewFormLayout(), widget.NewLabel("s3: "), sen3)
	subNivels3 := container.New(layout.NewGridLayout(2), subS3, widget.NewLabel("chumacera-accesorio"))

	s4 := binding.NewString()
	s4.Set("4")
	sen4 := widget.NewEntryWithData(s4)
	subS4 := container.New(layout.NewFormLayout(), widget.NewLabel("s4: "), sen4)
	subNivels4 := container.New(layout.NewGridLayout(2), subS4, widget.NewLabel("chumacera-accesorio"))

	s5 := binding.NewString()
	s5.Set("5")
	sen5 := widget.NewEntryWithData(s5)
	subS5 := container.New(layout.NewFormLayout(), widget.NewLabel("s5: "), sen5)
	subNivels5 := container.New(layout.NewGridLayout(2), subS5, widget.NewLabel("chumacera-accesorio"))

	//log:
	logM := binding.NewString()
	logM.Set("Hi!")

	//coneccion
	ip_log := binding.NewString()
	ip_log.Set("localhost:8080")

	ip := widget.NewEntryWithData(ip_log)
	conect := widget.NewButton("conectar", func() {
		log.Println("funcion")
		conectServidor(ip_log, id_log, pot_log, info_log, s1, s2, s3, s4, s5, logM, dmg)
	})
	subNivel1 := container.New(layout.NewFormLayout(), widget.NewLabel("IP: "), ip)
	subNivel2 := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), conect, layout.NewSpacer())

	nivel1 := container.New(layout.NewGridLayout(2), subNivel1, subNivel2)

	w.SetContent(
		fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
			widget.NewLabel("Conección: "),
			nivel1,
			widget.NewSeparator(),
			widget.NewLabel("Motor: "),
			nivel2,
			nivel3,
			nivel4,
			widget.NewSeparator(),
			widget.NewLabel("Sensores: "),
			subNivels1,
			subNivels2,
			subNivels3,
			subNivels4,
			subNivels5,
			widget.NewSeparator(),
			widget.NewLabel("Log: "),
			widget.NewLabelWithData(logM),
			//widget.NewEntryWithData(str),
		))

	w.Resize(fyne.Size{Height: 320, Width: 480})
	w.ShowAndRun()

}

func conectServidor(ip, id, pot, info binding.String, s1, s2, s3, s4, s5, logp binding.String, dmg binding.Float) {
	dir, err := ip.Get()
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}
	//idmotor a unsigned
	aux, err := id.Get()
	num, _ := strconv.ParseInt(aux, 10, 64)
	idMotor := uint64(num)
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}
	//potencia a float
	aux, err = pot.Get()
	num, _ = strconv.ParseInt(aux, 10, 64)
	potencia := float64(num)
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}

	informacion, err := info.Get()
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}

	var idSensor []uint64

	aux, err = s1.Get()
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}
	num, _ = strconv.ParseInt(aux, 10, 64)
	idSensor = append(idSensor, uint64(num))
	aux, err = s2.Get()
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}
	num, _ = strconv.ParseInt(aux, 10, 64)
	idSensor = append(idSensor, uint64(num))
	aux, err = s3.Get()
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}
	num, _ = strconv.ParseInt(aux, 10, 64)
	idSensor = append(idSensor, uint64(num))
	aux, err = s4.Get()
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}
	num, _ = strconv.ParseInt(aux, 10, 64)
	idSensor = append(idSensor, uint64(num))
	aux, err = s5.Get()
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}
	num, _ = strconv.ParseInt(aux, 10, 64)
	idSensor = append(idSensor, uint64(num))

	nivelD, err := dmg.Get()
	if err != nil {
		log.Fatalf("No se pudo obtener la info")
	}

	log.Println(dir + " " + informacion)
	log.Println(idMotor)
	log.Println(potencia)
	log.Println(nivelD)
	for _, sen := range idSensor {
		log.Println(sen)
	}
}
