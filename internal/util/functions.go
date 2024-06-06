package util

import (
	"fmt"
	"strconv"
)

func ConvertMapToStringVals(m map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for key, val := range m {
		if val == nil {
			out[key] = ""
		} else {
			switch v := val.(type) {
			case string:
				out[key] = v
			case bool:
				out[key] = strconv.FormatBool(v)
			case float64:
				out[key] = strconv.FormatFloat(v, 'f', -1, 64)
			case int:
				out[key] = strconv.Itoa(v)
			case int64:
				out[key] = strconv.FormatInt(v, 10)
			default:
				out[key] = fmt.Sprintf("%v", v)
			}
		}
	}
	return out
}
