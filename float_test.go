package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat_NewFloat(t *testing.T) {
	valid := NewFloat(123.45, true)
	assert.True(t, valid.Valid)
	assert.Equal(t, 123.45, valid.Float64)

	invalid := NewFloat(0, false)
	assert.False(t, invalid.Valid)
}

func TestFloat_FloatFrom(t *testing.T) {
	valid := FloatFrom(123.45)
	assert.True(t, valid.Valid)
	assert.Equal(t, 123.45, valid.Float64)
}

func TestFloat_FloatFromPtr(t *testing.T) {
	f := 123.45
	valid := FloatFromPtr(&f)
	assert.True(t, valid.Valid)
	assert.Equal(t, f, valid.Float64)

	invalid := FloatFromPtr(nil)
	assert.False(t, invalid.Valid)
}

func TestFloat_MarshalJSON(t *testing.T) {
	valid := NewFloat(123.45, true)
	data, err := json.Marshal(valid)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`123.45`), data)

	invalid := NewFloat(0, false)
	data, err = json.Marshal(invalid)
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), data)
}

func TestFloat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		inputJSON []byte
		want      Float
		wantErr   bool
	}{
		{name: "valid float", inputJSON: []byte(`123.45`), want: NewFloat(123.45, true)},
		{name: "zero", inputJSON: []byte(`0`), want: NewFloat(0, true)},
		{name: "null", inputJSON: []byte(`null`), want: NewFloat(0, false)},
		{name: "invalid json", inputJSON: []byte(`"not a float"`), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f Float
			err := json.Unmarshal(tt.inputJSON, &f)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, f)
			}
		})
	}
}

func TestFloat_UnmarshalText(t *testing.T) {
	tests := []struct {
		name      string
		inputText []byte
		want      Float
		wantErr   bool
	}{
		{name: "valid float", inputText: []byte(`123.45`), want: NewFloat(123.45, true)},
		{name: "negative float", inputText: []byte(`-54.321`), want: NewFloat(-54.321, true)},
		{name: "empty bytes", inputText: []byte(``), want: NewFloat(0, false)},
		{name: "null string", inputText: []byte(`null`), want: NewFloat(0, false)},
		{name: "invalid string", inputText: []byte(`abc`), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f Float
			err := f.UnmarshalText(tt.inputText)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, f)
			}
		})
	}
}

func TestFloat_EncodeValues(t *testing.T) {
	key := "test_float"

	t.Run("valid", func(t *testing.T) {
		v := &url.Values{}
		f := NewFloat(123.45, true)
		err := f.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.Equal(t, "123.45", v.Get(key))
	})

	t.Run("invalid", func(t *testing.T) {
		v := &url.Values{}
		f := NewFloat(0, false)
		err := f.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.False(t, v.Has(key))
	})
}

func TestFloat_ValueOrZero(t *testing.T) {
	valid := NewFloat(123.45, true)
	assert.Equal(t, 123.45, valid.ValueOrZero())

	invalid := NewFloat(0, false)
	assert.Equal(t, float64(0), invalid.ValueOrZero())
}

func TestFloat_GobEncoding(t *testing.T) {
	// Test valid Float
	var validBuf bytes.Buffer
	valid := NewFloat(123.45, true)
	err := gob.NewEncoder(&validBuf).Encode(valid)
	assert.NoError(t, err)

	var decodedValid Float
	err = gob.NewDecoder(&validBuf).Decode(&decodedValid)
	assert.NoError(t, err)
	assert.Equal(t, valid, decodedValid)

	// Test invalid Float
	var invalidBuf bytes.Buffer
	invalid := NewFloat(0, false)
	err = gob.NewEncoder(&invalidBuf).Encode(invalid)
	assert.NoError(t, err)

	var decodedInvalid Float
	err = gob.NewDecoder(&invalidBuf).Decode(&decodedInvalid)
	assert.NoError(t, err)
	assert.Equal(t, invalid, decodedInvalid)
}

func TestFloat_NullValue(t *testing.T) {
	// 验证有效 Float 的 NullValue
	validFloat := NewFloat(123.45, true)
	assert.True(t, validFloat.NullValue().Valid)
	assert.Equal(t, float64(123.45), validFloat.NullValue().Float64)

	// 验证无效 Float 的 NullValue
	invalidFloat := NewFloat(0, false)
	assert.False(t, invalidFloat.NullValue().Valid)
	assert.Equal(t, float64(0), invalidFloat.NullValue().Float64)
}
