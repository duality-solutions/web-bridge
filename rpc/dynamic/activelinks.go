package dynamic

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

// ActiveLinks stores the completed link list returned by
type ActiveLinks struct {
	Links       []Link `json:"link"`
	LockedLinks int    `json:"locked_links"`
}

func newActiveLinks() *ActiveLinks {
	var links ActiveLinks
	links.Links = []Link{}
	links.LockedLinks = 0
	return &links
}

func (d *Dynamicd) getLinks() (*ActiveLinks, error) {
	var linksGeneric map[string]interface{}
	req, _ := NewRequest("dynamic-cli link complete")
	rawResp := []byte(<-d.ExecCmdRequest(req))
	err := json.Unmarshal(rawResp, &linksGeneric)
	if err == nil {
		links := newActiveLinks()
		for k, v := range linksGeneric {
			if strings.HasPrefix(k, "link-") {
				b, err := json.Marshal(v)
				if err == nil {
					var link Link
					err = json.Unmarshal(b, &link)
					if err != nil {
						util.Error.Println("getLinks inner unmarshal error", err)
						return nil, err
					}
					links.Links = append(links.Links, link)
				}
			}
		}
		return links, nil
	}
	util.Error.Println("getLinks unmarshal error", err)
	return nil, err
}

// GetActiveLinks returns all the active links
func (d *Dynamicd) GetActiveLinks(timeout time.Duration) (*ActiveLinks, error) {
	activeLinks, err := d.getLinks()
	if err != nil {
		for {
			select {
			case <-time.After(time.Second * 5):
				activeLinks, err = d.getLinks()
				if err == nil {
					return activeLinks, nil
				}
			case <-time.After(timeout):
				return nil, fmt.Errorf("GetActiveLinks failed after timeout")
			}
		}
	} else {
		return activeLinks, nil
	}
}
