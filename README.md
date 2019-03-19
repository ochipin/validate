バリデーションライブラリ
===

使い方
---

### 単一の値を検証する
```go
package main

import (
	"fmt"

	"github.com/ochipin/validate"
)

func main() {
	// 検証したい値を Validator 関数に渡す
	v := validate.Validator("No number")
	// 数字ではない場合、"Not Number" エラーメッセージを出力する
	v.Number().Message("Not Number")
	// エラーが発生している場合は、Println でエラーメッセージ("Not Number") を出力する
	if v.HasErrors() {
		fmt.Println(v)
	}
}
```

### 複数の値を検証する

```go
package main

import (
	"fmt"

	"github.com/ochipin/validate"
)

func main() {
	// 検証したい複数のデータを Validators 関数に渡す(構造体も可)
	v, err := validate.Validators(map[string]interface{}{
		"username": "User Name",
		"age":      99,
		"start":    "2019-02-01",
	})
	if err != nil {
		panic(err)
	}
	// username は空文字列ではないことを検証
	v.Require("username") // .Message("...") 関数を使用しない場合、デフォルトのエラーメッセージとなる
	// "username" は、 User Name 文字列か検証する
	v.Match("username", `^User Name$`)
	// age には数字が格納されているか検証する
	v.Number("age")
	// 18-100 歳までOK
	v.Max("age", 80).Message("80歳までです")
	v.Min("age", 18)
	// 日付か否かを検証する
	v.Date("start")
	// エラーが発生している場合は表示する
	if v.HasErrors() {
		// ココのエラーでは、"age"のみがエラー出力される
		fmt.Println(v.ErrList()[0])
	}
}
```

バリデート可能一覧
---

* Number
* Max
* Min
* MaxLen
* MinLen
* URL
* Match
* EMail
* Require
* Date

[](validate_test.go) に詳細記載済み。

バリデート構造体をカスタムし、独自のバリデーション関数を追加する
---

```go
package main

import (
	"fmt"

	"github.com/ochipin/validate"
)

// Validate : カスタム用バリデーション構造体
type Validate struct {
	*validate.Validate
}

// Equal : 文字列マッチ用関数を追加
func (validate *Validate) Equal(equal string) validate.Result {
	// 文字列が一致していない場合、エラーとして扱う
	if fmt.Sprint(validate.Value) != equal {
		validate.Result().Message("Not match")
	}
	return validate.Result()
}

// Validates : 複数のバリデーションを行うカスタムバリデーション構造体
type Validates struct {
	*validate.Validates
}

// Equal : 文字列マッチ用関数を追加
func (validates *Validates) Equal(keyname, equal string) validate.Result {
	if v, ok := validates.Values[keyname]; ok {
		result := Validator(v).Equal(equal)
		if result.Error() != "" {
			validates.Errors[keyname] = result
		}
		return result
	}
	// 指定したキーが見つからない場合
	result := &validate.ValidResult{}
	result.Message(fmt.Sprintf("%s: not found", keyname))
	validates.Errors[keyname] = result
	return result
}

// Validator : 単一のバリデーションチェックに使用する
func Validator(value interface{}) *Validate {
	return &Validate{
		Validate: validate.Validator(value),
	}
}

// Validators : 複数のバリデーションチェックに使用する
func Validators(i interface{}) (*Validates, error) {
	v, err := validate.Validators(i)
	if err != nil {
		return nil, err
	}
	return &Validates{
		Validates: v,
	}, nil
}

func main() {
	v := Validator("User Name")
	v.Equal("User Name").Message("NO MATCH USER NAME")
	if v.HasErrors() {
		fmt.Println(v)
	}
}
```
