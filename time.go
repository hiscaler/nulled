package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"time"

	"gopkg.in/guregu/null.v4"
)

type Time null.Time

func NewTime(t time.Time, valid bool) Time {
	return Time(null.NewTime(t, valid))
}

func TimeFrom(t time.Time) Time {
	if t.IsZero() {
		return NewTime(time.Time{}, false)
	}

	return NewTime(t, true)
}

func TimeFromPtr(t *time.Time) Time {
	if t == nil || t.IsZero() {
		return NewTime(time.Time{}, false)
	}

	return NewTime(*t, true)
}

func (m Time) ValueOrZero() time.Time {
	if !m.Valid {
		return time.Time{}
	}

	return m.Time
}

func (m Time) EncodeValues(key string, v *url.Values) error {
	if !m.Valid {
		return nil
	}
	v.Set(key, m.Time.Format("2006-01-02 15:04:05"))
	return nil
}

func (m Time) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(m.Time)
}

func (m *Time) UnmarshalJSON(data []byte) error {
	var t null.Time
	if err := json.Unmarshal(data, &t); err != nil {
		m.Valid = false
		return err
	}
	m.Time = t.Time
	m.Valid = t.Valid
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (m *Time) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" || s == "null" {
		m.Valid = false
		return nil
	}
	// Use time.RFC3339 by default for robust parsing
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		m.Valid = false
		return err
	}
	m.Time = parsedTime
	m.Valid = true
	return nil
}

func (m Time) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m.Time)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(m.Valid)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Time) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&m.Time)
	if err != nil {
		return err
	}
	return dec.Decode(&m.Valid)
}
