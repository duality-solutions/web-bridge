package dynamic

import (
	"fmt"
	"io/ioutil"
	"strings"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

// ParseDynamicConfValue returns the given configuration path
func ParseDynamicConfValue(confPath, parameter string) (string, error) {
	// TODO: look for comment char (#)
	var val = ""
	if util.FileExists(confPath) {
		dat, err := ioutil.ReadFile(confPath)
		if err != nil {
			return val, fmt.Errorf("ParseDynamicConfValue failed reading file %v. %v", confPath, err)
		}
		conf := string(dat[:len(dat)])
		index := strings.LastIndex(conf, parameter)
		if index >= 0 {
			val = conf[index+len(parameter):]
			lineEnd := strings.Index(val, "\n")
			if lineEnd > 0 {
				val = val[:lineEnd]
			}
			val = strings.Trim(val, "\n\r ")
		}
	} else {
		return val, fmt.Errorf("ParseDynamicConfValue failed. %v doesn't exist", confPath)
	}
	return val, nil
}
