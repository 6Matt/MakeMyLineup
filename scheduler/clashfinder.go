package scheduler

import (
	"fmt"
	"time"
	"encoding/json"
	"os"
	"strings"
	"strconv"
	"sort"
    "html/template"
)

// Helpers
const cfEndPoint = "http://clashfinder.com/data/"
const ctLayout = "2006-01-02 15:04"
const dayLayout = "2006-01-02"
const friendlyLayout = "Monday, January 2"

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
func (e Event) FormatEventName() template.HTML {

	var formatted string = ""
	for _, rune_value := range e.Name {
		if rune_value < 128 {
			formatted += string(rune_value);
		} else {
			formatted += "&#" + strconv.FormatInt(int64(rune_value), 10) + ";"
		}
	}

    return template.HTML(formatted)
}
type ByStart []Event
func (a ByStart) Len() int           { return len(a) }
func (a ByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStart) Less(i, j int) bool { return a[i].Start.Time.Before(a[j].Start.Time) }

type SchedEvent struct {
	Name 		string
	Location 	string
	Scheduled 	bool
	Start 		ETime
	End 		ETime
}

type Location struct {
	Name 	string
	Events 	[]Event
}
func (l Location) FormatLocationName() template.HTML {

	var formatted string = ""
	for _, rune_value := range l.Name {
		if rune_value < 128 {
			formatted += string(rune_value);
		} else {
			formatted += "&#" + strconv.FormatInt(int64(rune_value), 10) + ";"
		}
	}

    return template.HTML(formatted)
}
type ByName []Location
func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type Day struct {
	Encoded 	string
	Date 		string
	Locations 	[]Location
}
type ByEncoded []Day
func (a ByEncoded) Len() int           { return len(a) }
func (a ByEncoded) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByEncoded) Less(i, j int) bool { return a[i].Encoded < a[j].Encoded }

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
	url := cfEndPoint + "events/all.json"
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
func ScheduleByLocation(id string) []Location {
	type Response struct {
		Locations []Location
	}

	schedule := Response{}
	url := cfEndPoint + "event/" + id + ".json"
    error := getJson(url, &schedule)
	if error != nil {
		fmt.Println(error)
	}
	return schedule.Locations
}

func (e *Event) getDay() string {
	return e.Start.Format(dayLayout)
}
func makeFriendly(enc string) string {
	time, _ := time.Parse(dayLayout, string(enc))
	return time.Format(friendlyLayout)
}

// Sort schedule by day, locations by name, and events by time
func ScheduleByDay(sched []Location) []Day {
	dayToLocToEvt := make(map[string]map[string][]Event)
	for _, loc := range sched {
    	for _, evt := range loc.Events {
			if _, isInMap := dayToLocToEvt[evt.getDay()]; !isInMap {
				dayToLocToEvt[evt.getDay()] = make(map[string][]Event)
			}
    		dayToLocToEvt[evt.getDay()][loc.Name] = append(dayToLocToEvt[evt.getDay()][loc.Name], evt)
    	}
    }

    days := make([]Day, 0, len(dayToLocToEvt))
	for d, lToE := range dayToLocToEvt {
		locations := make([]Location, 0, len(lToE))
    	for l, e := range lToE {
    		locations = append(locations, Location{l, e})
    	}
    	days = append(days, Day{d, makeFriendly(d), locations})
    }

    sort.Sort(ByEncoded(days))
    for _, d := range days {
    	sort.Sort(ByName(d.Locations))
    	for _, l := range d.Locations {
    		sort.Sort(ByStart(l.Events))
    	}
    }

    return days
}


func ArtistList(byLocation []Location) []Artist {
	names := make(map[string]bool)
	for _, loc := range byLocation {
    	for _, evt := range loc.Events {
    		if _, isInMap := names[evt.Name]; !isInMap {
				names[evt.Name] = true
			}
    	}
    }

    artists := make([]Artist, 0, len(names))
	for a, _ := range names {
		artists = append(artists, getAritstByName(a))
	}
	return artists
}


/*
// Simple main for testing
func main() {
	writeJsonFile(filterFestivals(getFestivals()), "festivals.json")
	//festivals := getFestivals()
	//fmt.Println("\nFESTIVALS:\n\n", festivals, "\n")
	//sched := getSchedule("osheaga2016official")
	//fmt.Println("\nSCHEDULE:\n\n", sched, "\n")
}
*/