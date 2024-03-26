package types

import (
	"fmt"
	"strings"
	"time"
)

var layoutDateTime = "02.01.2006 15:04"
var layoutDate = "02.01.2006"

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), "\"")
	if value == "" || value == "null" {
		return nil
	}

	var tParsed time.Time
	var err error

	if len(value) > len(layoutDate) {
		tParsed, err = time.Parse(layoutDateTime, value)
	} else {
		tParsed, err = time.Parse(layoutDate, value)
	}

	if err != nil {
		return err
	}

	*t = Time(tParsed)

	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	var layoutToUse string
	if time.Time(*t).Hour() == 0 && time.Time(*t).Minute() == 0 && time.Time(*t).Second() == 0 {
		layoutToUse = layoutDate
	} else {
		layoutToUse = layoutDateTime
	}

	stamp := fmt.Sprintf("\"%s\"", time.Time(*t).Format(layoutToUse))
	return []byte(stamp), nil
}
