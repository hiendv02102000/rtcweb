package utils

import "encoding/json"

func ConvertToPbject(i interface{}, o any) error {
	d, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, o)
}
