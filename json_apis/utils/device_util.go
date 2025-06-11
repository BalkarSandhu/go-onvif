package utils

import (
	"context"
	"fmt"
	"go-onvif/device"
	"go-onvif/onvif"
	"net/http"
	"sync"
	"time"
)

// Constants
const (
	maxConcurrentScans = 50
	onvifPort          = 80
)

// DeviceCache caches ONVIF devices to avoid repeated connections
type DeviceCache struct {
	devices      map[string]*onvif.Device
	lastAccessed map[string]time.Time
	mutex        sync.RWMutex
	ttl          time.Duration
}

// DeviceScanResult represents a single scan result
type DeviceScanResult struct {
	IP              string `json:"ip"`
	Alive           bool   `json:"alive"`
	Manufacturer    string `json:"manufacturer,omitempty"`
	Model           string `json:"model,omitempty"`
	FirmwareVersion string `json:"firmware_version,omitempty"`
	SerialNumber    string `json:"serial_number,omitempty"`
	HardwareID      string `json:"hardware_id,omitempty"`
	Error           string `json:"error,omitempty"`
}

// ScanRequest represents the scan input from user/API
type ScanRequest struct {
	StartIP  string `form:"start" binding:"required"`
	EndIP    string `form:"end" binding:"required"`
	Username string `form:"username"`
	Password string `form:"password"`
	Timeout  int    `form:"timeout"`
}

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
	cacheKey := xaddr + "|" + username
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

// PerformDeviceScan performs scanning for a list of IPs
func (s *DeviceScanResult) PerformDeviceScan(ips []string, req ScanRequest, acceptedData []byte, timeout time.Duration) []DeviceScanResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentScans)
	results := make([]DeviceScanResult, len(ips))

	for i, ip := range ips {
		wg.Add(1)
		go func(i int, ip string) {
			defer wg.Done()

			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				results[i] = DeviceScanResult{
					IP:    ip,
					Alive: false,
					Error: "Scan timeout",
				}
				return
			}

			results[i] = s.scanSingleDevice(ctx, ip, req, acceptedData)
		}(i, ip)
	}

	wg.Wait()
	return results
}

// scanSingleDevice handles ONVIF detection and info fetching
func (s *DeviceScanResult) scanSingleDevice(ctx context.Context, ip string, req ScanRequest, acceptedData []byte) DeviceScanResult {
	result := DeviceScanResult{IP: ip}

	select {
	case <-ctx.Done():
		result.Error = "Scan cancelled"
		return result
	default:
	}

	if req.Username == "" || req.Password == "" {
		if s.isONVIFDeviceAvailable(ctx, ip) {
			result.Alive = true
		}
		return result
	}

	dev, err := s.createONVIFDevice(ctx, ip, req.Username, req.Password)
	if err != nil {
		result.Error = fmt.Sprintf("Connection failed: %v", err)
		return result
	}

	devInfo, err := s.getDeviceInformation(ctx, dev, acceptedData)
	if err != nil {
		result.Alive = true
		result.Error = fmt.Sprintf("Failed to get device info: %v", err)
		return result
	}

	result.Alive = true
	result.Manufacturer = devInfo.Manufacturer
	result.Model = devInfo.Model
	result.FirmwareVersion = devInfo.FirmwareVersion
	result.SerialNumber = devInfo.SerialNumber
	result.HardwareID = devInfo.HardwareId

	return result
}

// isONVIFDeviceAvailable performs a lightweight availability check
func (s *DeviceScanResult) isONVIFDeviceAvailable(ctx context.Context, ip string) bool {
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return IsONVIFDeviceWithContext(checkCtx, ip)
}

// createONVIFDevice initializes an ONVIF device
func (s *DeviceScanResult) createONVIFDevice(ctx context.Context, ip, username, password string) (*onvif.Device, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	return onvif.NewDevice(onvif.DeviceParams{
		Xaddr:      fmt.Sprintf("%s:%d", ip, onvifPort),
		Username:   username,
		Password:   password,
		HttpClient: client,
	})
}

// getDeviceInformation fetches and type-asserts ONVIF device information
func (s *DeviceScanResult) getDeviceInformation(ctx context.Context, dev *onvif.Device, acceptedData []byte) (*device.GetDeviceInformationResponse, error) {
	resultChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	go func() {
		info, err := CallDeviceMethod("GetDeviceInformation", dev, acceptedData)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- info
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
	case raw := <-resultChan:
		if info, ok := raw.(device.GetDeviceInformationResponse); ok {
			return &info, nil
		}
		return nil, fmt.Errorf("unexpected response type")
	}
}

// CountAliveDevices returns number of devices marked alive
func CountAliveDevices(results []DeviceScanResult) int {
	count := 0
	for _, r := range results {
		if r.Alive {
			count++
		}
	}
	return count
}

// Placeholder: implement or adapt to your networking logic
func IsONVIFDeviceWithContext(ctx context.Context, ip string) bool {
	return IsONVIFDevice(ip)
}
