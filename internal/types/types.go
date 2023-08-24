package types

import (
	"database/sql"
	"encoding/json"
	"reflect"

	"github.com/sqlc-dev/pqtype"
)

type Json map[string]any

type NullInt32 sql.NullInt32

// Scan implements the Scanner interface for NullInt32
func (ni *NullInt32) Scan(value interface{}) error {
	var i sql.NullInt32
	if err := i.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = NullInt32{i.Int32, false}
	} else {
		*ni = NullInt32{i.Int32, true}
	}

	return nil
}

// MarshalJSON for NullInt32
func (ni NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(ni.Int32)
}

// UnmarshalJSON for NullInt32
func (ni *NullInt32) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int32)
	ni.Valid = (err == nil)
	return err
}

// NullString is an alias for sql.NullString data type
type NullString sql.NullString

// Scan implements the Scanner interface for NullString
func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}

	return nil
}

// MarshalJSON for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}

// ----------------------------------

type NullRawMessage pqtype.NullRawMessage

// Scan implements the Scanner interface for NullRawMessage
func (ni *NullRawMessage) Scan(value interface{}) error {
	var i pqtype.NullRawMessage
	if err := i.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = NullRawMessage{i.RawMessage, false}
	} else {
		*ni = NullRawMessage{i.RawMessage, true}
	}

	return nil
}

// MarshalJSON for NullRawMessage
func (ni NullRawMessage) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(ni.RawMessage)
}

// UnmarshalJSON for NullRawMessage
func (ni *NullRawMessage) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.RawMessage)
	ni.Valid = (err == nil)
	return err
}
