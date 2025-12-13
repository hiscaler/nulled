# Null<font color="red">E</font><font color="green">R</font>

`nulled` is a Go library that provides nullable types for `bool`, `float64`, `int64`, `string`, and `time.Time`. It is built as a wrapper around the popular `gopkg.in/guregu/null.v4` library, extending it with additional functionality and convenience methods.

This library is particularly useful for handling nullable values from databases, JSON/YAML APIs, or any other source where a value can be either present or null.

## Features

- Provides nullable `Bool`, `Float`, `Int`, `String`, and `Time` types.
- Seamlessly handles JSON encoding and decoding (`json.Marshaler`, `json.Unmarshaler`).
- Supports text encoding and decoding (`encoding.TextUnmarshaler`).
- Supports Gob encoding and decoding (`gob.GobEncoder`, `gob.GobDecoder`).
- Helper functions to create nullable types from values and pointers.
- Methods to safely retrieve values or a zero-value if null (`ValueOrZero`).
- `EncodeValues` method for encoding values into `net/url.Values`.

## Installation

To install the `nulled` package, use `go get`:

```bash
go get github.com/hiscaler/nulled
```

## Usage

Here are some examples of how to use the different nullable types provided by this library.

### String

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
)

func main() {
	// Valid string
	s1 := nulled.StringFrom("Hello, World!")
	fmt.Printf("s1: Value=%q, Valid=%t\n", s1.ValueOrZero(), s1.Valid)

	// Null string (from empty string)
	s2 := nulled.StringFrom("")
	fmt.Printf("s2: Value=%q, Valid=%t\n", s2.ValueOrZero(), s2.Valid)

	// JSON Marshalling
	jsonData, _ := json.Marshal(struct {
		S1 nulled.String
		S2 nulled.String
	}{S1: s1, S2: s2})
	fmt.Println("JSON:", string(jsonData)) // JSON: {"S1":"Hello, World!","S2":null}
}
```

### Int

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
)

func main() {
	// Valid int
	i1 := nulled.IntFrom(123)
	fmt.Printf("i1: Value=%d, Valid=%t\n", i1.ValueOrZero(), i1.Valid)

	// Null int (from pointer)
	var p *int64
	i2 := nulled.IntFromPtr(p)
	fmt.Printf("i2: Value=%d, Valid=%t\n", i2.ValueOrZero(), i2.Valid)

	// JSON Marshalling
	jsonData, _ := json.Marshal(struct {
		I1 nulled.Int
		I2 nulled.Int
	}{I1: i1, I2: i2})
	fmt.Println("JSON:", string(jsonData)) // JSON: {"I1":123,"I2":null}
}
```

### Float

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
)

func main() {
	// Valid float
	f1 := nulled.FloatFrom(123.45)
	fmt.Printf("f1: Value=%f, Valid=%t\n", f1.ValueOrZero(), f1.Valid)

	// Null float
	f2 := nulled.NewFloat(0, false)
	fmt.Printf("f2: Value=%f, Valid=%t\n", f2.ValueOrZero(), f2.Valid)

	// JSON Marshalling
	jsonData, _ := json.Marshal(struct {
		F1 nulled.Float
		F2 nulled.Float
	}{F1: f1, F2: f2})
	fmt.Println("JSON:", string(jsonData)) // JSON: {"F1":123.45,"F2":null}
}
```

### Bool

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
)

func main() {
	// Valid bool
	b1 := nulled.BoolFrom(true)
	fmt.Printf("b1: Value=%t, Valid=%t\n", b1.ValueOrZero(), b1.Valid)

	// Null bool
	b2 := nulled.NewBool(false, false)
	fmt.Printf("b2: Value=%t, Valid=%t\n", b2.ValueOrZero(), b2.Valid)

	// JSON Marshalling
	jsonData, _ := json.Marshal(struct {
		B1 nulled.Bool
		B2 nulled.Bool
	}{B1: b1, B2: b2})
	fmt.Println("JSON:", string(jsonData)) // JSON: {"B1":true,"B2":null}
}
```

### Time

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
	"time"
)

func main() {
	// Valid time
	t1 := nulled.TimeFrom(time.Now())
	fmt.Printf("t1: Value=%s, Valid=%t\n", t1.ValueOrZero(), t1.Valid)

	// Null time
	var p *time.Time
	t2 := nulled.TimeFromPtr(p)
	fmt.Printf("t2: Value=%s, Valid=%t\n", t2.ValueOrZero(), t2.Valid)

	// JSON Marshalling
	jsonData, _ := json.Marshal(struct {
		T1 nulled.Time
		T2 nulled.Time
	}{T1: t1, T2: t2})
	fmt.Println("JSON:", string(jsonData))
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.