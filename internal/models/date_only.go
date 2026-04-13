package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// DateOnly wraps time.Time to accept "YYYY-MM-DD" (and RFC3339) in JSON,
// serialize back as "YYYY-MM-DD", and map to a DATE column in Postgres.
type DateOnly struct {
	time.Time
}

const dateLayout = "2006-01-02"

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		d.Time = time.Time{}
		return nil
	}

	// Try plain date first (YYYY-MM-DD)
	if t, err := time.Parse(dateLayout, s); err == nil {
		d.Time = t
		return nil
	}

	// Fall back to RFC3339 (e.g. "2021-01-07T00:00:00Z")
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		d.Time = t
		return nil
	}

	return fmt.Errorf("cannot parse %q as a date: expected YYYY-MM-DD or RFC3339", s)
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format(dateLayout) + `"`), nil
}

// Value implements driver.Valuer so GORM can write a plain date to Postgres.
func (d DateOnly) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// Scan implements sql.Scanner so GORM can read a date back from Postgres.
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	}
	return fmt.Errorf("DateOnly.Scan: unsupported type %T", value)
}
