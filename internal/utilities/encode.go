package util

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

// EncodeString compress with zip and encode in base64
func EncodeString(in string) (string, error) {
	b := []byte(in)
	z, err := zipBytes(b)
	if err != nil {
		return in, fmt.Errorf("EncodeObject error after zipBytes %v", err)
	}
	return base64.StdEncoding.EncodeToString(z), nil
}

// DecodeString decodes the zip compressed input from base64
func DecodeString(in string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return in, fmt.Errorf("DecodeObject error after DecodeString %v", err)
	}
	b, err = unzipBytes(b)
	if err != nil {
		return in, fmt.Errorf("DecodeObject error after unzipBytes %v", err)
	}
	return string(b), nil
}

func zipBytes(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		return in, fmt.Errorf("zipBytes error after gz.Write() %v", err)
	}
	err = gz.Flush()
	if err != nil {
		return in, fmt.Errorf("zipBytes error after gz.Flush() %v", err)
	}
	err = gz.Close()
	if err != nil {
		return in, fmt.Errorf("zipBytes error after gz.Close() %v", err)
	}
	return b.Bytes(), nil
}

func unzipBytes(in []byte) ([]byte, error) {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		return in, fmt.Errorf("unzipBytes error after b.Write() %v", err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		return in, fmt.Errorf("unzipBytes error after gzip.NewReader() %v", err)
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		return in, fmt.Errorf("unzipBytes error after ioutil.ReadAll() %v", err)
	}
	return res, nil
}
