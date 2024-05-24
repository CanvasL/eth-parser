package utils

import (
	"fmt"
	"strconv"
)

func HexStringToInt(hexStr string) (int, error) {
	if val, err := strconv.ParseInt(hexStr[2:], 16, 64); err != nil {
		return -1, err
	} else {
		return int(val), nil
	}
}

func IntToHexString(num int) string {
	return fmt.Sprintf("0x%x", num)
}
