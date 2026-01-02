package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_NewString(t *testing.T) {
	validString := NewString("hello", true)
	assert.True(t, validString.Valid)
	assert.Equal(t, "hello", validString.String)

	invalidString := NewString("", false)
	assert.False(t, invalidString.Valid)
}

func TestString_StringFrom(t *testing.T) {
	validString := StringFrom("world")
	assert.True(t, validString.Valid)
	assert.Equal(t, "world", validString.String)

	invalidString := StringFrom("")
	assert.False(t, invalidString.Valid)
}

func TestString_StringFromPtr(t *testing.T) {
	s := "hello world"
	validString := StringFromPtr(&s)
	assert.True(t, validString.Valid)
	assert.Equal(t, s, validString.String)

	invalidString := StringFromPtr(nil)
	assert.False(t, invalidString.Valid)

	s = " "
	invalidWhitespace := StringFromPtr(&s)
	assert.False(t, invalidWhitespace.Valid, "should be invalid for whitespace only string")
}

func TestString_MarshalJSON(t *testing.T) {
	validString := NewString("hello", true)
	data, err := json.Marshal(validString)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"hello"`), data)

	invalidString := NewString("", false)
	data, err = json.Marshal(invalidString)
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), data)
}

func TestString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		inputJSON []byte
		want      String
		wantErr   bool
	}{
		{name: "valid string", inputJSON: []byte(`"hello"`), want: NewString("hello", true)},
		{name: "null", inputJSON: []byte(`null`), want: NewString("", false)},
		{name: "empty string", inputJSON: []byte(`""`), want: NewString("", false)},
		{name: "invalid json", inputJSON: []byte(`not a string`), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s String
			err := json.Unmarshal(tt.inputJSON, &s)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, s)
			}
		})
	}
}

func TestString_UnmarshalText(t *testing.T) {
	tests := []struct {
		name      string
		inputText []byte
		want      String
		wantErr   bool
	}{
		{name: "valid string", inputText: []byte(`hello`), want: NewString("hello", true)},
		{name: "empty bytes", inputText: []byte(``), want: NewString("", false)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s String
			err := s.UnmarshalText(tt.inputText)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, s)
			}
		})
	}
}

func TestString_EncodeValues(t *testing.T) {
	key := "test_string"

	t.Run("valid", func(t *testing.T) {
		v := &url.Values{}
		s := NewString("hello world", true)
		err := s.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", v.Get(key))
	})

	t.Run("invalid", func(t *testing.T) {
		v := &url.Values{}
		s := NewString("", false)
		err := s.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.False(t, v.Has(key))
	})

	t.Run("empty string", func(t *testing.T) {
		v := &url.Values{}
		s := NewString(" ", true) // Whitespace only
		err := s.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.False(t, v.Has(key), "should not encode whitespace only string")
	})
}

func TestString_ValueOrZero(t *testing.T) {
	validString := NewString("hello", true)
	assert.Equal(t, "hello", validString.ValueOrZero())

	invalidString := NewString("", false)
	assert.Equal(t, "", invalidString.ValueOrZero())
}

func TestString_GobEncoding(t *testing.T) {
	// Test valid String
	var validBuf bytes.Buffer
	validString := NewString("hello gob", true)
	err := gob.NewEncoder(&validBuf).Encode(validString)
	assert.NoError(t, err)

	var decodedValid String
	err = gob.NewDecoder(&validBuf).Decode(&decodedValid)
	assert.NoError(t, err)
	assert.Equal(t, validString, decodedValid)

	// Test invalid String
	var invalidBuf bytes.Buffer
	invalidString := NewString("", false)
	err = gob.NewEncoder(&invalidBuf).Encode(invalidString)
	assert.NoError(t, err)

	var decodedInvalid String
	err = gob.NewDecoder(&invalidBuf).Decode(&decodedInvalid)
	assert.NoError(t, err)
	assert.Equal(t, invalidString, decodedInvalid)
}

func TestString_NullValue(t *testing.T) {
	// 验证有效 String 的 NullValue
	validString := NewString("test", true)
	assert.True(t, validString.NullValue().Valid)
	assert.Equal(t, "test", validString.NullValue().String)

	// 验证无效 String 的 NullValue
	invalidString := NewString("", false)
	assert.False(t, invalidString.NullValue().Valid)
	assert.Equal(t, "", invalidString.NullValue().String)
}
