package providerlist

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/Ajnasz/dnsbl-check/stringutils"
)

func negate(f func(string) bool) func(string) bool {
	return func(str string) bool {
		r := f(str)
		return !r
	}
}

func isCommentLine(line string) bool {
	return strings.HasPrefix(line, "#")
}

func isEmptyString(str string) bool {
	return str == ""
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

	trimmed := stringutils.Map(lines, strings.TrimSpace)
	noEmpty := stringutils.Filter(trimmed, negate(isEmptyString))
	noComment := stringutils.Filter(noEmpty, negate(isCommentLine))

	return noComment, nil
}

func getLinesChan(reader io.Reader) (chan string, chan struct{}) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	text := make(chan string)
	done := make(chan struct{})

	go func() {
		defer func() {
			close(text)
			done <- struct{}{}
			close(done)
		}()
		for scanner.Scan() {
			text <- scanner.Text()
		}
	}()

	return text, done
}

func getProvidersFromFileChan(fn string) (chan string, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	text, done := getLinesChan(file)

	go func() {
		defer file.Close()
		<-done
	}()

	return text, nil
}

func getProvidersFromStdinChan() (chan string, error) {
	text, done := getLinesChan(os.Stdin)

	go func() {
		<-done
	}()

	return text, nil
}

func GetProvidersChan(fn string) (chan string, error) {
	var lines chan string
	var err error

	if fn == "" || fn == "-" {
		lines, err = getProvidersFromStdinChan()
	} else {
		lines, err = getProvidersFromFileChan(fn)
	}

	if err != nil {
		return nil, err
	}

	trimmed := stringutils.MapChan(strings.TrimSpace, lines)
	noEmpty := stringutils.FilterChan(negate(isEmptyString), trimmed)
	noComment := stringutils.FilterChan(negate(isCommentLine), noEmpty)

	return noComment, nil
}

func GetAddresses(addresses string) []string {
	return stringutils.Filter(strings.Split(addresses, ","), negate(isEmptyString))
}
