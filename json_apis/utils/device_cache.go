package utils

import (
	"sync"
	"time"

	"github.com/BalkarSandhu/go-onvif/onvif"
)

// DeviceCache provides caching for ONVIF devices to avoid repeated connections
type DeviceCache struct {
	devices      map[string]*onvif.Device
	lastAccessed map[string]time.Time
	mutex        sync.RWMutex
	ttl          time.Duration
}

// startCleanup periodically removes old cache entries
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
		// Remove entries that haven't been accessed within the TTL period
		if now.Sub(lastAccess) > c.ttl {
			// Remove from both maps
			delete(c.devices, key)
			delete(c.lastAccessed, key)
		}
	}
}

// GetDevice retrieves a device from cache or creates a new one
func (c *DeviceCache) GetDevice(xaddr, username, password string) (*onvif.Device, error) {
	cacheKey := xaddr + "|" + username

	// Try to get from cache first
	c.mutex.RLock()
	dev, found := c.devices[cacheKey]
	c.mutex.RUnlock()

	if found {
		return dev, nil
	}

	// Create new device connection
	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    xaddr,
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	// Store in cache
	c.mutex.Lock()
	c.devices[cacheKey] = dev
	c.mutex.Unlock()

	return dev, nil
}

// NewDeviceCache creates a new device cache with the specified TTL
func NewDeviceCache(ttl time.Duration) *DeviceCache {
	cache := &DeviceCache{
		devices: make(map[string]*onvif.Device),
		ttl:     ttl,
	}
	// Start a goroutine to clean up expired cache entries
	go cache.startCleanup()
	return cache
}
