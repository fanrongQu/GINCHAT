package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// MD5字符串
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempString := h.Sum(nil)
	return strings.ToLower(hex.EncodeToString(tempString))
}

// MD5字符串
func Md5ScalEncode(data string, salt string) string {
	return Md5Encode(data + salt)
}

// MD5字符串
func ValidMd5Scal(data string, salt string, md5Str string) bool {
	return Md5Encode(data+salt) == md5Str
}
