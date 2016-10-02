package scheduler

import (
	"fmt"
)

func GetAllEvents(lastFMID string, festivalID string) []SchedEvent {
	schedByLoc := ScheduleByLocation(festivalID)
	rankings := RankArtists(lastFMID, ArtistList(schedByLoc))
	
	for key, val := range rankings {
		fmt.Println(key.Name, val)
	}

	return nil
}