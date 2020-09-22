package bridge

import "sync"

// Controller hold all link WebRTC bridges
type Controller struct {
	bridgesMut  *sync.RWMutex
	connected   map[string]*Bridge
	unconnected map[string]*Bridge
}

// NewController creates a new controller struct and initizes the bridges map
func NewController() *Controller {
	return &Controller{
		bridgesMut:  new(sync.RWMutex),
		connected:   make(map[string]*Bridge),
		unconnected: make(map[string]*Bridge),
	}
}

// GetConnected returns the connected bridge by key from the controller
func (c *Controller) GetConnected(key string) *Bridge {
	c.bridgesMut.RLock()
	defer c.bridgesMut.RUnlock()
	return c.connected[key]
}

// PutConnected adds the connected bridge by key to the controller
func (c *Controller) PutConnected(bridge *Bridge) {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	c.connected[bridge.LinkID()] = bridge
}

// GetUnconnected returns the unconnected bridge by key from the controller
func (c *Controller) GetUnconnected(key string) *Bridge {
	c.bridgesMut.RLock()
	defer c.bridgesMut.RUnlock()
	return c.unconnected[key]
}

// PutUnconnected adds the unconnected bridge by key to the controller
func (c *Controller) PutUnconnected(bridge *Bridge) {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	c.unconnected[bridge.LinkID()] = bridge
}

// MoveConnectedToUnconnected moves a bridge from connected to unconnected map
func (c *Controller) MoveConnectedToUnconnected(bridge *Bridge) {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	if c.connected[bridge.LinkID()] != nil {
		delete(c.connected, bridge.LinkID())
	}
	c.unconnected[bridge.LinkID()] = bridge
}

// MoveUnconnectedToConnected moves a bridge from connected to unconnected map
func (c *Controller) MoveUnconnectedToConnected(bridge *Bridge) {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	if c.unconnected[bridge.LinkID()] != nil {
		delete(c.unconnected, bridge.LinkID())
	}
	c.connected[bridge.LinkID()] = bridge
}

// Connected returns a map containing the currently connected bridges
func (c *Controller) Connected() map[string]*Bridge {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	connected := make(map[string]*Bridge)
	for _, bridge := range c.connected {
		connected[bridge.LinkID()] = bridge
	}
	return connected
}

// Unconnected returns a map containing the currently unconnected bridges
func (c *Controller) Unconnected() map[string]*Bridge {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	unconnected := make(map[string]*Bridge)
	for _, bridge := range c.unconnected {
		unconnected[bridge.LinkID()] = bridge
	}
	return unconnected
}

// Count returns the total bridge count in both connected and unconnected maps
func (c *Controller) Count() uint16 {
	c.bridgesMut.RLock()
	defer c.bridgesMut.RUnlock()
	return uint16(len(c.connected) + len(c.unconnected))
}

// AllBridges returns a map containing all connected and unconnected bridges
func (c *Controller) AllBridges() map[string]*Bridge {
	c.bridgesMut.Lock()
	defer c.bridgesMut.Unlock()
	allBridges := make(map[string]*Bridge)
	for _, bridge := range c.connected {
		allBridges[bridge.LinkID()] = bridge
	}
	for _, bridge := range c.unconnected {
		allBridges[bridge.LinkID()] = bridge
	}
	return allBridges
}
