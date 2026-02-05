package utils

import "time"

var ArgentinaLocation *time.Location

func init() {
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		ArgentinaLocation = time.FixedZone("ART", -3*60*60)
	} else {
		ArgentinaLocation = loc
	}
}

func Now() time.Time {
	return time.Now().In(ArgentinaLocation)
}

func NowUTC() time.Time {
	return time.Now().UTC()
}
