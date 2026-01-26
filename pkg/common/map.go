package common

import (
	"encoding/json"

	"github.com/kaichao/gopkg/errors"
)

// ToMap 将 JSON 字符串转换为 map[string]interface{}
func ToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, errors.WrapE(err, "json.Unmarshal failed", "json-text", jsonStr)
	}
	return result, nil
}

// ToJSON 将 map[string]interface{} 转换为 JSON 字符串
func ToJSON(m map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return "", errors.WrapE(err, "json.Unmarshal failed", "map", m)
	}
	return string(jsonBytes), nil
}

// ToMapStringInt 将 JSON 字符串转换为 map[string]int
func ToMapStringInt(jsonStr string) (map[string]int, error) {
	var result map[string]int
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, errors.WrapE(err, "json.Unmarshal failed", "json-text", jsonStr)
	}
	return result, nil
}

// ToMapStringString 将 JSON 字符串转换为 map[string]string
func ToMapStringString(jsonStr string) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, errors.WrapE(err, "json.Unmarshal failed", "json-text", jsonStr)
	}
	return result, nil
}
