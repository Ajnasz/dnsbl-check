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

// getAddress returns the address what should be queried
// Combines the IP address (octets reversed) and the provider URL
func (provider GeneralDNSBLProvider) getAddress(address string) string {
	if provider.isIPAddress(address) {
		return fmt.Sprintf("%s.%s", reverseIPAddress(address), provider.URL)
	}

	return fmt.Sprintf("%s.%s", address, provider.URL)
}

func (provider GeneralDNSBLProvider) isIPAddress(address string) bool {
	return net.ParseIP(address) != nil
}

// GetName returns the name of the provider
// Now it's the URL
func (provider GeneralDNSBLProvider) GetName() string {
	return provider.URL
}

// GetReason returns the block reason for an IP address
func (provider GeneralDNSBLProvider) GetReason(address string) (string, error) {
	texts, err := net.LookupTXT(provider.getAddress(address))
	if err != nil {
		if isNoHostError(err) {
			return "", nil
		}

		return "", err
	}

	return strings.Join(texts, ""), nil
}

// IsBlacklisted returns if the IP address listed at a provider
func (provider GeneralDNSBLProvider) IsBlacklisted(address string) (bool, error) {
	names, err := net.LookupIP(provider.getAddress(address))

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
