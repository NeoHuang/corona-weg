package core

import "time"

type Epidemic struct {
	Infections int
	Deaths     int
	Timestamp  time.Time
}

type EpidemicMap map[string]Epidemic
