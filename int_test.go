package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt_NewInt(t *testing.T) {
	validInt := NewInt(10, true)
	assert.True(t, validInt.Valid)
	assert.Equal(t, int64(10), validInt.Int64)

	invalidInt := NewInt(0, false)
	assert.False(t, invalidInt.Valid)
}

func TestInt_IntFrom(t *testing.T) {
	i := int64(42)
	validInt := IntFrom(i)
	assert.True(t, validInt.Valid)
	assert.Equal(t, i, validInt.Int64)
}

func TestInt_IntFromPtr(t *testing.T) {
	i := int64(42)
	validInt := IntFromPtr(&i)
	assert.True(t, validInt.Valid)
	assert.Equal(t, i, validInt.Int64)

	nullInt := IntFromPtr(nil)
	assert.False(t, nullInt.Valid)
}

func TestInt_MarshalJSON(t *testing.T) {
	// 测试有效的 Int 值
	validInt := IntFrom(42)
	data, err := json.Marshal(validInt)
	assert.NoError(t, err)
	assert.Equal(t, []byte("42"), data)

	// 测试无效的 Int 值
	invalidInt := NewInt(0, false)
	data, err = json.Marshal(invalidInt)
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), data)
}

func TestInt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		inputJSON []byte
		wantInt   Int
		wantErr   bool
	}{
		{
			name:      "zero",
			inputJSON: []byte(`0`),
			wantInt:   IntFrom(0),
			wantErr:   false,
		},
		{
			name:      "positive number",
			inputJSON: []byte(`123`),
			wantInt:   IntFrom(123),
			wantErr:   false,
		},
		{
			name:      "negative number",
			inputJSON: []byte(`-42`),
			wantInt:   IntFrom(-42),
			wantErr:   false,
		},
		{
			name:      "null",
			inputJSON: []byte(`null`),
			wantInt:   NewInt(0, false),
			wantErr:   false,
		},
		{
			name:      "empty string",
			inputJSON: []byte(`""`),
			wantInt:   NewInt(0, false),
			wantErr:   false,
		},
		{
			name:      "number as string",
			inputJSON: []byte(`"123"`),
			wantInt:   NewInt(0, false), // current implementation does not support this
			wantErr:   true,
		},
		{
			name:      "invalid value",
			inputJSON: []byte(`"invalid"`),
			wantInt:   NewInt(0, false),
			wantErr:   true,
		},
		{
			name:      "invalid json",
			inputJSON: []byte(`invalid`),
			wantInt:   NewInt(0, false),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i Int
			err := json.Unmarshal(tt.inputJSON, &i)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantInt.Valid, i.Valid)
			assert.Equal(t, tt.wantInt.Int64, i.Int64)
		})
	}
}

func TestInt_EncodeValues(t *testing.T) {
	key := "test_int"

	// Test with a valid Int
	t.Run("valid", func(t *testing.T) {
		v := &url.Values{}
		i := IntFrom(123)
		err := i.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.Equal(t, "123", v.Get(key))
	})

	// Test with an invalid Int
	t.Run("invalid", func(t *testing.T) {
		v := &url.Values{}
		i := NewInt(0, false)
		err := i.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.False(t, v.Has(key))
	})
}

func TestInt_ValueOrZero(t *testing.T) {
	// 测试有效值
	validInt := IntFrom(42)
	assert.Equal(t, int64(42), validInt.ValueOrZero())

	// 测试无效值
	invalidInt := NewInt(0, false)
	assert.Equal(t, int64(0), invalidInt.ValueOrZero())
}

func TestInt_GobEncoding(t *testing.T) {
	// 测试有效的 Int 值
	var validBuf bytes.Buffer
	validInt := IntFrom(42)
	err := gob.NewEncoder(&validBuf).Encode(validInt)
	assert.NoError(t, err)

	var decodedValidInt Int
	err = gob.NewDecoder(&validBuf).Decode(&decodedValidInt)
	assert.NoError(t, err)
	assert.Equal(t, validInt, decodedValidInt)

	// 测试无效的 Int 值
	var invalidBuf bytes.Buffer
	invalidInt := NewInt(0, false)
	err = gob.NewEncoder(&invalidBuf).Encode(invalidInt)
	assert.NoError(t, err)

	var decodedInvalidInt Int
	err = gob.NewDecoder(&invalidBuf).Decode(&decodedInvalidInt)
	assert.NoError(t, err)
	assert.Equal(t, invalidInt, decodedInvalidInt)
}

func TestInt_NullValue(t *testing.T) {
	// 验证有效 Int 的 NullValue
	validInt := NewInt(123, true)
	assert.True(t, validInt.NullValue().Valid)
	assert.Equal(t, int64(123), validInt.NullValue().Int64)

	// 验证无效 Int 的 NullValue
	invalidInt := NewInt(0, false)
	assert.False(t, invalidInt.NullValue().Valid)
	assert.Equal(t, int64(0), invalidInt.NullValue().Int64)
}
