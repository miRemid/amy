package utils

import (
	"encoding/json"
	"bytes"
)

// LoadIntoMap return a map[string]interface{}
func LoadIntoMap(data []byte) (res map[string]interface{}) {
	decode := json.NewDecoder(bytes.NewReader(data))
	decode.UseNumber()
	decode.Decode(&res)
	return
}

// LoadIntoStruct return a error
func LoadIntoStruct(data []byte, v interface{}) error{
	decode := json.NewDecoder(bytes.NewReader(data))
	decode.UseNumber()
	return decode.Decode(&v)
}