package models

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
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

// 有効なメールアドレスチェック
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "無効なEmailアドレスです")
	}
}

// 新規のメールアドレスチェック
func (f *Form) NotSameEmail(field string, user User) {
	value := f.Get(field)
	if strings.TrimSpace(value) == user.Email {
		f.Errors.Add(field, "このメールアドレスは既に使われています。")
	}
}

// ログインチェック
func (f *Form) SameEmailAndPassword(field1 string, field2 string, user User) {
	value1 := f.Get(field1)
	value2 := Encrypt(f.Get(field2))
	if strings.TrimSpace(value1) != user.Email || strings.TrimSpace(value2) != user.Password {
		f.Errors.Add(field1, "メールアドレスかパスワードが違います。")
		f.Errors.Add(field2, "メールアドレスかパスワードが違います。")
	}
}

