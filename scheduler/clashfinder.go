package scheduler

import (
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"os"
	"strings"
	"strconv"
)

// Helpers
const endpoint = "http://clashfinder.com/data/"
const ctLayout = "2006-01-02 15:04"

func getJson(url string, target interface{}) error {
	fmt.Println("getting json from:", url)
	r, err := http.Get(url)
	if err != nil {
    	return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
func writeJsonFile(source interface{}, path string) error {
	b, err := json.Marshal(source)
	if err != nil {
    	return err
	}
	f, err := os.Create(path)
	if err != nil {
    	return err
	}
	defer f.Close()

	_, e := f.Write(b)
	return e
}

type ETime struct { time.Time }
func (ct *ETime) UnmarshalJSON(b []byte) (err error) {
    if b[0] == '"' && b[len(b)-1] == '"' {
        b = b[1 : len(b)-1]
    }
    ct.Time, err = time.Parse(ctLayout, string(b))
    return err
}
type DTime struct { time.Time }
func (ct *DTime) UnmarshalJSON(b []byte) (err error) {
    if b[0] == '"' && b[len(b)-1] == '"' {
        b = b[1 : len(b)-1]
    }
    unixTime, err := strconv.ParseInt(string(b), 10, 64)
    ct.Time = (time.Unix(unixTime, 0)).In(time.UTC)
    return err
}

type Festival struct {
    Id 			string	`json:"name"`
    Name 		string	`json:"desc"`
    StartDate	DTime	`json:"startDate"`
    IsCore		bool	`json:"coreClashfinder"`
}

type Event struct {
	Name 	string
	Start 	ETime
	End 	ETime
}

type Location struct {
	Name 	string
	Events 	[]Event
}

// Get list of festival names
func getFestivalNames(festivals []Festival, onlyCore bool) []string {
	names := make([]string, 1)
	for _, f := range festivals {
		if (onlyCore && !f.IsCore) { continue }
		names = append(names, f.Name)
	}
	return names
}

// Get recent festivals, make names unique, and prioritize core clashfinders
func filterFestivals(unfiltered []Festival) []Festival {
	nameToFestival := map[string]Festival{}
	for _, f := range unfiltered {
		// Discard festivals from more than 3 years ago
		if (f.StartDate.Before(time.Now().AddDate(-3, 0, 0))) { continue; }

		lcName := strings.ToLower(f.Name)
		if fInMap, isInMap := nameToFestival[lcName]; isInMap {
			if (!fInMap.IsCore && f.IsCore) {
				nameToFestival[lcName] = f;
			}
		} else {
			nameToFestival[lcName] = f;
		}
	}
	filtered := make([]Festival, 0, len(nameToFestival))
	for _, v := range nameToFestival {
    	filtered = append(filtered, v)
    }
	return filtered
}

// Get all Festivals
func getFestivals() []Festival {
	m := map[string]Festival{}
	url := endpoint + "events/all.json"
    error := getJson(url, &m)
	if error != nil {
		fmt.Println(error)
	}
	festivals := make([]Festival, 0, len(m))
	for _, v := range m {
    	festivals = append(festivals, v)
    }
	return festivals
}

// Get an Event schedule
func getSchedule(id string) []Location {
	type Response struct {
		Locations []Location
	}

	schedule := Response{}
	url := endpoint + "event/" + id + ".json"
    error := getJson(url, &schedule)
	if error != nil {
		fmt.Println(error)
	}
	return schedule.Locations
}

// Simple main for testing
func main() {
	writeJsonFile(filterFestivals(getFestivals()), "festivals.json")
	//festivals := getFestivals()
	//fmt.Println("\nFESTIVALS:\n\n", festivals, "\n")
	//sched := getSchedule("osheaga2016official")
	//fmt.Println("\nSCHEDULE:\n\n", sched, "\n")
}
