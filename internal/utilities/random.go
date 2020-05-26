package util

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
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

// RandomHashString returns the requested number of random charactors
func RandomHashString(length uint) (string, error) {
	var ranlen uint = 24
	if length > ranlen {
		ranlen = length
	}
	s, err := RandomString(ranlen)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	hash.Write([]byte(s))
	bs := hash.Sum(nil)
	hs := fmt.Sprintf("%x", bs)
	if len(hs) > int(length) {
		hs = hs[0 : length-1]
	}
	return hs, nil
}
