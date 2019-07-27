package main

import (
	"sort"
	"strings"

	"github.com/dansackett/levenshtein"
	"github.com/dansackett/radix"
)

type suggestion struct {
	Weight           float64
	Query            string
	HasQueryAsPrefix bool
}

type bySuggestionWeight []*suggestion

func (s bySuggestionWeight) Len() int {
	return len(s)
}

func (s bySuggestionWeight) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s bySuggestionWeight) Less(i, j int) bool {
	// check weights first
	if s[i].Weight < s[j].Weight {
		return true
	}

	if s[i].Weight > s[j].Weight {
		return false
	}

	// if weights are the same, check length of the suggestion
	if len(s[i].Query) < len(s[j].Query) {
		return true
	}

	if len(s[i].Query) > len(s[j].Query) {
		return false
	}

	// if weights and length are the same, check alphabetical order
	return s[i].Query < s[j].Query
}

// Filter results so only unique items appear in slice
func unique(s []*suggestion) []*suggestion {
	keys := make(map[string]bool)
	list := []*suggestion{}
	for _, entry := range s {
		q := strings.ToLower(entry.Query)
		if _, value := keys[q]; !value {
			keys[q] = true
			list = append(list, entry)
		}
	}
	return list
}

// Find potential typos to use for further suggestions
func gatherTypoSuggestions(tree *radix.Tree, allQueries []*suggestion, searchQuery string) []*suggestion {
	for word := range tree.Iter() {
		distance := levenshtein.CalculateDistance(searchQuery, word)

		if (len(word) <= 4 && distance <= 1) || (len(word) >= 8 && distance <= 2) {
			newQueryObj := &suggestion{
				Weight:           float64(distance),
				Query:            word,
				HasQueryAsPrefix: false,
			}

			allQueries = append(allQueries, newQueryObj)
		}
	}
	return allQueries
}

// Collect all suggestions for each query
func collectSuggestions(tree *radix.Tree, allQueries []*suggestion, searchQuery string) []*suggestion {
	numQueriesToProcess := len(allQueries)
	suggestionsCollectorChannel := make(chan []*suggestion, numQueriesToProcess)

	// find suggestions for each query that we've collected
	for _, query := range allQueries {
		go func(ch chan []*suggestion, queryObj *suggestion) {
			var suggestions []*suggestion

			for _, s := range tree.GetSuggestions(queryObj.Query) {
				weight := queryObj.Weight + 0.5

				if !queryObj.HasQueryAsPrefix {
					weight = weight * float64(levenshtein.CalculateDistance(searchQuery, s))
				}

				newQueryObj := &suggestion{
					Weight:           weight,
					Query:            s,
					HasQueryAsPrefix: queryObj.HasQueryAsPrefix,
				}

				suggestions = append(suggestions, newQueryObj)
			}

			suggestionsCollectorChannel <- suggestions
		}(suggestionsCollectorChannel, query)
	}

	// wait for the suggestions to finish collecting
	for res := range suggestionsCollectorChannel {
		allQueries = append(allQueries, res...)
		numQueriesToProcess--
		if numQueriesToProcess == 0 {
			break
		}
	}

	close(suggestionsCollectorChannel)
	return allQueries
}

// Convert suggestions to a slice of strings for return
func convertSuggestionsToStrings(allQueries []*suggestion) []string {
	var rankedSuggestions []string
	for _, result := range allQueries {
		rankedSuggestions = append(rankedSuggestions, result.Query)
	}
	return rankedSuggestions
}

// GetSuggestions does the work to gather and display results for a given
// search query.
func GetSuggestions(searchQuery string) []string {
	tree := radix.InitTreeFromDict(new(radix.LinuxDictionary))

	allQueries := []*suggestion{&suggestion{
		Weight:           0.0,
		Query:            searchQuery,
		HasQueryAsPrefix: true,
	}}

	allQueries = gatherTypoSuggestions(tree, allQueries, searchQuery)
	allQueries = collectSuggestions(tree, allQueries, searchQuery)
	allQueries = unique(allQueries)

	sort.Sort(bySuggestionWeight(allQueries))

	return convertSuggestionsToStrings(allQueries)
}
