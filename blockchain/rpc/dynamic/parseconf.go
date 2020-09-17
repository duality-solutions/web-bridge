package dynamic

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/duality-solutions/web-bridge/internal/util"
)

// RemoveConfigComments removes the comments from the config string
func RemoveConfigComments(conf string) string {
	for strings.Index(conf, "#") > 0 {
		i := strings.Index(conf, "#")
		n := strings.Index(conf[i:], "\n")
		if n < 0 {
			n = len(conf) - i - 1
		}
		remove := conf[i : i+n+1]
		conf = strings.Replace(conf, remove, "", -1)
	}
	return conf
}

// GetDynamicConfig returns the configuration file string for the given path
func GetDynamicConfig(confPath string) (string, error) {
	if !util.FileExists(confPath) {
		return "", fmt.Errorf("GetDynamicConfig failed. %v doesn't exist", confPath)
	}
	dat, err := ioutil.ReadFile(confPath)
	if err != nil {
		return "", fmt.Errorf("GetDynamicConfig failed reading file %v. %v", confPath, err)
	}
	dataLen := len(dat)
	if dataLen == 0 {
		return "", nil
	}
	return RemoveConfigComments(string(dat[:dataLen])), nil
}

// ParseDynamicConfigValue returns paramter value for the given configuration path file and paramter name
func ParseDynamicConfigValue(conf, parameter string) (string, error) {
	var val = ""
	if len(parameter) < 2 {
		return val, fmt.Errorf("ParseDynamicConfigValue parameter (%v) not long enough %v", parameter, len(parameter))
	}
	if parameter[len(parameter)-1:] != "=" {
		parameter += "="
	}
	index := strings.LastIndex(conf, parameter)
	if index >= 0 {
		val = conf[index+len(parameter):]
		lineEnd := strings.Index(val, "\n")
		if lineEnd > 0 {
			val = val[:lineEnd]
		}
		val = strings.Trim(val, "\n\r ")
	}
	return val, nil
}
