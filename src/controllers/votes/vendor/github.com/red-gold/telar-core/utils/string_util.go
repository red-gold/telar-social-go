package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

// Hash convert txt to hash
func Hash(text string) ([]byte, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

// CompareHash if hash and plan are equal
func CompareHash(hash []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}

// GetInterfaceType get inteface type
func GetInterfaceType(input interface{}) string {
	s := fmt.Sprintf("%T", input)
	return GetType(s)
}

// GetInterfaceType get inteface type
func GetInterfaceTypeLower(input interface{}) string {
	return LowerFirst(GetInterfaceType(input))
}

// GetType get inteface type without package name
func GetType(s string) string {
	inType := ""
	if n := strings.IndexByte(s, '.'); n >= 0 {
		inType = s[(n + 1):]
	} else {
		inType = strings.Replace(s, "*", "", -1)
	}
	return inType
}

// GetTypeLowerFirst get inteface type without package name with lowercase first character
func GetTypeLowerFirst(s string) string {
	return LowerFirst(GetType(s))
}

// LowerFirst lower case the first character of string
func LowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

// GenerateDigits generate fix length of digits
func GenerateDigits(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// ParseHtmlTemplate parse html template return in string
func ParseHtmlTemplate(templateFilePath string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ParseHtmlBytesTemplate parse html template return in []byte
func ParseHtmlBytesTemplate(templateFilePath string, data interface{}) ([]byte, error) {
	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
