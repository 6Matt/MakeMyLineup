package handlers

import (
    "github.com/6Matt/se390-internal/libhttp"
    "github.com/gorilla/mux"
    "html/template"
    "net/http"
)

type SchedData struct {
    One string
    Two string
}

func GetSched(w http.ResponseWriter, r *http.Request) {
    d := SchedData{mux.Vars(r)["lastID"], mux.Vars(r)["festID"]}

    w.Header().Set("Content-Type", "text/html")

    tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/sched.html.tmpl")
    if err != nil {
        libhttp.HandleErrorJson(w, err)
        return
    }

    tmpl.Execute(w, d)
}