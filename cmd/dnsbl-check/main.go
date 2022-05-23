package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Ajnasz/dnsbl-check/dnsblprovider"
	"github.com/Ajnasz/dnsbl-check/dnsblprovider/io"
	"github.com/Ajnasz/dnsbl-check/providerlist"
)

func processLookupResult(result io.LookupResult) {
	if result.Error != nil {
		fmt.Println(fmt.Sprintf("ERR\t%s\t%s\t%s", result.Address, result.Provider.GetName(), result.Error))
		return
	}
	if result.IsBlacklisted {
		var reason string

		if result.Reason == "" {
			reason = "unknown reason"
		} else {
			reason = result.Reason
		}

		fmt.Println(fmt.Sprintf("FAIL\t%s\t%s\t%s", result.Address, result.Provider.GetName(), reason))
	} else {
		fmt.Println(fmt.Sprintf("OK\t%s\t%s", result.Address, result.Provider.GetName()))
	}
}

func main() {
	var domainsFile = flag.String("p", "", "path to file which stores list of dnsbl checks, empty or - for stdin")
	var addressesParam = flag.String("i", "", "IP Address to check, separate by comma for a list")

	flag.Parse()
	list, err := providerlist.GetProvidersChan(*domainsFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading domains")
		os.Exit(1)
	}

	var providers []dnsblprovider.DNSBLProvider

	for item := range list {
		provider := dnsblprovider.GeneralProvider{
			URL: item,
		}

		providers = append(providers, provider)
	}

	addresses := providerlist.GetAddresses(*addressesParam)
	for result := range io.LookupMany(addresses, providers) {
		processLookupResult(result)
	}
}
