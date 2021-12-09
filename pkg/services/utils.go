package services

import (
	"net/http"
	"strconv"
)

func parseQuery(r *http.Request, params map[string]string) (map[string]interface{}, error) {
	var result = make(map[string]interface{}, len(params))
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	for name, typ := range params {
		if values, ok := r.Form[name]; ok {
			if values[0] == "" {
				continue
			}
			switch typ {
			default:
				result[name] = values[0]
			case "stringArray":
				result[name] = values
			case "uintArray":
				result[name] = make([]uint64, len(values))
				for i, value := range values {
					result[name].([]uint64)[i] = parseString(value, "uint").(uint64)
				}
			}
		}
	}

	return result, nil
}

func parseString(value, typ string) interface{} {
	switch typ {
	case "bool":
		parsedValue, err := strconv.ParseBool(value)
		if err == nil {
			return parsedValue
		}
	case "uint":
		parsedValue, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			return parsedValue
		}
	case "int":
		parsedValue, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return parsedValue
		}
	case "float":
		parsedValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return parsedValue
		}
	}

	return value
}
