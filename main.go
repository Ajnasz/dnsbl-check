package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"unsafe"
)

func getToken(str, delimiter, out string) string {
	if str == "" {
		return out
	}

	if string(str[0]) == delimiter {
		return out
	}

	return getToken(str[1:], delimiter, out+string(str[0]))
}

func reverseStringByToken(str string, delimiter string, out string) string {
	if str == "" {
		return out
	}

	token := getToken(str, delimiter, "")

	if out == "" {
		out = token
	} else {
		out = token + delimiter + out
	}

	if len(str) <= len(token)+1 {
		return out
	}

	return reverseStringByToken(str[len(token)+1:], delimiter, out)
}

func reverseIPAddress(str string) string {
	return reverseStringByToken(str, ".", "")
}

// DNSBLProvider interface should be implemented to be able to query a provider
type DNSBLProvider interface {
	GetName() string
	IsBlacklisted(string) (bool, error)
	GetReason(string) (string, error)
}

// LookupResult stores the query result with reason
type LookupResult struct {
	IsBlacklisted bool
	Address       string
	Reason        string
	Provider      DNSBLProvider
}

func getBlacklists(address string, providers []DNSBLProvider) chan LookupResult {
	var wg sync.WaitGroup
	wg.Add(len(providers))

	results := make(chan LookupResult)

	for _, provider := range providers {
		go func(provider DNSBLProvider) {
			defer wg.Done()
			isListed, err := provider.IsBlacklisted(address)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if isListed {
				desc, err := provider.GetReason(address)

				if err != nil {
					fmt.Println("ERROR", err.Error())
				}

				results <- LookupResult{
					Address:       address,
					IsBlacklisted: true,
					Provider:      provider,
					Reason:        desc,
				}
				return
			}

			results <- LookupResult{
				Address:       address,
				IsBlacklisted: false,
				Provider:      provider,
			}
		}(provider)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func isValidIP(address string) bool {
	return net.ParseIP(address) != nil
}

func isCommentLine(line string) bool {
	return strings.HasPrefix(line, "#")
}

func filterString(lines []string, test func(string) bool) []string {
	var out []string

	for _, line := range lines {
		if test(line) {
			out = append(out, line)
		}
	}

	return out
}

func mapString(lines []string, conv func(string) string) []string {
	var out []string

	for _, line := range lines {
		out = append(out, conv(line))
	}

	return out
}

func isEmptyString(str string) bool {
	return str == ""
}

func negate(f func(string) bool) func(string) bool {
	return func(str string) bool {
		r := f(str)
		return !r
	}
}

func getProviders(fn string) ([]string, error) {
	f, err := os.ReadFile(fn)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(f), "\n")

	trimmed := mapString(lines, strings.TrimSpace)
	noEmpty := filterString(trimmed, negate(isEmptyString))
	noComment := filterString(noEmpty, negate(isCommentLine))

	return noComment, nil
}

func main() {
	var domainsFile = flag.String("p", "./providers", "path to file which stores list of dnsbl checks")
	var addressesParam = flag.String("i", "", "IP Address to check, separate by comma for a list")

	flag.Parse()

	list, err := getProviders(*domainsFile)

	if err != nil {
		fmt.Println("Error reading domains")
		os.Exit(1)
	}

	var providers []DNSBLProvider

	s := 0
	for _, item := range list {
		provider := GeneralDNSBLProvider{
			URL: item,
		}
		s += int(unsafe.Sizeof(provider))
		providers = append(providers, provider)
	}

	addresses := filterString(strings.Split(*addressesParam, ","), negate(isEmptyString))

	var addressWg sync.WaitGroup
	addressWg.Add(len(addresses))

	for _, address := range addresses {
		go func(address string) {
			defer addressWg.Done()
			for result := range getBlacklists(address, providers) {
				if result.IsBlacklisted {
					var reason string

					if result.Reason == "" {
						reason = "unkown reason"
					} else {
						reason = result.Reason
					}

					fmt.Println(fmt.Sprintf("ERR\t%s\t%s\t%s", result.Address, result.Provider.GetName(), reason))
				} else {
					fmt.Println(fmt.Sprintf("OK\t%s\t%s", result.Address, result.Provider.GetName()))
				}
			}
		}(address)
	}
	addressWg.Wait()
}
