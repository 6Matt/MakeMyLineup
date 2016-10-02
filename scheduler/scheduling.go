package scheduler

import (
	//"fmt"
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
	//rankings := RankingsByName(RankArtists(lastFMID, ArtistList(schedByLoc)))
	events := ToSchedEvent(schedByLoc)

	/*
	SCHEDULE events
	*/

	return EventsByDay(events)
}