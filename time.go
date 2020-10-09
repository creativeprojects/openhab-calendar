package main

import "time"

func getTomorrow(from time.Time) time.Time {
	tomorrow := from.AddDate(0, 0, 1)
	return tomorrow
}
