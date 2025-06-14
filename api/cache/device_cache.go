package cache

import (
	"sync"
	"time"

	onvif "go-onvif/internal"
)

// DeviceCache caches ONVIF devices to avoid repeated connections
type DeviceCache struct {
	devices      map[string]*onvif.Device
	lastAccessed map[string]time.Time
	mutex        sync.RWMutex
	ttl          time.Duration
}

var GlobalDeviceCache *DeviceCache

// NewDeviceCache creates and returns a new device cache
func NewDeviceCache(ttl time.Duration) *DeviceCache {
	cache := &DeviceCache{
		devices:      make(map[string]*onvif.Device),
		lastAccessed: make(map[string]time.Time),
		ttl:          ttl,
	}
	go cache.startCleanup()
	return cache
}

// startCleanup periodically removes old entries
func (c *DeviceCache) startCleanup() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup removes expired devices
func (c *DeviceCache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	now := time.Now()
	for key, lastAccess := range c.lastAccessed {
		if now.Sub(lastAccess) > c.ttl {
			delete(c.devices, key)
			delete(c.lastAccessed, key)
		}
	}
}

// GetDevice returns a cached device or creates and caches a new one
func (c *DeviceCache) GetDevice(xaddr, username, password string) (*onvif.Device, error) {

	cacheKey := xaddr
	c.mutex.RLock()
	dev, found := c.devices[cacheKey]
	c.mutex.RUnlock()

	if found {
		c.mutex.Lock()
		c.lastAccessed[cacheKey] = time.Now()
		c.mutex.Unlock()
		return dev, nil
	}

	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    xaddr,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	c.mutex.Lock()
	c.devices[cacheKey] = dev
	c.lastAccessed[cacheKey] = time.Now()
	c.mutex.Unlock()

	return dev, nil
}

func (c *DeviceCache) GetAllDevices() map[string]*onvif.Device {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Create a copy to avoid race conditions
	copy := make(map[string]*onvif.Device, len(c.devices))
	for key, dev := range c.devices {
		copy[key] = dev
	}
	return copy
}
