package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// b, err := utils.SendRequest("GET", "", "", nil)
	// fmt.Println(string(b))
	jsonData, err := json.Marshal(nil)
	fmt.Println(string(jsonData), err)
}
