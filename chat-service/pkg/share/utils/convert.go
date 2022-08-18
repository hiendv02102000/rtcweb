package utils

import "encoding/json"

func ConvertToObject(i interface{}, o any) error {
	d, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, o)
}
func ConvertObjectToJsonString(o interface{}) string {
	res := []byte{}
	err := json.Unmarshal(res, o)
	if err != nil {
		return "err"
	}
	return string(res)
}
