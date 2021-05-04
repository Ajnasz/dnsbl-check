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
	IsBlacklistedIP(string) (bool, error)
	GetReason(string) (string, error)
}

// LookupResult stores the query result with reason
type LookupResult struct {
	IsBlacklisted bool
	IP            string
	Reason        string
	Provider      DNSBLProvider
}

func getBlacklists(ip string, providers []DNSBLProvider) chan LookupResult {
	var wg sync.WaitGroup
	wg.Add(len(providers))
	errChan := make(chan LookupResult)
	for _, provider := range providers {
		go func(provider DNSBLProvider) {
			defer wg.Done()
			isListed, err := provider.IsBlacklistedIP(ip)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if isListed {
				desc, err := provider.GetReason(ip)

				if err != nil {
					fmt.Println("ERROR", err.Error())
				}

				errChan <- LookupResult{
					IP:            ip,
					IsBlacklisted: true,
					Provider:      provider,
					Reason:        desc,
				}
				return
			}

			errChan <- LookupResult{
				IP:            ip,
				IsBlacklisted: false,
				Provider:      provider,
			}
		}(provider)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func isNotCommentLine(line string) bool {
	return !strings.HasPrefix(line, "#")
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

func isNotEmptyString(str string) bool {
	return str != ""
}

func getProviders(fn string) ([]string, error) {
	f, err := os.ReadFile(fn)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(f), "\n")

	return filterString(filterString(mapString(lines, strings.TrimSpace), isNotEmptyString), isNotCommentLine), nil

}

func main() {
	var domainsFile = flag.String("p", "./providers", "path to file which stores list of dnsbl checks")
	var ipAddress = flag.String("i", "", "IP Address to check, separate by comma for a list")

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

	var ipwg sync.WaitGroup
	ipAddresses := filterString(filterString(strings.Split(*ipAddress, ","), isNotEmptyString), isValidIP)
	ipwg.Add(len(ipAddresses))
	for _, ip := range ipAddresses {
		go func(ip string) {
			defer ipwg.Done()
			for result := range getBlacklists(ip, providers) {
				if result.IsBlacklisted {
					var reason string

					if result.Reason == "" {
						reason = "unkown reason"
					} else {
						reason = result.Reason
					}

					fmt.Println(fmt.Sprintf("ERR\t%s\t%s\t%s", result.IP, result.Provider.GetName(), reason))
				} else {
					fmt.Println(fmt.Sprintf("OK\t%s\t%s", result.IP, result.Provider.GetName()))
				}
			}
		}(ip)
	}
	ipwg.Wait()
}
