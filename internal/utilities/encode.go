package util

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// EncodeObject marshals the object to JSON, compresses the JSON with zip and returns a base64 encoded string
func EncodeObject(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("EncodeObject error after Marshal %v", err)
	}
	z, err := zipBytes(b)
	if err != nil {
		return "", fmt.Errorf("EncodeObject error after zipBytes %v", err)
	}
	return base64.StdEncoding.EncodeToString(z), nil
}

// DecodeObject decodes the base64 encoded string, unzips the compressed bytes and marshals the JSON into an object
func DecodeObject(in string, obj interface{}) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return fmt.Errorf("DecodeObject error after DecodeString %v", err)
	}

	uz, err := unzipBytes(b)
	if err != nil {
		return fmt.Errorf("DecodeObject error after unzipBytes %v", err)
	}

	err = json.Unmarshal(uz, obj)
	if err != nil {
		return fmt.Errorf("DecodeObject error after Unmarshal %v", err)
	}

	return nil
}

// EncodeString compresses the input string using zip and returns a base64 encoded string
func EncodeString(in string) (string, error) {
	b := []byte(in)
	z, err := zipBytes(b)
	if err != nil {
		return in, fmt.Errorf("EncodeObject error after zipBytes %v", err)
	}
	return base64.StdEncoding.EncodeToString(z), nil
}

// DecodeString decodes the base64 string, unzips the data and returns the original encoded string
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
