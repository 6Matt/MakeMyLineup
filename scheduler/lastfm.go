package scheduler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Helpers
func getJson(url string, target interface{}) error {
	url = strings.Replace(url, " ", "%20", -1)
	 fmt.Println("getting json from:", url)
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	// fmt.Println("response:", r)
	return json.NewDecoder(r.Body).Decode(target)
}

const apiKey = "7c5cb1bf8a097b3491633081c1e62ff8"
const endpoint = "http://ws.audioscrobbler.com/2.0/?method="

func getLastFMJson(query string, properties map[string]string, limit int, autocorrect bool, target interface{}) {
	url := endpoint + query + "&api_key=" + apiKey
	for key, val := range properties {
		url = url + "&" + key + "=" + val
	}
	url = url + "&limit=" + strconv.Itoa(limit)
	if (autocorrect) { url = url + "&autocorrect=1" }
	url = url + "&format=json"
	error := getJson(url, &target)
	if error != nil {
		fmt.Printf("getLastFMJson(%s)\nerror: %#v\n", url, error)
	}
}

type Artist struct {
	Name string
	Mbid string

	// add other needed properties here like:
	//		Url string
	//
	// NOTES:
	// 	- 	Var names must start with capitals for JSON parser to see it
	// 	- 	Var names must match json name (and entire structer) or be
	//			followed with `json: "<JSON NAME>"`
	//			OUTER QUOTES ARE NOT ' THEY ARE ` ONE WORKS, OTHER DOESN'T
	// 	- 	https://mholt.github.io/json-to-go/
	//			json -> go-struct converter
}

// Get Similar Artists to Artist
func getSimilarArtists(propertyName, propertyValue string, count int) []Artist {
	type Response struct {
		Similarartists struct {
			Artist_info []Artist `json:"artist"`
		}
	}

	similarArtists := Response{}
	getLastFMJson("artist.getsimilar", map[string]string{propertyName: propertyValue}, count, false, &similarArtists)
	return similarArtists.Similarartists.Artist_info
}

func getSimilarArtistsByName(name string, count int) []Artist {
	return getSimilarArtists("artist", name, count)
}

func getSimilarArtistsByID(mbid string, count int) []Artist {
	return getSimilarArtists("mbid", mbid, count)
}

// Artist fetching
func getAritstByName(name string) Artist {
	type Response struct {
		Artist_info Artist `json:"artist"`
	}

	artist := Response{}
	getLastFMJson("artist.getInfo", map[string]string{"artist": name}, 1, true, &artist)
	return artist.Artist_info
}

// Get Artists from user's Library
func getArtistsForUser(user string, count int) []Artist {
	type Response struct {
		Artists struct {
			Artist_info []Artist `json:"artist"`
		}
	}

	library := Response{}
	getLastFMJson("library.getartists", map[string]string{"user": user}, count, false, &library)
	return library.Artists.Artist_info
}

// Simple main for testing
/*
func main() {
	artists := getArtistsForUser("devrevan", 20)
	fmt.Println("\nARTISTS:\n\n", artists, "\n")

	if len(artists) > 1 {
		artist := artists[0]
		similar := getSimilarArtistsByName(artist.Name, 5)
		fmt.Println("\nSIMILAR TO", artist.Name, ":\n", similar)

		artist = artists[1]
		similar = getSimilarArtistsByID(artist.Mbid, 7)
		fmt.Println("\nSIMILAR TO", artist.Name, " (by mbid):\n", similar)
	} else {
		fmt.Println("No artists to find similarities with")
	}
}
*/
