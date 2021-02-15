package util

import (
	"strings"
	"time"
)

// JSONTime shouldn't exist, but
// the default Golang time serialization
// is dumb.
type JSONTime struct {
	time.Time
}

// MarshalJSON handles serializing JSONTime (i.e., time.Time)
// into a byte array following "%Y-%m-%dT%H:%M:%S.%f" (in normal syntax)
func (j JSONTime) MarshalJSON() ([]byte, error) {
	stamp := j.Time.Format("2006-01-02T15:04:05.000000")
	return []byte(stamp), nil
}

// UnmarshalJSON handles deserializing a byte array
// into a time.Time object according to "%Y-%m-%dT%H:%M:%S.%f" (in normal syntax)
func (j *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02T15:04:05.000000", s)
	if err != nil {
		return err
	}
	j.Time = t
	return nil
}
