package conver

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"strconv"
)

// IntToString int转string
func IntToString(num int) string {
	string := strconv.Itoa(num)
	return string
}

// Int64ToString int64转string
func Int64ToString(num int64) string {
	string := strconv.FormatInt(num, 10)
	return string
}

// StringToInt string转int
func StringToInt(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}

// StringToInt64 string转int64
func StringToInt64(str string) int64 {
	num, _ := strconv.ParseInt(str, 10, 64)
	return num
}

// Utf8ToGbk utf8转Gbk
func Utf8ToGbk(text string) string {
	r := bytes.NewReader([]byte(text))
	decoder := transform.NewReader(r, simplifiedchinese.GBK.NewDecoder()) //GB18030
	content, _ := io.ReadAll(decoder)
	return string(content)
}

// GbkToUtf8 Gbk转Utf8
func GbkToUtf8(b []byte) []byte {
	tfr := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(tfr)
	if e != nil {
		return nil
	}
	return d
}
