package scheduler

import (
	// "fmt"
	"sort"
)

func RankingsByName(rankings map[Artist]int64) map[string]int64 {
	nameToRank := make(map[string]int64)
	for key, val := range rankings {
		nameToRank[key.Name] = val
	}
	return nameToRank
}

func ScheduledEventsByDay(lastFMID string, festivalID string) []SchedDay {
	schedByLoc := ScheduleByLocation(festivalID)
	rankings := RankingsByName(RankArtists(lastFMID, ArtistList(schedByLoc)))
	events := ToSchedEvent(schedByLoc)


	scheduledEvents := weightedIntreval(events, rankings)

	return EventsByDay(scheduledEvents)
}


//ByFinishTime implements sort.Interface for []SchedEvent based on End time of the event
type ByFinishTime []SchedEvent

func (a ByFinishTime) Len() int 				{ return len(a) }
func (a ByFinishTime) Swap(i, j int) 			{ a[i], a[j] = a[j], a[i] }
func (a ByFinishTime) Less(i,j int) bool 		{ return a[i].End.Time.Before(a[j].End.Time) }

func max(a, b int64) int64 {
	if a < b {
		return b
	} 
	return a
} 

func weightedIntreval(schedByLoc []SchedEvent, rankings map[string]int64) []SchedEvent {

//1 - sort by end time
	sort.Sort(ByFinishTime(schedByLoc))

//2 - calculate the event the ends last before the start of each event
	lastBefore := make([]int, len(schedByLoc))

	//initialize to -1 for the starting events
	for i := range lastBefore {
		lastBefore[i] = -1;
	}

	for key, val := range schedByLoc {

		for i := key - 1; i >= 0; i-- {
			//assume people can teleport
			if !val.Start.Time.Before(schedByLoc[i].End.Time) {
				lastBefore[key] = i
				break;
			}
		}
	}
 
//3 - do weighted interval scheduling
	maxValues := make([]int64, len(schedByLoc)+1)
	maxValues[0] = 0;

	for key, val := range schedByLoc {
		value := rankings[val.Name] + 1 
		maxValues[key+1] = max(value + maxValues[lastBefore[key]+1], maxValues[key])
	}

//4 - change the Scheduled bool to true for chosen events (backtracking)
	
	for length := len(schedByLoc)-1; length >= 0; length-- {
		value := rankings[schedByLoc[length].Name] + 1

		if maxValues[lastBefore[length]+1] + value > maxValues[length] {
			schedByLoc[length].Scheduled = true;
			length = lastBefore[length]+1
		} else {
			length--;
		}
	}

//5 - return scheduled items

	return schedByLoc;
}