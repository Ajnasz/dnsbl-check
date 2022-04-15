package iputils

import "strings"

func reverseStringByToken(str string, delimiter string) string {
	if str == "" {
		return ""
	}

	before, after, found := strings.Cut(str, delimiter)

	if !found {
		return before
	}

	return reverseStringByToken(after, delimiter) + delimiter + before
}

// ReverseIPv4Address will return ip address octets in reverse order:
// 1.2.3.4 will be 4.3.2.1
func ReverseIPv4Address(str string) string {
	return reverseStringByToken(str, ".")
}
