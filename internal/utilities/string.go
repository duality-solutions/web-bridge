package util

import "strconv"

// IsNumeric determines if the string val is a number
func IsNumeric(val string) bool {
	_, err := strconv.Atoi(val)
	return err == nil
}

// ToInt converts the string val to an int. Returns zero (0) if cast fails.
func ToInt(val string) int {
	ret, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return ret
}
