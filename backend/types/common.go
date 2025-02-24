package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

type NameLocalization struct {
	Name string `json:"name"`
}

type NullString sql.NullString

func (s NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte("null"), nil
}

func (s NullString) Value() (driver.Value, error) {
	if s.Valid {
		return s.String, nil
	}
	return []byte("null"), nil
}

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

type NullTime sql.NullTime

func (s NullTime) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Time)
	}
	return []byte("null"), nil
}

func (s NullTime) Value() (driver.Value, error) {
	if s.Valid {
		return s.Time, nil
	}
	return []byte("null"), nil
}

func (ns *NullTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullTime{s.Time, false}
	} else {
		*ns = NullTime{s.Time, true}
	}

	return nil
}

type NullInt32 sql.NullInt32

func (s NullInt32) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Int32)
	}
	return []byte("null"), nil
}

func (s NullInt32) Value() (driver.Value, error) {
	if s.Valid {
		return s.Int32, nil
	}
	return []byte("null"), nil
}

func (ns *NullInt32) Scan(value interface{}) error {
	var s sql.NullInt32
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullInt32{s.Int32, false}
	} else {
		*ns = NullInt32{s.Int32, true}
	}

	return nil
}

type NullFloat64 sql.NullFloat64

func (s NullFloat64) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Float64)
	}
	return []byte("null"), nil
}

func (s NullFloat64) Value() (driver.Value, error) {
	if s.Valid {
		return s.Float64, nil
	}
	return []byte("null"), nil
}

func (ns *NullFloat64) Scan(value interface{}) error {
	var s sql.NullFloat64
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullFloat64{s.Float64, false}
	} else {
		*ns = NullFloat64{s.Float64, true}
	}

	return nil
}

type NullUserProfileSimple struct {
	Profile UserProfileSimple
	Valid   bool
}

func (s NullUserProfileSimple) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Profile)
	}
	return []byte("null"), nil
}

func (s NullUserProfileSimple) Value() (driver.Value, error) {
	if s.Valid {
		return s.Profile, nil
	}
	return []byte("null"), nil
}
