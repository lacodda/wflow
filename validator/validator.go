package validator

import (
	"errors"
	"strconv"
)

func IsNumber(input interface{}) error {
	str, _ := input.(string)
	_, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return errors.New("Invalid number")
	}
	return nil
}
