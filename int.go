package nulled

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"strconv"

	"gopkg.in/guregu/null.v4"
)

type Int null.Int

func NewInt(i int64, valid bool) Int {
	return Int(null.NewInt(i, valid))
}

func IntFrom(i int64) Int {
	return NewInt(i, true)
}

func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

func (i Int) ValueOrZero() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

func (i Int) EncodeValues(key string, v *url.Values) error {
	if !i.Valid {
		return nil
	}

	v.Set(key, strconv.FormatInt(i.Int64, 10))
	return nil
}

func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte(`null`), nil
	}

	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// GobEncode implements the gob.GobEncoder interface.
func (i Int) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i.Int64)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(i.Valid)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode implements the gob.GobDecoder interface.
func (i *Int) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&i.Int64)
	if err != nil {
		return err
	}
	return dec.Decode(&i.Valid)
}

func (i *Int) UnmarshalJSON(bytes []byte) error {
	// 通过类型转换获取值
	n := null.Int(*i)

	// 处理空字符串或 null 的情况
	if string(bytes) == `""` || string(bytes) == `null` || string(bytes) == `nil` {
		n.Valid = false
		n.Int64 = 0
		*i = Int(n)
		return nil
	}

	// 解析数值
	err := json.Unmarshal(bytes, &n.Int64)
	if err != nil {
		n.Valid = false
		n.Int64 = 0
		*i = Int(n)
		return err
	}

	n.Valid = true
	*i = Int(n)
	return nil
}
