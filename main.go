package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	searchQuery := flag.String("query", "", "Query to get suggestions for")
	numResults := flag.Int("num-results", 10, "Number of results to return")

	flag.Parse()

	if *searchQuery == "" {
		log.Fatalf("Must provide a query parameter")
	}

	rankedSuggestions := GetSuggestions(*searchQuery)

	if len(rankedSuggestions) <= *numResults {
		fmt.Println(rankedSuggestions)
	} else {
		fmt.Println(rankedSuggestions[:*numResults])
	}
}
