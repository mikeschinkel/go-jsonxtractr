package test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mikeschinkel/go-jsonxtractr"
)

// TestFuzzCorpus reads each fuzz corpus file and tests it with timeout detection
func TestFuzzCorpus(t *testing.T) {
	// Test ExtractValuesFromReader corpus
	testCorpus(t, "FuzzExtractValuesFromReader", func(input1, input2 string) error {
		reader := bytes.NewReader([]byte(input1))
		selectors := []jsonxtractr.Selector{jsonxtractr.Selector(input2)}
		_, _, _ = jsonxtractr.ExtractValuesFromReader(reader, selectors)
		return nil
	})

	// Test ExtractValueFromReader corpus
	testCorpus(t, "FuzzExtractValueFromReader", func(input1, input2 string) error {
		reader := bytes.NewReader([]byte(input1))
		_, _ = jsonxtractr.ExtractValueFromReader(reader, jsonxtractr.Selector(input2))
		return nil
	})
}

func testCorpus(t *testing.T, fuzzFuncName string, testFunc func(string, string) error) {
	corpusDir := filepath.Join("testdata", "fuzz", fuzzFuncName)
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		// No corpus directory is fine - fuzzing hasn't been run yet
		return
	}

	infiniteLoops := []string{}
	parseErrors := []string{}
	successes := []string{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Read the fuzz corpus file
		path := filepath.Join(corpusDir, entry.Name())
		f, err := os.Open(path)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(f)
		lineNum := 0
		var input1, input2 string
		for scanner.Scan() {
			lineNum++
			if lineNum == 2 { // Second line contains first string
				line := scanner.Text()
				if strings.HasPrefix(line, "string(") && strings.HasSuffix(line, ")") {
					strLiteral := line[7 : len(line)-1]
					unquoted, err := strconv.Unquote(strLiteral)
					if err != nil {
						break
					}
					input1 = unquoted
				}
			} else if lineNum == 3 { // Third line contains second string
				line := scanner.Text()
				if strings.HasPrefix(line, "string(") && strings.HasSuffix(line, ")") {
					strLiteral := line[7 : len(line)-1]
					unquoted, err := strconv.Unquote(strLiteral)
					if err != nil {
						break
					}
					input2 = unquoted
				}
				break
			}
		}
		_ = f.Close()

		if input1 == "" && input2 == "" {
			continue
		}

		// Test this input with timeout
		done := make(chan struct{})
		var testErr error

		go func() {
			defer func() {
				if r := recover(); r != nil {
					testErr = fmt.Errorf("PANIC: %v", r)
				}
				close(done)
			}()

			testErr = testFunc(input1, input2)
		}()

		select {
		case <-done:
			// Test completed
			if testErr != nil {
				parseErrors = append(parseErrors, entry.Name())
			} else {
				successes = append(successes, entry.Name())
			}
		case <-time.After(10 * time.Second):
			infiniteLoops = append(infiniteLoops, entry.Name())
			t.Errorf("%-20s INFINITE LOOP: json=%q selector=%q", entry.Name(), input1, input2)
		}
	}

	if len(infiniteLoops) > 0 {
		t.Fatalf("Found %d infinite loop(s) in %s", len(infiniteLoops), fuzzFuncName)
	}
}
