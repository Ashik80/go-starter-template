package valueobject

import "time"

type Time time.Time

func (t *Time) ToString() string {
	return time.Time(*t).Format("January 2, 2006 - 3:04PM")
}

func (t *Time) ToTime() time.Time {
	return time.Time(*t)
}
