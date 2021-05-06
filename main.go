package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
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
	Error         error
}

func lookup(address string, provider DNSBLProvider) LookupResult {
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

func getBlacklists(addresses []string, providers []DNSBLProvider) chan LookupResult {
	var wg sync.WaitGroup
	results := make(chan LookupResult)
	for _, address := range addresses {
		for _, provider := range providers {
			wg.Add(1)
			go func(address string, provider DNSBLProvider) {
				defer wg.Done()
				results <- lookup(address, provider)
			}(address, provider)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
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

func readLines(reader io.Reader) []string {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	return text
}

func getProvidersFromFile(fn string) ([]string, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readLines(file), nil
}

func getProvidersFromStdin() ([]string, error) {
	text := readLines(os.Stdin)

	return text, nil
}

func getProviders(fn string) ([]string, error) {
	var lines []string
	var err error

	if fn == "" || fn == "-" {
		lines, err = getProvidersFromStdin()
	} else {
		lines, err = getProvidersFromFile(fn)
	}

	if err != nil {
		return nil, err
	}

	trimmed := mapString(lines, strings.TrimSpace)
	noEmpty := filterString(trimmed, negate(isEmptyString))
	noComment := filterString(noEmpty, negate(isCommentLine))

	return noComment, nil
}

func processLookupResult(result LookupResult) {
	if result.Error != nil {
		fmt.Println(fmt.Sprintf("ERR\t%s\t%s\t%s", result.Address, result.Provider.GetName(), result.Error))
		return
	}
	if result.IsBlacklisted {
		var reason string

		if result.Reason == "" {
			reason = "unkown reason"
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
	list, err := getProviders(*domainsFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading domains")
		os.Exit(1)
	}

	var providers []DNSBLProvider

	for _, item := range list {
		provider := GeneralDNSBLProvider{
			URL: item,
		}
		providers = append(providers, provider)
	}

	addresses := filterString(strings.Split(*addressesParam, ","), negate(isEmptyString))

	for result := range getBlacklists(addresses, providers) {
		processLookupResult(result)
	}
}
