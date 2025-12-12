package stats

import (
	"fmt"
	"time"
)

func durationToString(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%dm%ds%03dms", minutes, seconds, milliseconds)
}
