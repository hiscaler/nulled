package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"strings"

	"gopkg.in/guregu/null.v4"
)

type String null.String

func NewString(s string, valid bool) String {
	return String(null.NewString(s, valid))
}
func StringFrom(s string) String {
	s = strings.TrimSpace(s)
	if s == "" {
		return NewString("", false)
	}
	return NewString(s, true)
}

func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	return StringFrom(*s)
}

func (m String) ValueOrZero() string {
	if !m.Valid {
		return ""
	}
	return m.String
}

func (m String) EncodeValues(key string, v *url.Values) error {
	if !m.Valid || strings.TrimSpace(m.String) == "" {
		return nil
	}
	v.Set(key, m.String)
	return nil
}

func (m String) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(m.String)
}

func (m *String) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		m.Valid = false
		m.String = ""
	} else {
		// an empty string is considered null
		m.Valid = *s != ""
		m.String = *s
	}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (m *String) UnmarshalText(text []byte) error {
	s := string(text)
	// an empty string is considered null
	if s == "" {
		m.Valid = false
		m.String = ""
	} else {
		m.Valid = true
		m.String = s
	}
	return nil
}

func (m String) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m.String)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(m.Valid)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *String) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&m.String)
	if err != nil {
		return err
	}
	return dec.Decode(&m.Valid)
}

func (m String) NullValue() null.String {
	if m.Valid {
		return null.StringFrom(m.String)
	}
	return null.NewString("", false)
}
