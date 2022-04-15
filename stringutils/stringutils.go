package stringutils

func Filter(lines []string, test func(string) bool) []string {
	var out []string

	for _, line := range lines {
		if test(line) {
			out = append(out, line)
		}
	}

	return out
}

func Map(lines []string, conv func(string) string) []string {
	var out []string

	for _, line := range lines {
		out = append(out, conv(line))
	}

	return out
}

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
