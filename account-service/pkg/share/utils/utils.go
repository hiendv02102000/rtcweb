package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const time_layout = "2006-01-02 15:04:05"

func EncryptPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(h[:])
}
func FormatTime(time time.Time) string {
	return time.Format(time_layout)
}
func FormatStringSpace(s string) string {

	words := strings.Fields(s)
	res := ""
	for _, word := range words {
		if len(word) == 0 {
			continue
		}
		res += word + " "
	}
	res = strings.TrimSpace(res)
	return res
}
func ConverInterfaceToString(vals []interface{}) ([]string, error) {
	rs := make([]string, len(vals))
	for i, v := range vals {
		val, ok := v.(string)
		if !ok {
			return nil, errors.New("this is not string array")
		}
		rs[i] = val
	}
	return rs, nil
}
func RemoveEndFile(s string) string {
	res := Reverse(s)
	res = strings.Replace(res, ".", "", 1)
	res = Reverse(res)
	return res
}
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func ConvertMapToObject(valMap map[string]interface{}, object interface{}) error {
	jsonStr, err := json.Marshal(valMap)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonStr, object)
	return err
}
