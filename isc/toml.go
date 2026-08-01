package isc

import (
	"log"

	"github.com/pelletier/go-toml/v2"
)

// TomlToMap toml 字符串转嵌套 map[string]any（table 变嵌套 map，数组变切片）
func TomlToMap(contentOfToml string) (map[string]any, error) {
	resultMap := make(map[string]any)
	err := toml.Unmarshal([]byte(contentOfToml), &resultMap)
	if err != nil {
		log.Printf("TomlToMap, error: %v, content: %v", err, contentOfToml)
		return nil, err
	}
	return resultMap, nil
}

// TomlToProperties toml 字符串转扁平 properties（key=value 多行）
func TomlToProperties(contentOfToml string) (string, error) {
	dataMap, err := TomlToMap(contentOfToml)
	if err != nil {
		return "", err
	}
	return MapToProperties(dataMap)
}
