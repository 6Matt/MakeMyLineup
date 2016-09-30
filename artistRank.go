package main

import (
	"fmt"
	"math"
)

type Rank struct {
	Val  float64
	Tier int
}

func rankValue(r Rank) int64 {
	return int64(r.Val * math.Pow(penalty, float64(r.Tier)))
}

var depthLimit = 3
var penalty float64 = 2 / 3
var adjCount = []int64{5, 2, 1}

func rankArtists(username string, artists []Artist) map[Artist]int64 {
	fmt.Println("artists to rank: ", artists)

	// get user's library
	library := getArtistsForUser("devrevan", 100)

	fmt.Println("\n library of ", username, " size ", len(library), "\n", library)

	// assign ranks to artists in user's library
	var cache = map[Artist]Rank{}
	length := int64(len(library))
	for i := length - 1; i >= 0; i-- {
		cache[library[i]] = Rank{float64(length - i), 0}
	}

	// rank given from artists
	var rankedArtists = map[Artist]int64{}
	for _, artist := range artists {
		rank := cache[artist]
		rankedArtists[artist] = rankValue(rank)
	}

	fmt.Println("\nartists ranks\n", rankedArtists)

	return rankedArtists
}

func main() {
	a := getArtistsForUser("devrevan", 3)
	a = append(a, Artist{"Artist not in library", "1234"})
	rankArtists("devrevan", a)
}
