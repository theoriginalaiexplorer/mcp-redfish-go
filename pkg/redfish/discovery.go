package redfish

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	ssdpAddr = "239.255.255.250"
	ssdpPort = 1900
	ssdpMX   = 2
	ssdpST   = "urn:dmtf-org:service:redfish-rest:1"
)

// SSDPDiscovery handles SSDP discovery of Redfish endpoints
type SSDPDiscovery struct {
	timeout time.Duration
	logger  *slog.Logger
}

// NewSSDPDiscovery creates a new SSDP discovery instance
func NewSSDPDiscovery(timeout time.Duration, logger *slog.Logger) *SSDPDiscovery {
	if logger == nil {
		logger = slog.Default()
	}
	return &SSDPDiscovery{
		timeout: timeout,
		logger:  logger,
	}
}

// Discover performs SSDP M-SEARCH and returns discovered Redfish endpoints
func (d *SSDPDiscovery) Discover() ([]DiscoveredHost, error) {
	d.logger.Info("Starting SSDP discovery")

	// Create UDP socket
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(ssdpAddr),
		Port: ssdpPort,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP socket: %w", err)
	}
	defer conn.Close()

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(d.timeout))

	// Create M-SEARCH message
	message := fmt.Sprintf("M-SEARCH * HTTP/1.1\r\n"+
		"HOST: %s:%d\r\n"+
		"MAN: \"ssdp:discover\"\r\n"+
		"MX: %d\r\n"+
		"ST: %s\r\n\r\n",
		ssdpAddr, ssdpPort, ssdpMX, ssdpST)

	// Send M-SEARCH request
	_, err = conn.Write([]byte(message))
	if err != nil {
		return nil, fmt.Errorf("failed to send M-SEARCH: %w", err)
	}

	d.logger.Info("SSDP M-SEARCH sent, waiting for responses")

	var hosts []DiscoveredHost
	buffer := make([]byte, 1024)

	// Read responses until timeout
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				d.logger.Info("SSDP discovery timeout reached")
				break
			}
			d.logger.Warn("Error reading SSDP response", "error", err)
			continue
		}

		response := string(buffer[:n])
		alURI := d.parseAL(response)
		if alURI != "" && d.isValidServiceRoot(alURI) {
			host := DiscoveredHost{
				Address:     addr.IP.String(),
				ServiceRoot: alURI,
			}
			hosts = append(hosts, host)
			d.logger.Info("Discovered Redfish endpoint",
				"address", host.Address,
				"service_root", host.ServiceRoot)
		} else {
			d.logger.Debug("Received SSDP response but no valid AL header found",
				"address", addr.IP.String())
		}
	}

	d.logger.Info("SSDP discovery completed", "hosts_found", len(hosts))
	return hosts, nil
}

// parseAL extracts the AL (Alternate Location) header from SSDP response
func (d *SSDPDiscovery) parseAL(response string) string {
	// SSDP responses may have multiline headers, split and search each line
	lines := strings.Split(response, "\n")
	alRegex := regexp.MustCompile(`(?i)^AL:\s*(.+)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if matches := alRegex.FindStringSubmatch(line); len(matches) > 0 {
			return strings.TrimSpace(matches[1])
		}
	}
	return ""
}

// isValidServiceRoot validates that the URI is a Redfish service root endpoint
func (d *SSDPDiscovery) isValidServiceRoot(uri string) bool {
	parsed, err := url.Parse(uri)
	if err != nil {
		d.logger.Debug("Service root URI parse error", "uri", uri, "error", err)
		return false
	}

	// Must use HTTPS
	if parsed.Scheme != "https" {
		d.logger.Debug("Service root URI rejected (not https)", "uri", uri)
		return false
	}

	// Must have a host
	if parsed.Host == "" {
		d.logger.Debug("Service root URI rejected (missing host)", "uri", uri)
		return false
	}

	// Must end with /redfish/v1/ (allow optional trailing slash)
	pathRegex := regexp.MustCompile(`^/redfish/v1/?$`)
	if !pathRegex.MatchString(parsed.Path) {
		d.logger.Debug("Service root URI rejected (invalid path)", "uri", uri, "path", parsed.Path)
		return false
	}

	return true
}
