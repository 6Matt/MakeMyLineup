package scheduler

import (
	//"fmt"
	"math"
	"sync"
)

type Rank struct {
	Val  float64
	Tier int
}

var penalty float64 = 2.0 / 3.0
var adjCount = []int{3, 2}
var zeroRank = Rank{0, 0}

func rankValue(r Rank) int64 {
	return int64(r.Val * math.Pow(penalty, float64(r.Tier)))
}

func rankForNeighbourOf(r Rank) Rank {
	return Rank{r.Val, r.Tier + 1}
}

func getSimilarArtsitsWithCache(artist Artist, limit int, similarArtistCache map[Artist]([]Artist), similarArtistsLock sync.Mutex) []Artist {
	similarArtistsLock.Lock()
	simliarArtists := similarArtistCache[artist]
	similarArtistsLock.Unlock()
	if len(simliarArtists) == 0 {
		simliarArtists = getSimilarArtistsByName(artist.Name, limit)
		similarArtistsLock.Lock()
		similarArtistCache[artist] = simliarArtists
		similarArtistsLock.Unlock()
	}
	return simliarArtists
}

type QueuedArtist struct {
	Artist
	Depth  int
	Parent *QueuedArtist
}

func appendArtists(queue []QueuedArtist, artists []Artist, parent QueuedArtist) []QueuedArtist {
	for _, artist := range artists {
		queue = append(queue, QueuedArtist{artist, parent.Depth + 1, &parent})
	}
	return queue
}

func computeRankForArtist(artist Artist, rankCache map[Artist]Rank, rankCacheLock sync.Mutex, similarArtistCache map[Artist]([]Artist), similarArtistsLock sync.Mutex) {
	// fmt.Println("computeRankForArtist", artist)
	rankCacheLock.Lock()
	_, ok := rankCache[artist]
	rankCacheLock.Unlock()
	if ok {
		return
	}

	rootArtist := QueuedArtist{artist, 0, nil}
	simliarArtists := getSimilarArtsitsWithCache(artist, adjCount[0], similarArtistCache, similarArtistsLock)
	queue := appendArtists([]QueuedArtist{}, simliarArtists, rootArtist)

	foundDepth := len(adjCount) + 1
	for len(queue) > 0 {
		queuedArtist := queue[0]
		queue = queue[1:]
		// fmt.Println("\nevaluating artist", queuedArtist /*, "\nqueue", queue*/)

		if queuedArtist.Depth > foundDepth {
			break
		}

		rankCacheLock.Lock()
		val, ok := rankCache[queuedArtist.Artist]
		rankCacheLock.Unlock()
		if ok {
			// fmt.Println("\nevaluating artist", queuedArtist, "found similar artist", queuedArtist)
			foundDepth = queuedArtist.Depth
			childRank := val
			parent := queuedArtist.Parent
			for parent != nil {
				rankCacheLock.Lock()
				currentRank := rankCache[(*parent).Artist]
				if newRank := rankForNeighbourOf(childRank); rankValue(newRank) > rankValue(currentRank) {
					rankCache[(*parent).Artist] = newRank
					rankCacheLock.Unlock()
					childRank = newRank
					parent = (*parent).Parent
				} else {
					rankCacheLock.Unlock()
					break
				}
			}
		} else if queuedArtist.Depth < len(adjCount) {
			// fmt.Println("\nevaluating artist", queuedArtist, "artist not similar", queuedArtist)
			simliarArtists = getSimilarArtsitsWithCache(queuedArtist.Artist, adjCount[queuedArtist.Depth], similarArtistCache, similarArtistsLock)
			queue = appendArtists(queue, simliarArtists, queuedArtist)
		}
	}
}

func RankArtists(username string, artists []Artist) map[Artist]int64 {
	// get user's library
	library := getArtistsForUser(username, 100)
	// fmt.Println("\n library of ", username, " size ", len(library), "\n", library)

	// assign ranks to artists in user's library
	var cache = map[Artist]Rank{}
	cacheLock := sync.Mutex{}
	for i := len(library) - 1; i >= 0; i-- {
		cache[library[i]] = Rank{float64(len(library) - i), 0}
	}

	// rank given from artists
	var wg sync.WaitGroup
	var rankedArtists = map[Artist]int64{}
	rankedArtistsLock := sync.Mutex{}
	var similarArtistCache = map[Artist]([]Artist){}
	similarArtistCacheLock := sync.Mutex{}
	for _, artist := range artists {
		wg.Add(1)
		// fmt.Println("looping for", artist)
		_artist := artist
		go func() {
			defer wg.Done()
			// fmt.Println("calling computeRankForArtist for", _artist)
			computeRankForArtist(_artist, cache, cacheLock, similarArtistCache, similarArtistCacheLock)
			cacheLock.Lock()
			rank := cache[_artist]
			cacheLock.Unlock()
			rankedArtistsLock.Lock()
			rankedArtists[_artist] = rankValue(rank)
			rankedArtistsLock.Unlock()
		}()
	}

	wg.Wait()
	return rankedArtists
}

/*
// Simple main for testing
func main() {
	a := []Artist{}
	a = append(a, getAritstByName("Twenty One Pilots")) // artist in library (high)
	a = append(a, getAritstByName("Three Days Grace"))  // aritst in library (low)
	a = append(a, getAritstByName("Halsey"))            // not similar to in library (any depth)
	a = append(a, getAritstByName("Tyler Joseph"))      // similar to artist in library (depth 1)
	a = append(a, getAritstByName("Leathermouth"))      // similar to artist in library (depth 2)
	rankedArtists := RankArtists("devrevan", a)
	for key, val := range rankedArtists {
		fmt.Println(key.Name, val)
	}
}
*/
