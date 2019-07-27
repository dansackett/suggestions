# Suggestions

Suggestion is a toy project to see how a basic autocompletion / spell checker
would work. It uses a [Radix Tree](https://github.com/dansackett/radix) and the
[Levenshtein Distance](https://github.com/dansackett/levenshtein) formula to
calcuate and return a slice of strings which may be a match or has potential to
be used as an autocomplete term.

Usage:

```
go build

Usage of ./suggestions:
  -num-results int
        Number of results to return (default 10)
  -query string
        Query to get suggestions for
```
