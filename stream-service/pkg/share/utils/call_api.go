package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const HOST_ACCOUNT_SERVICE = "http://localhost:8080"

//"GET" POST
func SendRequest(method string, url string, token string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		bearer := "Bearer " + token
		req.Header.Add("Authorization", bearer)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
