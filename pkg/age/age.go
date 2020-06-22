package age

import (
	"fmt"
	"time"
)

// GetAge returns a human-readable age string based on the amount of time elapsed since the given time
func GetAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d.Seconds() <= 60:
		return fmt.Sprintf("%vs", int(d.Seconds()))
	case d.Minutes() <= 60:
		return fmt.Sprintf("%vm", int(d.Minutes()))
	case d.Hours() <= 24:
		return fmt.Sprintf("%vh", int(d.Hours()))
	case d.Hours() > 24:
		return fmt.Sprintf("%vd", int(d.Hours()/24))
	default:
		return ""
	}
}
