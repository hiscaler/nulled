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

func (t Time) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}

	return t.Time
}

func (t Time) EncodeValues(key string, v *url.Values) error {
	if !t.Valid {
		return nil
	}
	v.Set(key, t.Time.Format("2006-01-02 15:04:05"))
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var tt null.Time
	if err := json.Unmarshal(data, &tt); err != nil {
		t.Valid = false
		return err
	}
	t.Time = tt.Time
	t.Valid = tt.Valid
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Time) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" || s == "null" {
		t.Valid = false
		return nil
	}
	// Use time.RFC3339 by default for robust parsing
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Valid = false
		return err
	}
	t.Time = parsedTime
	t.Valid = true
	return nil
}

func (t Time) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(t.Time)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(t.Valid)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *Time) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&t.Time)
	if err != nil {
		return err
	}
	return dec.Decode(&t.Valid)
}

func (t Time) NullValue() null.Time {
	if t.Valid {
		return null.TimeFrom(t.Time)
	}
	return null.NewTime(time.Time{}, false)
}
