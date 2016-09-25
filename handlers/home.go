package handlers

import (
	"github.com/6Matt/se390-internal/libhttp"
	"html/template"
	"net/http"

	"encoding/json"
)

type Event struct {
    Name 		string	`json:"name"`
    Description string	`json:"desc"`
    Core		bool	`json:"coreClashfinder"`
}

func getEvents(url string, e *[]Event) error {
    r, err := http.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    dec := json.NewDecoder(r.Body)
    var m map[string]Event
    if err := dec.Decode(&m); err != nil {
        println(err)
        return err
    }
    for _, v := range m {
    	*e = append(*e, v)
    }
    return nil
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	events := make([]Event, 1)
    getEvents("http://clashfinder.com/data/events/all.json", &events)

    /*
    for _, e := range events {
    	println(e.Name)
    	println(e.Description)
	}
	*/

	tmpl.Execute(w, nil)
}
