package valueobject

import (
	"errors"
	"time"
)

var (
	ErrTimeIsRequired = errors.New("Time is required")
)

type Time struct {
	value time.Time
}

func NewTime(t time.Time) Time {
	return Time{value: t}
}

func NewCurrentTime() Time {
	return Time{value: time.Now()}
}

func (t Time) ToString() string {
	return t.value.Format("January 2, 2006 - 3:04PM")
}

func (t Time) ToTime() time.Time {
	return t.value
}

func (t Time) ExtendByHour(hours int) Time {
	return Time{value: t.value.Add(time.Hour * time.Duration(hours))}
}
