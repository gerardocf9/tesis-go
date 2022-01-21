package router

import(
    "html/template"
    "net/http"
    "path"
)

type Book struct{

    Title string
    Author string
}

func VistaGeneral( w http.ResponseWriter, r *http.Request){
    book := Book{"Building web app with go","Jeremy"}
    fp := path.Join("static","template/index.html")

    tmpl,err := template.ParseFiles(fp)
    if err != nil {

        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
    }
    if err := tmpl.Execute(w,book); err!=nil{
        http.Error(w,err.Error(),http.StatusInternalServerError)
    }


}
