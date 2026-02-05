package utils

import "time"

// ArgentinaLocation returns the Argentina timezone (UTC-3)
var ArgentinaLocation *time.Location

func init() {
	// Load Argentina timezone
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		// Fallback to fixed offset if timezone data is not available
		ArgentinaLocation = time.FixedZone("ART", -3*60*60) // UTC-3
	} else {
		ArgentinaLocation = loc
	}
}

// Now returns the current time in Argentina timezone (UTC-3)
func Now() time.Time {
	return time.Now().In(ArgentinaLocation)
}

// NowUTC returns the current time in UTC (for cases where UTC is explicitly needed)
func NowUTC() time.Time {
	return time.Now().UTC()
}
