package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/mikeschinkel/go-jsonxtractr"
)

func main() {
	// Example JSON data
	jsonData := `{
		"user": {
			"name": "John Doe",
			"email": "john@example.com",
			"age": 30
		},
		"address": {
			"city": "San Francisco",
			"country": "USA"
		}
	}`

	// Create a reader from the JSON data
	reader := bytes.NewReader([]byte(jsonData))

	// Define selectors for values we want to extract
	selectors := []jsonxtractr.Selector{
		"user.name",
		"user.email",
		"address.city",
	}

	// Extract values
	values, notFound, err := jsonxtractr.ExtractValuesFromReader(reader, selectors)
	if err != nil {
		log.Fatalf("Error extracting values: %v", err)
	}

	// Display extracted values
	fmt.Println("Extracted values:")
	for selector, value := range values {
		fmt.Printf("  %s: %v\n", selector, value)
	}

	// Display not found selectors
	if len(notFound) > 0 {
		fmt.Println("\nNot found:")
		for _, selector := range notFound {
			fmt.Printf("  %s\n", selector)
		}
	}

	// Example of extracting a single value
	reader2 := bytes.NewReader([]byte(jsonData))
	userName, err := jsonxtractr.ExtractValueFromReader(reader2, "user.name")
	if err != nil {
		log.Fatalf("Error extracting single value: %v", err)
	}
	fmt.Printf("\nSingle value extraction:\n  user.name: %v\n", userName)
}
