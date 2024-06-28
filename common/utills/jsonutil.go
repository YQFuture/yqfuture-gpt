package utills

import (
	"encoding/json"
)

func AnyToString(source any) (string, error) {
	marshal, err := json.Marshal(source)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func StringToAny(source string, target any) error {
	err := json.Unmarshal([]byte(source), target)
	if err != nil {
		return err
	}
	return nil
}
