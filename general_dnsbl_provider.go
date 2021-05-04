package main

import (
	"fmt"
	"net"
	"strings"
)

func isNoHostError(err error) bool {
	if serr, ok := err.(*net.DNSError); ok {
		return serr.IsNotFound
	}

	return false
}

// GeneralDNSBLProvider implements DNSBLProvider
// URL is a required property which should be the ending of the dnsbl hostname
type GeneralDNSBLProvider struct {
	URL string
}

// GetName returns the name of the provider
// Now it's the URL
func (provider GeneralDNSBLProvider) GetName() string {
	return provider.URL
}

// getAddress returns the address what should be queried
// Combines the IP address (octets reversed) and the provider URL
func (provider GeneralDNSBLProvider) getAddress(ip string) string {
	return fmt.Sprintf("%s.%s", reverseIPAddress(ip), provider.URL)
}

// GetReason returns the block reason for an IP address
func (provider GeneralDNSBLProvider) GetReason(ip string) (string, error) {
	texts, err := net.LookupTXT(provider.getAddress(ip))
	if err != nil {
		if isNoHostError(err) {
			return "", nil
		}

		return "", err
	}

	return strings.Join(texts, ""), nil
}

// IsBlacklistedIP returns if the IP address listed at a provider
func (provider GeneralDNSBLProvider) IsBlacklistedIP(ip string) (bool, error) {
	names, err := net.LookupIP(provider.getAddress(ip))

	if err != nil {
		if isNoHostError(err) {
			return false, nil
		}

		return false, nil
	}

	if len(names) == 0 {
		return false, nil
	}

	return true, nil

}
