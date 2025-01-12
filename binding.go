package xin

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/fengjx/go-halo/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()
var validate *validator.Validate

func init() {
	decoder.IgnoreUnknownKeys(true)
	decoder.SetAliasTag("json")
	decoder.RegisterConverter([]string{}, func(s string) reflect.Value {
		return reflect.ValueOf(strings.Split(s, ","))
	})
	validate = validator.New()
	validate.SetTagName("binding")
}

// ShouldBind 从参数url参数和form表单解析参数
func ShouldBind(r *http.Request, obj any) error {
	values := r.URL.Query()
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			return err
		}
		for key, val := range r.Form {
			values[key] = val
		}
	}
	err := decoder.Decode(obj, values)
	if err != nil {
		return err
	}
	return validate.Struct(obj)
}

// ShouldBindJSON 从body解析json
func ShouldBindJSON(r *http.Request, obj any) error {
	err := json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		return err
	}
	return validate.Struct(obj)
}

// GetQuery 获取URL查询参数，如果参数不存在返回空字符串
func GetQuery(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetQueryDefault 获取URL查询参数，如果参数不存在返回默认值
func GetQueryDefault(r *http.Request, key, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetForm 获取表单参数，如果参数不存在返回空字符串
func GetForm(r *http.Request, key string) string {
	return r.FormValue(key)
}

// GetFormDefault 获取表单参数，如果参数不存在返回默认值
func GetFormDefault(r *http.Request, key, defaultValue string) string {
	_ = r.ParseForm()
	value := r.FormValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetHeader 获取请求头，如果不存在返回空字符串
func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

// GetHeaderDefault 获取请求头，如果不存在返回默认值
func GetHeaderDefault(r *http.Request, key, defaultValue string) string {
	value := r.Header.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetCookie 获取Cookie值，如果不存在返回空字符串和错误
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetCookieDefault 获取Cookie值，如果不存在返回默认值
func GetCookieDefault(r *http.Request, name, defaultValue string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return defaultValue
	}
	return cookie.Value
}
