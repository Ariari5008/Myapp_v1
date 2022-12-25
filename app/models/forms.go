package models

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

)

type TemplateData struct {
	Form Form 
	Data map[string]interface{}
}

type Form struct {
	url.Values
	Errors errors
}

// コンストラクタ
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// エラーの存在チェック
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// 必須項目
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "この項目は必須です")
		}
	}
}

// 最小文字数
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	str := utf8.RuneCountInString(x)
	if str < length {
		f.Errors.Add(field, fmt.Sprintf(
			"このフィールドの最小文字数は%d文字です", length))
		return false
	}
	return true
}

// 最大文字数
func (f *Form) MaxLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	str := utf8.RuneCountInString(x)
	if str > length {
		f.Errors.Add(field, fmt.Sprintf(
			"このフィールドの最大文字数は%d文字です", length))
		return false
	}
	return true
}


// 新規のメールアドレスチェック
func (f *Form) NotSameName(field string, user User) {
	value := f.Get(field)
	if strings.TrimSpace(value) == user.Name {
		f.Errors.Add(field, "この名前は既に使われています。")
	}
}

// ログインチェック
func (f *Form) SameNameAndPassword(field1 string, field2 string, user User) {
	value1 := f.Get(field1)
	value2 := Encrypt(f.Get(field2))
	if strings.TrimSpace(value1) != user.Name || strings.TrimSpace(value2) != user.Password {
		f.Errors.Add(field1, "名前かパスワードが違います。")
		f.Errors.Add(field2, "名前かパスワードが違います。")
	}
}

