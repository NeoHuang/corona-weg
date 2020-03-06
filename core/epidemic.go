package core

import "time"

type Epidemic struct {
	Infections int
	Deaths     int
	Timestamp  time.Time
	SourceApi  string
}

type EpidemicMap map[string]Epidemic
