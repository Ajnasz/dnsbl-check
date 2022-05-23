package io

import (
	"sync"

	"github.com/Ajnasz/dnsbl-check/dnsblprovider"
)

// LookupResult stores the query result with reason
type LookupResult struct {
	IsBlacklisted bool
	Address       string
	Reason        string
	Provider      dnsblprovider.DNSBLProvider
	Error         error
}

// Lookup will check if an address is listed on a balcklist, returns a structured result
func Lookup(address string, provider dnsblprovider.DNSBLProvider) LookupResult {
	isListed, err := provider.IsBlacklisted(address)
	if err != nil {
		return LookupResult{
			Provider: provider,
			Address:  address,
			Error:    err,
		}
	}

	if isListed {
		desc, err := provider.GetReason(address)

		return LookupResult{
			Error:         err,
			Address:       address,
			IsBlacklisted: true,
			Provider:      provider,
			Reason:        desc,
		}
	}

	return LookupResult{
		Address:       address,
		IsBlacklisted: false,
		Provider:      provider,
	}
}

// LookupMany will check multiple addresses on multiple providers, return the result in a chan
func LookupMany(addresses []string, providers []dnsblprovider.DNSBLProvider) chan LookupResult {
	var wg sync.WaitGroup
	results := make(chan LookupResult)

	for _, address := range addresses {
		wg.Add(len(providers))

		for _, provider := range providers {
			go func(address string, provider dnsblprovider.DNSBLProvider) {
				defer wg.Done()
				results <- Lookup(address, provider)
			}(address, provider)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
