package util

import (
	"crypto/rand"
)

const (
	letters     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 62 possibilities
	lettersBits = 6                                                                // 6 bits to represent 64 possibilities / indexes
	bitmask     = 1<<lettersBits - 1                                               // All 1-bits, as many as lettersBits
)

// RandomString returns the requested number of random charactors
func RandomString(length uint) (string, error) {
	r := make([]byte, length)
	bs := int(float64(length) * 1.3)
	var err error
	for i, j, rb := 0, 0, []byte{}; uint(i) < length; j++ {
		if j%bs == 0 {
			rb, err = RandomBytes(uint(bs))
			if err != nil {
				return "", err
			}
		}
		if idx := uint(rb[j%int(length)] & bitmask); idx < uint(len(letters)) {
			r[i] = letters[idx]
			i++
		}
	}

	return string(r), nil
}

// RandomBytes returns the requested number of random bytes
func RandomBytes(length uint) ([]byte, error) {
	var rb = make([]byte, length)
	_, err := rand.Read(rb)
	if err != nil {
		return nil, err
	}
	return rb, nil
}
