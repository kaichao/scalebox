package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kaichao/gopkg/errors"
)

// SetJSONAttribute ...
// 设置属性（值统一处理为字符串类型）
func SetJSONAttribute(jsonString string, attrName string, attrValue string) (string, error) {
	data := make(map[string]interface{})
	if jsonString == "" {
		jsonString = "{}"
	}
	// 解析原始JSON
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		// logrus.Errorf("json.Unmarshal(), err-info:%v", err)
		return "", errors.WrapE(err, "json.Unmarshal failed", "json-text", jsonString)
	}

	// 智能类型转换
	var value interface{} = attrValue
	if num, err := json.Number(attrValue).Int64(); err == nil {
		value = num
	} else if b, err := parseBool(attrValue); err == nil {
		value = b
	}

	data[attrName] = value
	// 生成新JSON
	result, err := json.Marshal(data)
	if err != nil {
		// logrus.Errorf("json.Marshal(), err-info:%v", err)
		return "", errors.WrapE(err, "json.Marshal failed", "json-data", data)
	}
	return string(result), nil
}

// RemoveJSONAttribute ...
// 删除属性
func RemoveJSONAttribute(jsonString string, attrName string) (string, error) {
	data := make(map[string]interface{})

	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return "", errors.WrapE(err, "json.Unmarshal failed", "json-text", jsonString)
	}

	delete(data, attrName)

	result, err := json.Marshal(data)
	if err != nil {
		return "", errors.WrapE(err, "json.Marshal failed", "json-data", data)
	}
	return string(result), nil
}

// GetJSONAttribute ...
// 获取属性（自动类型转换）
func GetJSONAttribute(jsonString string, attrName string) (string, error) {
	data := make(map[string]interface{})

	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return "", errors.WrapE(err, "json.Unmarshal failed", "json-text", jsonString)
	}

	value, exists := data[attrName]
	if !exists {
		return "", errors.E("json attr not found", "attr-name", attrName)
	}

	// 类型敏感转换
	switch v := value.(type) {
	case string:
		return v, nil
	case json.Number:
		return v.String(), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	default:
		// 复杂类型序列化
		if bytes, err := json.Marshal(v); err == nil {
			return string(bytes), nil
		}
	}
	return "", nil
}

// 辅助函数：解析布尔值
func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean")
}
