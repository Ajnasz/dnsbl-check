package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/Ajnasz/dnsbl-check/stringutils"
)

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

func getProvidersChan(fn string) (chan string, error) {
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
	noEmpty := filterStringChan(negate(isEmptyString), trimmed)
	noComment := filterStringChan(negate(isCommentLine), noEmpty)

	return noComment, nil
}

func filterStringChan(test func(string) bool, lines chan string) chan string {
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
