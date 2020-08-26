package util

import (
	"fmt"
	"math/big"
	"strconv"
)

// ScientificNotationToInt64 converts a scientific notation string to a 64 bit integer if possible
func ScientificNotationToInt64(strScientificNotation string) (int64, error) {
	flt, _, err := big.ParseFloat(strScientificNotation, 10, 0, big.ToNearestEven)
	if err != nil {
		return 0, err
	}
	fltVal := fmt.Sprintf("%.0f", flt)
	intVal, err := strconv.ParseInt(fltVal, 10, 64)
	if err != nil {
		return 0, err
	}
	return intVal, nil
}
