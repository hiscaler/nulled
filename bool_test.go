package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool_NewBool(t *testing.T) {
	valid := NewBool(true, true)
	assert.True(t, valid.Valid)
	assert.True(t, valid.Bool)

	invalid := NewBool(false, false)
	assert.False(t, invalid.Valid)
}

func TestBool_BoolFrom(t *testing.T) {
	valid := BoolFrom(true)
	assert.True(t, valid.Valid)
	assert.True(t, valid.Bool)
}

func TestBool_BoolFromPtr(t *testing.T) {
	b := true
	valid := BoolFromPtr(&b)
	assert.True(t, valid.Valid)
	assert.True(t, valid.Bool)

	invalid := BoolFromPtr(nil)
	assert.False(t, invalid.Valid)
}

func TestBool_MarshalJSON(t *testing.T) {
	validTrue := NewBool(true, true)
	data, err := json.Marshal(validTrue)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`true`), data)

	validFalse := NewBool(false, true)
	data, err = json.Marshal(validFalse)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`false`), data)

	invalid := NewBool(false, false)
	data, err = json.Marshal(invalid)
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), data)
}

func TestBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		inputJSON []byte
		want      Bool
		wantErr   bool
	}{
		{name: "valid true", inputJSON: []byte(`true`), want: NewBool(true, true)},
		{name: "valid false", inputJSON: []byte(`false`), want: NewBool(false, true)},
		{name: "null", inputJSON: []byte(`null`), want: NewBool(false, false)},
		{name: "invalid json", inputJSON: []byte(`"not a bool"`), wantErr: true},
		{name: "empty string", inputJSON: []byte(`""`), wantErr: true}, // Unmarshaling "" into a bool is an error
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Bool
			err := json.Unmarshal(tt.inputJSON, &b)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, b)
			}
		})
	}
}

func TestBool_UnmarshalText(t *testing.T) {
	tests := []struct {
		name      string
		inputText []byte
		want      Bool
		wantErr   bool
	}{
		{name: "valid true", inputText: []byte(`true`), want: NewBool(true, true)},
		{name: "valid 1", inputText: []byte(`1`), want: NewBool(true, true)},
		{name: "valid false", inputText: []byte(`false`), want: NewBool(false, true)},
		{name: "valid 0", inputText: []byte(`0`), want: NewBool(false, true)},
		{name: "empty bytes", inputText: []byte(``), want: NewBool(false, false)},
		{name: "null string", inputText: []byte(`null`), want: NewBool(false, false)},
		{name: "invalid string", inputText: []byte(`abc`), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Bool
			err := b.UnmarshalText(tt.inputText)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, b)
			}
		})
	}
}

func TestBool_EncodeValues(t *testing.T) {
	key := "test_bool"

	t.Run("valid true", func(t *testing.T) {
		v := &url.Values{}
		b := NewBool(true, true)
		err := b.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.Equal(t, "1", v.Get(key))
	})

	t.Run("valid false", func(t *testing.T) {
		v := &url.Values{}
		b := NewBool(false, true)
		err := b.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.Equal(t, "0", v.Get(key))
	})

	t.Run("invalid", func(t *testing.T) {
		v := &url.Values{}
		b := NewBool(false, false)
		err := b.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.False(t, v.Has(key))
	})
}

func TestBool_ValueOrZero(t *testing.T) {
	validTrue := NewBool(true, true)
	assert.True(t, validTrue.ValueOrZero())

	validFalse := NewBool(false, true)
	assert.False(t, validFalse.ValueOrZero())

	invalid := NewBool(false, false)
	assert.False(t, invalid.ValueOrZero())
}

func TestBool_GobEncoding(t *testing.T) {
	// Test valid Bool
	var validBuf bytes.Buffer
	valid := NewBool(true, true)
	err := gob.NewEncoder(&validBuf).Encode(valid)
	assert.NoError(t, err)

	var decodedValid Bool
	err = gob.NewDecoder(&validBuf).Decode(&decodedValid)
	assert.NoError(t, err)
	assert.Equal(t, valid, decodedValid)

	// Test invalid Bool
	var invalidBuf bytes.Buffer
	invalid := NewBool(false, false)
	err = gob.NewEncoder(&invalidBuf).Encode(invalid)
	assert.NoError(t, err)

	var decodedInvalid Bool
	err = gob.NewDecoder(&invalidBuf).Decode(&decodedInvalid)
	assert.NoError(t, err)
	assert.Equal(t, invalid, decodedInvalid)
}
