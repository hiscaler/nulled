# Null<font color="red">E</font><font color="green">D</font>

空值编码/解码

`nulled` 是一个 Go 库，为 `bool`, `float64`, `int64`, `string` 和 `time.Time` 提供了可空类型。它是在流行的 `gopkg.in/guregu/null.v4` 库的基础上构建的，通过额外的功能和便利方法对其进行了扩展。

该库对于处理来自数据库、JSON/YAML API 或任何其他值可能存在或为空的来源的可空值特别有用。

## 特性

-   提供可空的 `Bool`, `Float`, `Int`, `String`, 和 `Time` 类型。
-   无缝处理 JSON 编码和解码 (`json.Marshaler`, `json.Unmarshaler`)。
-   实现 `sql.Scanner` 和 `driver.Valuer`，便于数据库集成。
-   支持文本编码和解码 (`encoding.TextUnmarshaler`)。
-   支持 Gob 编码和解码 (`gob.GobEncoder`, `gob.GobDecoder`)。
-   用于从值和指针创建可空类型的辅助函数。
-   安全检索值的方法，如果为空则返回零值 (`ValueOrZero`)。
-   用于将值编码为 `net/url.Values` 的 `EncodeValues` 方法。
-   `NullValue` 方法可获取底层的 `gopkg.in/guregu/null.v4` 类型。

## 安装

要安装 `nulled` 包，请使用 `go get`:

```bash
go get github.com/hiscaler/nulled
```

## 用法

以下是如何使用此库提供的不同可空类型的一些示例。

### String (字符串)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
)

func main() {
	// 有效字符串
	s1 := nulled.StringFrom("Hello, World!")
	fmt.Printf("s1: Value=%q, Valid=%t\n", s1.ValueOrZero(), s1.Valid)

	// 空字符串 (来自空字符串)
	s2 := nulled.StringFrom("")
	fmt.Printf("s2: Value=%q, Valid=%t\n", s2.ValueOrZero(), s2.Valid)

	// JSON 编码
	jsonData, _ := json.Marshal(struct {
		S1 nulled.String
		S2 nulled.String
	}{S1: s1, S2: s2})
	fmt.Println("JSON:", string(jsonData)) // JSON: {"S1":"Hello, World!","S2":null}
}
```

### Int (整数)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
)

func main() {
	// 有效整数
	i1 := nulled.IntFrom(123)
	fmt.Printf("i1: Value=%d, Valid=%t\n", i1.ValueOrZero(), i1.Valid)

	// 空整数 (来自指针)
	var p *int64
	i2 := nulled.IntFromPtr(p)
	fmt.Printf("i2: Value=%d, Valid=%t\n", i2.ValueOrZero(), i2.Valid)

	// JSON 编码
	jsonData, _ := json.Marshal(struct {
		I1 nulled.Int
		I2 nulled.Int
	}{I1: i1, I2: i2})
	fmt.Println("JSON:", string(jsonData)) // JSON: {"I1":123,"I2":null}
}
```

### Float (浮点数)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
)

func main() {
	// 有效浮点数
	f1 := nulled.FloatFrom(123.45)
	fmt.Printf("f1: Value=%f, Valid=%t\n", f1.ValueOrZero(), f1.Valid)

	// 空浮点数
	f2 := nulled.NewFloat(0, false)
	fmt.Printf("f2: Value=%f, Valid=%t\n", f2.ValueOrZero(), f2.Valid)

	// JSON 编码
	jsonData, _ := json.Marshal(struct {
		F1 nulled.Float
		F2 nulled.Float
	}{F1: f1, F2: f2})
	fmt.Println("JSON:", string(jsonData)) // JSON: {"F1":123.45,"F2":null}
}
```

### Bool (布尔值)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
)

func main() {
	// 有效布尔值
	b1 := nulled.BoolFrom(true)
	fmt.Printf("b1: Value=%t, Valid=%t\n", b1.ValueOrZero(), b1.Valid)

	// 空布尔值
	b2 := nulled.NewBool(false, false)
	fmt.Printf("b2: Value=%t, Valid=%t\n", b2.ValueOrZero(), b2.Valid)

	// JSON 编码
	jsonData, _ := json.Marshal(struct {
		B1 nulled.Bool
		B2 nulled.Bool
	}{B1: b1, B2: b2})
	fmt.Println("JSON:", string(jsonData)) // JSON: {"B1":true,"B2":null}
}
```

### Time (时间)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/nulled"
	"time"
)

func main() {
	// 有效时间
	t1 := nulled.TimeFrom(time.Now())
	fmt.Printf("t1: Value=%s, Valid=%t\n", t1.ValueOrZero(), t1.Valid)

	// 空时间
	var p *time.Time
	t2 := nulled.TimeFromPtr(p)
	fmt.Printf("t2: Value=%s, Valid=%t\n", t2.ValueOrZero(), t2.Valid)

	// JSON 编码
	jsonData, _ := json.Marshal(struct {
		T1 nulled.Time
		T2 nulled.Time
	}{T1: t1, T2: t2})
	fmt.Println("JSON:", string(jsonData))
}
```

## 数据库集成

所有 `nulled` 类型都实现了 `database/sql.Scanner` 和 `database/sql/driver.Valuer` 接口，使它们可以与 `database/sql` 开箱即用。

### 扫描空值

当您从数据库中扫描一个 `NULL` 值时，`nulled` 类型将被标记为无效 (`Valid: false`)。

```go
var name nulled.String
err := db.QueryRow("SELECT name FROM users WHERE id = 1").Scan(&name)
// 如果 'name' 列为 NULL，name.Valid 将为 false。
```

### 存储空值

要将 `NULL` 值存储到数据库中，您可以使用一个无效的 `nulled` 类型。

```go
invalidName := nulled.NewString("", false)
_, err := db.Exec("UPDATE users SET name = ? WHERE id = 1", invalidName)
// 这会将 'name' 列设置为 NULL。
```

## gopkg.in/guregu/null.v4 兼容性

为了保持与底层的 `gopkg.in/guregu/null.v4` 库的兼容性，每个 `nulled` 类型都有一个 `NullValue()` 方法，该方法返回相应的 `null.v4` 类型。

```go
import "gopkg.in/guregu/null.v4"

// 获取底层的 null.v4 类型
nulledString := nulled.StringFrom("hello")
nullV4String := nulledString.NullValue() // 这是一个 null.String

// 您可以将其与期望 null.v4 类型的函数一起使用
func processNullV4(s null.String) {
    // ...
}
processNullV4(nullV4String)
```

## 与 go-querystring/query 集成

`nulled` 类型可以与 `github.com/google/go-querystring/query` 无缝集成，以将结构体编码为 URL 查询参数。这是通过在每个 `nulled` 类型上实现的 `EncodeValues` 方法来促成的，该方法能正确处理有效值和空值。

首先，请确保您已安装 `go-querystring`：

```bash
go get github.com/google/go-querystring/query
```

这是一个例子：

```go
package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/hiscaler/nulled"
)

type MyQueryParams struct {
	Name    nulled.String `url:"name"`
	Age     nulled.Int    `url:"age"`
	Active  nulled.Bool   `url:"active"`
	Amount  nulled.Float  `url:"amount"`
	Created nulled.Time   `url:"created"`
	Search  nulled.String `url:"search,omitempty"` // omitempty 如果为空或 nil 则会跳过
}

func main() {
	params := MyQueryParams{
		Name:    nulled.StringFrom("John Doe"),
		Age:     nulled.IntFrom(30),
		Active:  nulled.BoolFrom(true),
		Amount:  nulled.FloatFrom(123.45),
		Created: nulled.TimeFrom(time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)),
		Search:  nulled.StringFrom("Go Lang"),
	}

	v, _ := query.Values(params)
	fmt.Println("包含所有有效值的查询:", v.Encode())
	// 输出: 包含所有有效值的查询: active=true&age=30&amount=123.45&created=2023-01-15+10%3A30%3A00&name=John+Doe&search=Go+Lang

	nullParams := MyQueryParams{
		Name:    nulled.NewString("", false), // 空字符串
		Age:     nulled.NewInt(0, false),    // 空整数
		Active:  nulled.NewBool(false, false), // 空布尔值
		Amount:  nulled.NewFloat(0, false),  // 空浮点数
		Created: nulled.NewTime(time.Time{}, false), // 空时间
		Search:  nulled.NewString("", false), // 带 omitempty 的空字符串
	}

	nv, _ := query.Values(nullParams)
	fmt.Println("包含空值的查询:", nv.Encode())
	// 输出: 包含空值的查询:
}
```

## 贡献

欢迎贡献！请随时提交拉取请求或开启一个 issue。

## 许可证

该项目根据 MIT 许可证授权 - 有关详细信息，请参阅 [LICENSE](LICENSE) 文件。
