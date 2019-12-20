package utils

import "fmt"

func MarshalError(code string, message string) []byte {
	return []byte(fmt.Sprintf(`{
		"error": {
			"message": "%s",
			"code": "%s" 
		}
	}`, message, code))
}
