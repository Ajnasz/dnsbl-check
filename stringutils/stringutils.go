package stringutils

// Filter filters the slice of string by the given test func and keeps only
// those entries where the test returned true
func Filter(lines []string, test func(string) bool) []string {
	var out []string

	for _, line := range lines {
		if test(line) {
			out = append(out, line)
		}
	}

	return out
}

// Map returns the output conv function which executed on all entires
func Map(lines []string, conv func(string) string) []string {
	var out []string

	for _, line := range lines {
		out = append(out, conv(line))
	}

	return out
}

// MapChan enables to alter the input strings and emit them from a channel
func MapChan(conv func(string) string, lines chan string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for line := range lines {
			out <- conv(line)
		}
	}()

	return out
}

// FilterChan reads items from a string channel, and emits only those where the
// test function returned true
func FilterChan(test func(string) bool, lines chan string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for line := range lines {
			if test(line) {
				out <- line
			}
		}
	}()

	return out
}
