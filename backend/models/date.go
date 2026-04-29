package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type DateOnly struct {
	Valid bool
	Time  time.Time
}

func NewDateOnly(t time.Time) DateOnly {
	return DateOnly{Valid: true, Time: t}
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(d.Time.Format("2006-01-02"))
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == `""` {
		d.Valid = false
		return nil
	}
	// Strip quotes
	if len(s) >= 2 && s[0] == '"' {
		s = s[1 : len(s)-1]
	}

	// Try date-only
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		d.Time = t
		d.Valid = true
		return nil
	}
	// Try ISO 8601
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		d.Time = t
		d.Valid = true
		return nil
	}
	return fmt.Errorf("cannot parse date: %q", s)
}

func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Valid = false
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		d.Valid = true
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			t, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}
		}
		d.Time = t
		d.Valid = true
		return nil
	}
	return fmt.Errorf("cannot scan date: %v", value)
}

func (d DateOnly) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time.Format("2006-01-02"), nil
}
