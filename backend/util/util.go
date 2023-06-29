package util

import (
	"time"
)

func LoadLocalTime() *time.Location {
	gmt4 := time.FixedZone("GMT+4", 4*60*60)
	return gmt4
}
