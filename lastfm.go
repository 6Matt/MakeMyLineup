package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

// Helpers
func getJson(url string, target interface{}) error {
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
func getSimilarArtists(propertyName, propertyValue string) []Artist {
	type Response struct {
		Similarartists struct {
			Artist_info []Artist `json:"artist"`
		}
	}

	similarArtists := Response{}
	url := endpoint + "artist.getsimilar&api_key=" +
			apiKey + "&" + propertyName + "=" + propertyValue + "&format=json"
    error := getJson(url, &similarArtists)
	if error != nil {
		fmt.Println(error)
	}
	return similarArtists.Similarartists.Artist_info
}

func getSimilarArtistsByName(name string) []Artist {
	return getSimilarArtists("artist", name)
}

func getSimilarArtistsByID(mbid string) []Artist {
	return getSimilarArtists("mbid", mbid)
}

// Get Artists from user's Library
func getArtistsForUser(user string) []Artist {
	type Response struct {
		Artists struct {
			Artist_info []Artist `json:"artist"`
		}
	}

	library := Response{}
	url := endpoint + "library.getartists&api_key=" +
			apiKey + "&user=" + user + "&format=json"
    error := getJson(url, &library)
	if error != nil {
		fmt.Println(error)
	}
	return library.Artists.Artist_info
}

// Simple main for testing
func main() {
	artists := getArtistsForUser("devrevan")
	fmt.Println("\nARTISTS:\n\n", artists, "\n")

	if len(artists) > 1 {
		artist := artists[0]
		similar := getSimilarArtistsByName(artist.Name)
		fmt.Println("\nSIMILAR TO", artist.Name, ":\n", similar)

		artist = artists[1]
		similar = getSimilarArtistsByID(artist.Mbid)
		fmt.Println("\nSIMILAR TO", artist.Name, " (by mbid):\n", similar)	
	} else {
		fmt.Println("No artists to find similarities with")	
	}
}
