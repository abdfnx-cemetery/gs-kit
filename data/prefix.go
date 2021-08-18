package data

import (
	"strings"
	"unicode"
)

func PrefixedBy(input interface{}, prefix string) (interface{}, error) {
	normalized, err := Normalize(input)
	if err != nil {
		return input, err
	}

	input = normalized

	if inputMap, ok := input.(map[string]interface{}); ok {
		converted := make(map[string]interface{}, len(inputMap))
		for k, v := range inputMap {
			if strings.HasPrefix(k, prefix) {
				key := uncapitalize(strings.TrimPrefix(k, prefix))
				converted[key] = v
			}
		}

		return converted, nil
	} else if inputMap, ok := input.(map[string]string); ok {
		converted := make(map[string]string, len(inputMap))
		for k, v := range inputMap {
			if strings.HasPrefix(k, prefix) {
				key := uncapitalize(strings.TrimPrefix(k, prefix))
				converted[key] = v
			}
		}

		return converted, nil
	}

	return input, nil
}

func uncapitalize(str string) string {
	if len(str) == 0 {
		return str
	}

	vv := []rune(str)
	vv[0] = unicode.ToLower(vv[0])

	return string(vv)
}
