package dynamic

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ActiveLinks stores the completed link list returned by
type ActiveLinks struct {
	Links       []Link `json:"link"`
	LockedLinks int    `json:"locked_links"`
}

// GetActiveLinks returns all the active links
func (d *Dynamicd) GetActiveLinks() (*ActiveLinks, error) {
	var linksGeneric map[string]interface{}
	req, _ := NewRequest("dynamic-cli link complete")
	rawResp := []byte(<-d.ExecCmdRequest(req))
	errUnmarshal := json.Unmarshal(rawResp, &linksGeneric)
	if errUnmarshal != nil {
		fmt.Println("Outer error", errUnmarshal)
		return nil, errUnmarshal
	}
	var links ActiveLinks
	for k, v := range linksGeneric {
		if strings.HasPrefix(k, "link-") {
			b, err := json.Marshal(v)
			if err == nil {
				var link Link
				errUnmarshal = json.Unmarshal(b, &link)
				if errUnmarshal != nil {
					fmt.Println("Inner error", errUnmarshal)
					return nil, errUnmarshal
				}

				links.Links = append(links.Links, link)
			}
		}
	}
	return &links, nil
}
