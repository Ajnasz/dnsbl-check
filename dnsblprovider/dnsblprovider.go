package dnsblprovider

// DNSBLProvider interface should be implemented to be able to query a provider
type DNSBLProvider interface {
	GetName() string
	IsBlacklisted(string) (bool, error)
	GetReason(string) (string, error)
}
