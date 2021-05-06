package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

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

	trimmed := mapStringChan(strings.TrimSpace, lines)
	noEmpty := filterStringChan(negate(isEmptyString), trimmed)
	noComment := filterStringChan(negate(isCommentLine), noEmpty)

	return noComment, nil
}

func mapStringChan(conv func(string) string, lines chan string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for line := range lines {
			out <- conv(line)
		}
	}()

	return out
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
