package handlers

import (
    "github.com/6Matt/se390-internal/libhttp"
    "html/template"
    "net/http"
)

func GetSched(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/sched.html.tmpl")
    if err != nil {
        libhttp.HandleErrorJson(w, err)
        return
    }

    tmpl.Execute(w, nil)
}