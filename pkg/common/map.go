package common

import (
	"encoding/json"
)

// ToMap 将 JSON 字符串转换为 map[string]interface{}
func ToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ToJSON 将 map[string]interface{} 转换为 JSON 字符串
func ToJSON(m map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ToMapStringInt 将 JSON 字符串转换为 map[string]int
func ToMapStringInt(jsonStr string) (map[string]int, error) {
	var result map[string]int
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ToMapStringString 将 JSON 字符串转换为 map[string]string
func ToMapStringString(jsonStr string) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/*
func main() {
	// 示例 1: {"sema1":1,"sema2":2,"sema3":3} 与 map[string]int 相互转换
	jsonStr1 := `{"sema1":1,"sema2":2,"sema3":3}`
	mapInt, err := ToMapStringInt(jsonStr1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Map[string]int:", mapInt)
	}

	jsonStrFromMapInt, err := ToJSON(map[string]interface{}{"sema1": 1, "sema2": 2, "sema3": 3})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("JSON from map[string]int:", jsonStrFromMapInt)
	}

	// 示例 2: {"var1":"s1","var2":"s2","var3":"s3"} 与 map[string]string 相互转换
	jsonStr2 := `{"var1":"s1","var2":"s2","var3":"s3"}`
	mapString, err := ToMapStringString(jsonStr2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Map[string]string:", mapString)
	}

	jsonStrFromMapString, err := ToJSON(map[string]interface{}{"var1": "s1", "var2": "s2", "var3": "s3"})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("JSON from map[string]string:", jsonStrFromMapString)
	}
}
*/
