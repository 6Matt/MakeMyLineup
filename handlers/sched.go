package handlers

import (
    "github.com/6Matt/se390-internal/libhttp"
    "github.com/6Matt/se390-internal/scheduler"
    "github.com/gorilla/mux"
    "html/template"
    "net/http"
)

func GetSched(w http.ResponseWriter, r *http.Request) {
    //sched := scheduler.ScheduleByDay(ScheduleByLocation(mux.Vars(r)["festID"]))
    sched := scheduler.GetAllEvents(mux.Vars(r)["lastID"], mux.Vars(r)["festID"])

    w.Header().Set("Content-Type", "text/html")
    tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/sched.html.tmpl")
    if err != nil {
        libhttp.HandleErrorJson(w, err)
        return
    }

    tmpl.Execute(w, sched)
}