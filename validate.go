package validate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidResult : バリデーションの復帰値
type ValidResult struct {
	message string
}

// Message : 出力エラーメッセージを設定する
func (result *ValidResult) Message(message string) {
	result.message = message
}

// Error : エラーメッセージを出力する
func (result *ValidResult) Error() string {
	return result.message
}

// Result : バリデーションの復帰値を取り扱うインタフェース
type Result interface {
	Message(string)
	Error() string
}

// Validate : 単一の値を検証するバリデート構造体
type Validate struct {
	Value  interface{}
	result *ValidResult
}

// Result : 復帰値情報を返却する
func (validate *Validate) Result() Result { return validate.result }

func (validate *Validate) confirm(fn func(value string) bool) {
	value := reflect.ValueOf(validate.Value)
	if value.Kind() == reflect.Slice {
		// スライスの場合、全要素を検証する
		for i := 0; i < value.Len(); i++ {
			if fn(fmt.Sprint(value.Index(i).Interface())) == false {
				break
			}
		}
	} else {
		// スライス以外の場合、要素を検証する
		fn(fmt.Sprint(validate.Value))
	}
}

// Require : 空文字列か確認する
func (validate *Validate) Require() Result {
	validate.confirm(func(value string) bool {
		// 空文字列の場合、エラーとして扱う
		if value == "" {
			validate.Result().Message("It is a required input item")
			return false
		}
		return true
	})
	return validate.Result()
}

// MaxLen : 文字数の長さ制限
func (validate *Validate) MaxLen(max int) Result {
	validate.confirm(func(value string) bool {
		// 文字列が、指定した数より長い場合エラーとして扱う
		length := len(value)
		if length > max {
			validate.Result().Message(fmt.Sprintf("String too long. %s(%d) > max(%d)", value, length, max))
			return false
		}
		return true
	})
	return validate.Result()
}

// MinLen : 文字数の短さ制限
func (validate *Validate) MinLen(min int) Result {
	validate.confirm(func(value string) bool {
		// 文字列が、指定した数より短い場合エラーとして扱う
		length := len(fmt.Sprint(value))
		if length < min {
			validate.Result().Message(fmt.Sprintf("String too short. %s(%d) < min(%d)", value, length, min))
			return false
		}
		return true
	})
	return validate.Result()
}

// EMail : E-MAIL アドレスのチェック
func (validate *Validate) EMail() Result {
	var result = validate.Result()
	validate.confirm(func(email string) bool {
		// 空文字列の場合はスルーする
		if email == "" {
			return true
		}
		// E-MAIL アドレスを @ で分割し、ローカルパートとドメイン部分を分ける
		address := strings.Split(email, "@")
		// 2つ以外に分割された場合、またはローカルパート、ドメイン部分が空の場合、メールアドレスが正しくないためエラーとする
		if len(address) != 2 || address[0] == "" || address[1] == "" {
			result.Message("E-MAIL address is wrong")
			return false
		}
		localpart, domain := address[0], address[1]
		// ローカルパートがascii文字だけで構成されているか確認する
		for _, v := range []byte(localpart) {
			if v > 127 {
				result.Message("E-MAIL address is wrong")
				return false
			}
		}
		// ドメイン部分がascii文字だけで構成されているか確認する
		for _, v := range []byte(domain) {
			if v > 127 {
				result.Message("E-MAIL address is wrong")
				return false
			}
		}
		return true
	})
	return result
}

// Number : 数字か否かを確認する
func (validate *Validate) Number() Result {
	validate.confirm(func(value string) bool {
		// 数値変換可能かチェック
		_, err := strconv.Atoi(value)
		if err != nil {
			validate.Result().Message(err.Error())
			return false
		}
		return true
	})
	return validate.Result()
}

// Max : 指定した数字が、max最大値を超過しているか確認
func (validate *Validate) Max(max int) Result {
	validate.confirm(func(value string) bool {
		// 数値変換可能かチェック
		n, err := strconv.Atoi(value)
		if err != nil {
			validate.Result().Message(err.Error())
			return false
		}
		if n > max {
			validate.Result().Message(fmt.Sprintf("Exceeds the maximum value. %d > %d", n, max))
			return false
		}
		return true
	})
	return validate.Result()
}

// Min : 指定した数字が、min最小値を下回っているか確認
func (validate *Validate) Min(min int) Result {
	validate.confirm(func(value string) bool {
		// 数値変換可能かチェック
		n, err := strconv.Atoi(value)
		if err != nil {
			validate.Result().Message(err.Error())
			return false
		}
		if n < min {
			validate.Result().Message(fmt.Sprintf("Exceeds the min value. %d < %d", n, min))
			return false
		}
		return true
	})
	return validate.Result()
}

// Date : 日付か否かを確認する
func (validate *Validate) Date() Result {
	validate.confirm(func(value string) bool {
		// 2019-01-02, 2018/01/02, 2018/1/2, 2018/2/02 を許可する
		date := strings.Replace(value, "/", "-", -1)
		match, err := regexp.MatchString(`^\d{4}[\-\/](\d{2}|\d)[\-\/](\d{2}|\d)$`, date)
		if !match || err != nil {
			validate.Result().Message("Not date")
			return false
		}
		if _, err := time.Parse("2006-01-02", date); err != nil {
			validate.Result().Message(err.Error())
			return false
		}
		return true
	})
	return validate.Result()
}

// URL : 指定された文字列がURLか確認する
func (validate *Validate) URL() Result {
	validate.confirm(func(value string) bool {
		if value == "" {
			return true
		}
		match, err := regexp.MatchString(`^https?://[\w/:%#\$&\?\(\)~\.=\+\-]+$`, value)
		if !match || err != nil {
			validate.Result().Message("Not URL")
			return false
		}
		return true
	})
	return validate.Result()
}

// Match : 正規表現とマッチしているか確認する
func (validate *Validate) Match(regex string) Result {
	// 正規表現オブジェクトを生成する
	r, err := regexp.Compile(regex)
	if err != nil {
		validate.Result().Message(err.Error())
		return validate.Result()
	}

	validate.confirm(func(value string) bool {
		if value == "" {
			return true
		}
		if r.MatchString(value) == false {
			validate.Result().Message(regex + " no match")
			return false
		}
		return true
	})
	return validate.Result()
}

// HasErrors : バリデート時にエラーが発生しているか確認する
func (validate *Validate) HasErrors() bool {
	return validate.Result().Error() != ""
}

// エラーメッセージを返却する
func (validate *Validate) Error() string {
	return validate.Result().Error()
}

// Validates : 複数の値を検証するバリデート構造体
type Validates struct {
	Values   map[string]interface{}
	Errors   map[string]ValidError
	Validate *Validate
	Keyname  string
}

func (validates *Validates) callValidate(keyname, typename string, values ...interface{}) Result {
	// 指定したキーの値検証を行う
	if v, ok := validates.Values[keyname]; ok {
		var result Result
		switch typename {
		case "Require":
			result = Validator(v).Require()
		case "MaxLen":
			result = Validator(v).MaxLen(values[0].(int))
		case "MinLen":
			result = Validator(v).MinLen(values[0].(int))
		case "EMail":
			result = Validator(v).EMail()
		case "Min":
			result = Validator(v).Min(values[0].(int))
		case "Max":
			result = Validator(v).Max(values[0].(int))
		case "Number":
			result = Validator(v).Number()
		case "Date":
			result = Validator(v).Date()
		case "URL":
			result = Validator(v).URL()
		case "Match":
			result = Validator(v).Match(values[0].(string))
		}
		if result.Error() != "" {
			validates.Errors[keyname] = result
		}
		return result
	}
	// 指定したキーが見つからない場合、エラーとする
	result := &ValidResult{}
	result.Message(fmt.Sprintf("%s: not found", keyname))
	validates.Errors[keyname] = result
	return result
}

// Require : 空文字列か確認する
func (validates *Validates) Require(keyname string) Result {
	return validates.callValidate(keyname, "Require")
}

// MaxLen : 文字数の長さ制限
func (validates *Validates) MaxLen(keyname string, max int) Result {
	return validates.callValidate(keyname, "MaxLen", max)
}

// MinLen : 文字数の短さ制限
func (validates *Validates) MinLen(keyname string, min int) Result {
	return validates.callValidate(keyname, "MinLen", min)
}

// EMail : E-MAIL アドレスのチェック
func (validates *Validates) EMail(keyname string) Result {
	return validates.callValidate(keyname, "EMail")
}

// Number : 数字か否かを確認する
func (validates *Validates) Number(keyname string) Result {
	return validates.callValidate(keyname, "Number")
}

// Max : 指定した数字が、max最大値を超過しているか確認
func (validates *Validates) Max(keyname string, max int) Result {
	return validates.callValidate(keyname, "Max", max)
}

// Min : 指定した数字が、min最小値を下回っているか確認
func (validates *Validates) Min(keyname string, min int) Result {
	return validates.callValidate(keyname, "Min", min)
}

// Date : 日付か否かを確認する
func (validates *Validates) Date(keyname string) Result {
	return validates.callValidate(keyname, "Date")
}

// URL : 指定された文字列がURLか確認する
func (validates *Validates) URL(keyname string) Result {
	return validates.callValidate(keyname, "URL")
}

// Match : 正規表現とマッチしているか確認する
func (validates *Validates) Match(keyname, regexp string) Result {
	return validates.callValidate(keyname, "Match", regexp)
}

// HasErrors : バリデート時にエラーが発生しているか確認する
func (validates *Validates) HasErrors() bool {
	if len(validates.ErrList()) > 0 {
		return true
	}
	return false
}

// ErrList : エラー内容一覧を返却する
func (validates *Validates) ErrList() []error {
	var errors []error
	for _, v := range validates.Errors {
		errors = append(errors, v)
	}
	return errors
}

// ValidError : Validates で値検証時に、存在しないキーを指定された場合のエラー型
type ValidError error

// Validator : 単一のバリデーションチェックに使用する
func Validator(value interface{}) *Validate {
	return &Validate{
		Value:  value,
		result: &ValidResult{},
	}
}

// Validators : 複数のバリデーションチェックに使用する
func Validators(i interface{}) (*Validates, error) {
	if i == nil {
		return nil, fmt.Errorf("marshal error")
	}
	// 渡されたデータを一旦JSON文字列へ変換
	buf, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	// JSON文字列をmap[string]interface{}へ変換する
	var values map[string]interface{}
	if err := json.Unmarshal(buf, &values); err != nil {
		return nil, err
	}
	// 複数のバリデーションを行うための構造体を返却する
	return &Validates{
		Values: values,
		Errors: make(map[string]ValidError),
	}, nil
}
