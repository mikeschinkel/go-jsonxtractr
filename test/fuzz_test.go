package test

import (
	"bytes"
	"testing"

	"github.com/mikeschinkel/go-jsonxtractr"
)

// FuzzExtractValuesFromReader tests ExtractValuesFromReader with random JSON inputs
func FuzzExtractValuesFromReader(f *testing.F) {
	// Seed corpus with valid JSON and selectors
	seeds := []struct {
		json     string
		selector string
	}{
		{`{"name":"John"}`, "name"},
		{`{"user":{"name":"Jane"}}`, "user.name"},
		{`{"items":[1,2,3]}`, "items"},
		{`{"nested":{"deep":{"value":42}}}`, "nested.deep.value"},
		{`{}`, "missing"},
		{`{"a":1}`, "b"},
		{`null`, "key"},
		{`[]`, "key"},
		{`"string"`, "key"},
		{`123`, "key"},
		{`true`, "key"},
	}

	for _, seed := range seeds {
		f.Add(seed.json, seed.selector)
	}

	f.Fuzz(func(t *testing.T, jsonStr, selector string) {
		// Just ensure it doesn't panic
		reader := bytes.NewReader([]byte(jsonStr))
		selectors := []jsonxtractr.Selector{jsonxtractr.Selector(selector)}
		_, _, _ = jsonxtractr.ExtractValuesFromReader(reader, selectors)
	})
}

// FuzzExtractValueFromReader tests ExtractValueFromReader with random JSON inputs
func FuzzExtractValueFromReader(f *testing.F) {
	// Seed corpus
	seeds := []struct {
		json     string
		selector string
	}{
		{`{"key":"value"}`, "key"},
		{`{"a":{"b":"c"}}`, "a.b"},
		{`{"array":[1,2,3]}`, "array"},
		{`{}`, "key"},
		{`null`, "key"},
	}

	for _, seed := range seeds {
		f.Add(seed.json, seed.selector)
	}

	f.Fuzz(func(t *testing.T, jsonStr, selector string) {
		// Just ensure it doesn't panic
		reader := bytes.NewReader([]byte(jsonStr))
		_, _ = jsonxtractr.ExtractValueFromReader(reader, jsonxtractr.Selector(selector))
	})
}
