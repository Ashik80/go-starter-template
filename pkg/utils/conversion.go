package utils

import (
	"fmt"
	"strconv"
)

func ParseToInt(value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %w", err)
	}
	return n, nil
}

func ParseToString(value int) string {
	n := strconv.Itoa(value)
	return n
}
