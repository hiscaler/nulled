package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"strconv"

	"gopkg.in/guregu/null.v4"
)

type Bool null.Bool

func NewBool(b bool, valid bool) Bool {
	return Bool(null.NewBool(b, valid))
}

func BoolFrom(b bool) Bool {
	return NewBool(b, true)
}

func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

func (b Bool) ValueOrZero() bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

func (b Bool) EncodeValues(key string, v *url.Values) error {
	if !b.Valid {
		return nil
	}

	if b.Bool {
		v.Set(key, "1")
	} else {
		v.Set(key, "0")
	}
	return nil
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(b.Bool)
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	var bo null.Bool
	if err := json.Unmarshal(data, &bo); err != nil {
		b.Valid = false
		return err
	}
	b.Bool = bo.Bool
	b.Valid = bo.Valid
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *Bool) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" || s == "null" {
		b.Valid = false
		return nil
	}
	var err error
	b.Bool, err = strconv.ParseBool(s)
	if err != nil {
		b.Valid = false
		return err
	}
	b.Valid = true
	return nil
}

func (b Bool) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(b.Bool)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(b.Valid)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *Bool) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&b.Bool)
	if err != nil {
		return err
	}
	return dec.Decode(&b.Valid)
}
