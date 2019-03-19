package validate

import (
	"fmt"
	"testing"
)

func Test__Validate(t *testing.T) {
	validator := Validator("")

	validator.Require()
	validator.EMail()
	validator.URL()
	if validator.HasErrors() == false {
		t.Fatal("HasError()")
	}
	fmt.Println(validator)

	validator = Validator("漢字@ochipin.net")
	validator.EMail()
	if validator.HasErrors() == false {
		t.Fatal("HasError()")
	}

	validator = Validator("ascii.name@メールホスト")
	validator.EMail()
	if validator.HasErrors() == false {
		t.Fatal("HasError()")
	}
}

func Test__Validate2(t *testing.T) {
	validator := Validator("n")
	validator.Max(30)
	validator.Min(30)
	validator.Date()
	if validator.HasErrors() == false {
		t.Fatal("HasError()")
	}

	validator = Validator("name")
	validator.Match(`name[`)
	if validator.HasErrors() == false {
		t.Fatal("HasError()")
	}

	validator = Validator("name")
	validator.Match(`^name2$`)
	if validator.HasErrors() == false {
		t.Fatal("HasError()")
	}
}
func Test__Validates(t *testing.T) {
	validators, err := Validators(map[string]interface{}{
		"username": "User Name",
		"mail":     "ochipin@works.net",
		"url":      "https://ochipin.net",
		"date":     "2019-03-14",
		"min":      "0",
		"max":      "10",
		"number":   "999",
	})
	if err != nil {
		t.Fatal(err)
	}
	// ユーザ名は空ではないことを確認
	validators.Require("username")
	validators.MaxLen("username", 9)
	validators.MinLen("username", 0)
	validators.Date("date")
	validators.EMail("mail").Message("EMAIL ERROR")
	validators.Min("min", 0)
	validators.Max("max", 15)
	validators.URL("url")
	validators.Number("number")
	validators.Match("username", `^User Name$`)
	if validators.HasErrors() {
		t.Fatal("HasErrors() Error")
	}
}

func Test__ValidatesError1(t *testing.T) {
	// 入力データ
	validators, err := Validators(map[string]interface{}{
		"username": "",
		"minlen":   "minlen",
		"maxlen":   "maxlen",
		"mail":     "@ochipinworks.net",
		"url":      "htt://ochipin.net",
		"date":     "2019/99/14",
		"min":      "4",
		"max":      "10",
		"number":   "999s",
	})
	if err != nil {
		t.Fatal(err)
	}
	// ユーザ名は空ではないことを確認
	validators.Require("username")
	validators.Require("username")
	validators.MaxLen("maxlen", 5)
	validators.MaxLen("maxlen", 5)
	validators.MinLen("minlen", 7)
	validators.MinLen("minlen", 7)
	validators.Date("date")
	validators.Date("date")
	validators.EMail("mail")
	validators.EMail("mail")
	validators.Min("min", 5)
	validators.Min("min", 5)
	validators.Max("max", 9)
	validators.Max("max", 9)
	validators.URL("url")
	validators.URL("url")
	validators.Number("number")
	validators.Number("number")
	validators.Match("username", `^String$`)
	validators.Match("username", `^String$`)
	if validators.HasErrors() == false {
		t.Fatal("HasErrors() Error")
	}
	for k, v := range validators.Errors {
		fmt.Println(k, v)
	}
}

func Test__ValidatesError2(t *testing.T) {
	// 入力データ
	validators, err := Validators(map[string]interface{}{
		"username": "",
		"minlen":   "minlen",
		"maxlen":   "maxlen",
		"mail":     "@ochipinworks.net",
		"url":      "htt://ochipin.net",
		"date":     "2019/99/14",
		"min":      "4",
		"max":      "10",
		"number":   "999s",
	})
	if err != nil {
		t.Fatal(err)
	}
	// 存在しない入力データのキーを渡した場合は、エラーとなる
	validators.Require("username2").Message("NG")
	validators.MaxLen("maxlen2", 5)
	validators.MinLen("minlen2", 7)
	validators.Date("date2")
	validators.EMail("mail2")
	validators.Min("min2", 5)
	validators.Max("max2", 9)
	validators.URL("url2")
	validators.Number("number2")
	validators.Match("url3", `^String$`)
	if validators.HasErrors() == false {
		t.Fatal("HasErrors() Error")
	}
	for k, v := range validators.Errors {
		fmt.Println(k, v)
	}
}
func Test__ValidatesError3(t *testing.T) {
	if _, err := Validators(200); err == nil {
		t.Fatal("Validates Error")
	}
	if _, err := Validators(nil); err == nil {
		t.Fatal("Validates Error")
	}
}
