package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Custom general types
type SQLNullString struct {
	sql.NullString
}

func (ns SQLNullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil) // or return `[]byte("null")` for explicit `null`
}
func (ns *SQLNullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.String = *s
		ns.Valid = true
	} else {
		ns.Valid = false
	}
	return nil
}

type SQLNullTime struct {
	sql.NullTime
}
