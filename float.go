package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"strconv"

	"gopkg.in/guregu/null.v4"
)

type Float null.Float

func NewFloat(i float64, valid bool) Float {
	return Float(null.NewFloat(i, valid))
}

func FloatFrom(i float64) Float {
	return NewFloat(i, true)
}

func FloatFromPtr(f *float64) Float {
	if f == nil {
		return NewFloat(0, false)
	}
	return NewFloat(*f, true)
}

func (f Float) ValueOrZero() float64 {
	if !f.Valid {
		return 0
	}
	return f.Float64
}

func (f Float) EncodeValues(key string, v *url.Values) error {
	if !f.Valid {
		return nil
	}

	v.Set(key, strconv.FormatFloat(f.Float64, 'f', -1, 64))
	return nil
}

func (f Float) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(f.Float64)
}

func (f *Float) UnmarshalJSON(data []byte) error {
	temp := null.Float{}
	if err := json.Unmarshal(data, &temp); err != nil {
		f.Valid = false
		return err
	}
	f.Float64 = temp.Float64
	f.Valid = temp.Valid
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (f *Float) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" || s == "null" {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(s, 64)
	if err != nil {
		f.Valid = false
		return err
	}
	f.Valid = true
	return nil
}

func (f Float) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(f.Float64)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(f.Valid)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (f *Float) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&f.Float64)
	if err != nil {
		return err
	}
	return dec.Decode(&f.Valid)
}
