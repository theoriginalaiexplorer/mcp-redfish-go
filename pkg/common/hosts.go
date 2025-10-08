package common

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"github.com/theoriginalaiexplorer/mcp-redfish-go/pkg/config"
	"github.com/theoriginalaiexplorer/mcp-redfish-go/pkg/redfish"
)

// HostManager manages both static and discovered Redfish hosts
type HostManager struct {
	staticHosts     []config.HostConfig
	discoveredHosts []redfish.DiscoveredHost
	mu              sync.RWMutex
	logger          *slog.Logger
}

// NewHostManager creates a new host manager
func NewHostManager(logger *slog.Logger) *HostManager {
	if logger == nil {
		logger = slog.Default()
	}

	hm := &HostManager{
		logger: logger,
	}

	// Load static hosts from environment
	hm.loadStaticHosts()

	return hm
}

// loadStaticHosts loads static hosts from REDFISH_HOSTS environment variable
func (hm *HostManager) loadStaticHosts() {
	hostsJSON := os.Getenv("REDFISH_HOSTS")
	if hostsJSON == "" {
		hostsJSON = `[{"address": "127.0.0.1"}]`
	}

	var hosts []config.HostConfig
	if err := json.Unmarshal([]byte(hostsJSON), &hosts); err != nil {
		hm.logger.Error("Failed to parse REDFISH_HOSTS", "error", err)
		hosts = []config.HostConfig{{Address: "127.0.0.1"}}
	}

	hm.mu.Lock()
	hm.staticHosts = hosts
	hm.mu.Unlock()

	hm.logger.Info("Loaded static hosts", "count", len(hosts))
}

// UpdateDiscoveredHosts updates the list of discovered hosts
func (hm *HostManager) UpdateDiscoveredHosts(hosts []redfish.DiscoveredHost) {
	hm.mu.Lock()
	hm.discoveredHosts = hosts
	hm.mu.Unlock()

	hm.logger.Info("Updated discovered hosts", "count", len(hosts))
}

// GetHosts returns the merged list of static and discovered hosts
// Static hosts take precedence over discovered hosts with the same address
func (hm *HostManager) GetHosts() []config.HostConfig {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	// Start with static hosts
	allHosts := make(map[string]config.HostConfig)
	for _, host := range hm.staticHosts {
		allHosts[host.Address] = host
	}

	// Add discovered hosts (only if not already present)
	for _, discovered := range hm.discoveredHosts {
		if _, exists := allHosts[discovered.Address]; !exists {
			// Convert discovered host to config format
			host := config.HostConfig{
				Address: discovered.Address,
				// Use defaults for other fields since discovered hosts don't provide them
			}
			allHosts[discovered.Address] = host
		}
	}

	// Convert map back to slice
	result := make([]config.HostConfig, 0, len(allHosts))
	for _, host := range allHosts {
		result = append(result, host)
	}

	return result
}

// GetHostByAddress finds a host by address
func (hm *HostManager) GetHostByAddress(address string) (config.HostConfig, bool) {
	hosts := hm.GetHosts()
	for _, host := range hosts {
		if host.Address == address {
			return host, true
		}
	}
	return config.HostConfig{}, false
}

// GetAddresses returns just the addresses of all hosts
func (hm *HostManager) GetAddresses() []string {
	hosts := hm.GetHosts()
	addresses := make([]string, len(hosts))
	for i, host := range hosts {
		addresses[i] = host.Address
	}
	return addresses
}
