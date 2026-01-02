package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testTime     = time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)
	testTimeStr  = `"2023-10-27T10:00:00Z"`
	testTimeText = "2023-10-27T10:00:00Z"
)

func TestTime_NewTime(t *testing.T) {
	valid := NewTime(testTime, true)
	assert.True(t, valid.Valid)
	assert.Equal(t, testTime, valid.Time)

	invalid := NewTime(time.Time{}, false)
	assert.False(t, invalid.Valid)
}

func TestTime_TimeFrom(t *testing.T) {
	valid := TimeFrom(testTime)
	assert.True(t, valid.Valid)
	assert.Equal(t, testTime, valid.Time)

	invalid := TimeFrom(time.Time{})
	assert.False(t, invalid.Valid)
}

func TestTime_TimeFromPtr(t *testing.T) {
	valid := TimeFromPtr(&testTime)
	assert.True(t, valid.Valid)
	assert.Equal(t, testTime, valid.Time)

	invalid := TimeFromPtr(nil)
	assert.False(t, invalid.Valid)
}

func TestTime_MarshalJSON(t *testing.T) {
	valid := NewTime(testTime, true)
	data, err := json.Marshal(valid)
	assert.NoError(t, err)
	assert.Equal(t, []byte(testTimeStr), data)

	invalid := NewTime(time.Time{}, false)
	data, err = json.Marshal(invalid)
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), data)
}

func TestTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		inputJSON []byte
		want      Time
		wantErr   bool
	}{
		{name: "valid time", inputJSON: []byte(testTimeStr), want: NewTime(testTime, true)},
		{name: "null", inputJSON: []byte(`null`), want: NewTime(time.Time{}, false)},
		{name: "invalid json", inputJSON: []byte(`"not a time"`), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ti Time
			err := json.Unmarshal(tt.inputJSON, &ti)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Valid, ti.Valid)
				// comparing time objects directly can be tricky, using Equal is fine for this
				assert.Equal(t, tt.want.Time, ti.Time)
			}
		})
	}
}

func TestTime_UnmarshalText(t *testing.T) {
	tests := []struct {
		name      string
		inputText []byte
		want      Time
		wantErr   bool
	}{
		{name: "valid RFC3339", inputText: []byte(testTimeText), want: NewTime(testTime, true)},
		{name: "empty bytes", inputText: []byte(``), want: NewTime(time.Time{}, false)},
		{name: "null string", inputText: []byte(`null`), want: NewTime(time.Time{}, false)},
		{name: "invalid string", inputText: []byte(`abc`), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ti Time
			err := ti.UnmarshalText(tt.inputText)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Valid, ti.Valid)
				assert.True(t, tt.want.Time.Equal(ti.Time))
			}
		})
	}
}

func TestTime_EncodeValues(t *testing.T) {
	key := "test_time"

	t.Run("valid", func(t *testing.T) {
		v := &url.Values{}
		ti := NewTime(testTime, true)
		err := ti.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.Equal(t, "2023-10-27 10:00:00", v.Get(key))
	})

	t.Run("invalid", func(t *testing.T) {
		v := &url.Values{}
		ti := NewTime(time.Time{}, false)
		err := ti.EncodeValues(key, v)
		assert.NoError(t, err)
		assert.False(t, v.Has(key))
	})
}

func TestTime_ValueOrZero(t *testing.T) {
	valid := NewTime(testTime, true)
	assert.Equal(t, testTime, valid.ValueOrZero())

	invalid := NewTime(time.Time{}, false)
	assert.True(t, invalid.ValueOrZero().IsZero())
}

func TestTime_GobEncoding(t *testing.T) {
	// Test valid Time
	var validBuf bytes.Buffer
	valid := NewTime(testTime, true)
	err := gob.NewEncoder(&validBuf).Encode(valid)
	assert.NoError(t, err)

	var decodedValid Time
	err = gob.NewDecoder(&validBuf).Decode(&decodedValid)
	assert.NoError(t, err)
	assert.Equal(t, valid, decodedValid)

	// Test invalid Time
	var invalidBuf bytes.Buffer
	invalid := NewTime(time.Time{}, false)
	err = gob.NewEncoder(&invalidBuf).Encode(invalid)
	assert.NoError(t, err)

	var decodedInvalid Time
	err = gob.NewDecoder(&invalidBuf).Decode(&decodedInvalid)
	assert.NoError(t, err)
	assert.Equal(t, invalid, decodedInvalid)
}

func TestTime_NullValue(t *testing.T) {
	// 验证有效 Time 的 NullValue
	validTime := NewTime(testTime, true)
	assert.True(t, validTime.NullValue().Valid)
	assert.Equal(t, testTime, validTime.NullValue().Time)

	// 验证无效 Time 的 NullValue
	invalidTime := NewTime(time.Time{}, false)
	assert.False(t, invalidTime.NullValue().Valid)
	assert.True(t, invalidTime.NullValue().Time.IsZero())
}
