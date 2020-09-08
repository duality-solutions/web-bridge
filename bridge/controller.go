package bridge

import "sync"

// Controller hold all link WebRTC bridges
type Controller struct {
	bridgesMut sync.RWMutex
	bridges    map[string]*Bridge
}

// NewController creates a new controller struct and initizes the bridges map
func NewController() *Controller {
	return &Controller{
		bridges: make(map[string]*Bridge),
	}
}

// Get returns the bridge by key from the controller
func (c *Controller) Get(key string) *Bridge {
	c.bridgesMut.RLock()
	defer c.bridgesMut.RUnlock()
	return c.bridges[key]
}

// Put adds the bridge by key t0 the controller
func (c *Controller) Put(bridge *Bridge) {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	c.bridges[bridge.LinkID()] = bridge
}

// Connected returns a map containing the currently connected bridges
func (c *Controller) Connected() map[string]*Bridge {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	connected := make(map[string]*Bridge)
	for _, bridge := range c.bridges {
		if bridge.State == StateOpenConnection {
			connected[bridge.LinkID()] = bridge
		}
	}
	return connected
}

// Unconnected returns a map containing the currently unconnected bridges
func (c *Controller) Unconnected() map[string]*Bridge {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	unconnected := make(map[string]*Bridge)
	for _, bridge := range c.bridges {
		if bridge.State != StateOpenConnection {
			unconnected[bridge.LinkID()] = bridge
		}
	}
	return unconnected
}
