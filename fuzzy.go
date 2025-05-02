package fuzzy

import (
	"runtime"
	"slices"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// ChunkFind performs a parallelized fuzzy search using the standard matching algorithm.
// It splits the source slice into chunks and processes them concurrently for better performance
// on large datasets, then combines the results.
func ChunkFind(query string, source []string) []Match {
	return chunkFind(query, source, Find)
}

// ChunkLevenshteinFind performs a parallelized fuzzy search using the Levenshtein distance algorithm.
// It splits the source slice into chunks and processes them concurrently for better performance
// on large datasets, then combines the results.
func ChunkLevenshteinFind(query string, source []string) []Match {
	return chunkFind(query, source, LevenshteinFind)
}

// chunkFind is a helper function that splits the source into chunks and runs the algorithm on each chunk.
func chunkFind(query string, source []string, algo func(q string, s []string) []Match) []Match {
	cpu := min(4, runtime.NumCPU()/2)

	if cpu <= 1 || len(source) <= cpu*500 {
		return algo(query, source)
	}

	var wg sync.WaitGroup
	cs := len(source) / cpu
	cc := (len(source) + cs - 1) / cs
	rChan := make(chan []Match, cc)

	wg.Add(cc)
	for i := range cc {
		go func(chunk []string) {
			defer wg.Done()
			mm := algo(query, chunk)
			for j := range mm {
				mm[j].Position += i * cs
			}
			rChan <- mm
		}(source[i*cs : min(((i*cs)+cs), len(source))])
	}

	go func() {
		wg.Wait()
		close(rChan)
	}()

	r := make([]Match, 0, len(source))
	for mm := range rChan {
		r = append(r, mm...)
	}

	return r
}

// SortMatches sorts the matches by score and position.
//   - if the scores are equal, the position is used to determine the order
//   - if the scores are different, the score is used to determine the order
//
// Lower the score, better the match.
func SortMatches(m []Match) []Match {
	slices.SortFunc(m, func(a, b Match) int {
		if a.Score == b.Score {
			return a.Position - b.Position
		}
		return a.Score - b.Score
	})
	return m
}

// Find searches for the value in the source and returns the matches.
// The result is a slice of Match structs, which contain the score (lower is better)
// and the position of the match in the source (e.g. source[Match.Position]).
//
// By default, the search is case insensitive, but it can be made case sensitive by capitalizing the query.
// This behavior isn't affected by the filters.
//
// It it possible to add filters to the query (value), each filter must be separated by a space:
//   - if the filter starts with *, the source line must contain the filter
//   - if the filter starts with $, the source line must end with the filter
//   - if the filter starts with ^, the source line must start with the filter
//
// e.g. "query *filter1 $filter2 ^filter3" or "*filter1 $filter2 query ^filter3"
//
// The score are calculated in the following way:
//   - if the query is equal to the value, the score is 0
//   - if the query is a substring of the value, the score is the length of the value minus the length of the query
//   - else the score is the length of the value minus the length of the query plus the distance between every character of the query in the value
//
// e.g. "ca" in "cart" has a score of 2, "ca" in "clap" has a score of 3
//
// The result is unsorted.
// If you want to sort the result, use the SortMatches function.
func Find(queryValue string, source []string) []Match {
	return find(queryValue, source, matchScore)
}

// LevenshteinFind acts the same as Find, but it uses the Levenshtein distance to calculate the score.
// In this case the matches are more approximate, in fact to have a match the source line must contain at least 60% of the query.
// This is useful when the query is misspelled or when the source contains typos.
func LevenshteinFind(queryValue string, source []string) []Match {
	c := make([]int, len(queryValue)+1, len(queryValue)+1)
	return find(queryValue, source, func(q, s string, f func(string) (string, bool)) int {
		return levenshteinScore(q, s, f, c)
	})
}

// Match is a struct that contains the score and the position (in the source slice) of the match.
type Match struct {
	Score    int
	Position int
}

// find searches for the query in the source and returns the matches.
func find(q string, s []string, fn func(string, string, func(string) (string, bool)) int) []Match {
	var f func(string) (string, bool)
	q, f = input(q)
	m := make([]Match, 0, len(s))

	for i, l := range s {
		score := fn(q, l, f)
		if score >= 0 {
			m = append(m, Match{Score: score, Position: i})
		}
	}

	return m
}

// matchScore calculates the score of the match.
func matchScore(q, s string, f func(string) (string, bool)) int {
	var found bool
	s, found = f(s)
	ql, sl := len(q), len(s)

	// preliminary check to optimize algorithm speed
	switch {
	case !found || sl < ql:
		return -1
	case q == s, q == "":
		return 0
	case strings.Contains(s, q):
		return sl - ql
	}

	distance := 0
Outer:
	for index, qr := range q {
		for i, sr := range s {
			if qr == sr {
				s = s[i+utf8.RuneLen(sr):]
				if index > 0 {
					distance += i
				}
				continue Outer
			}
		}
		return -1
	}

	return sl - ql + distance
}

// levenshteinScore calculates the score of the match using the Levenshtein distance.
func levenshteinScore(q, s string, f func(string) (string, bool), column []int) int {
	var found bool
	s, found = f(s)
	ql, sl := len(q), len(s)

	// preliminary check to optimize algorithm speed
	switch {
	case !found || sl < ql:
		return -1
	case q == s, q == "":
		return 0
	case strings.Contains(s, q):
		return sl - ql
	}

	founded := 0
	minFind := int(float64(ql) * 0.6)

	sc := s
Outer:
	for _, qr := range q {
		if founded >= minFind {
			break
		}
		for i, sr := range sc {
			if qr == sr {
				sc = sc[i+utf8.RuneLen(sr):]
				founded++
				continue Outer
			}
		}
		return -1
	}

	for i := 0; i <= ql; i++ {
		column[i] = i
	}

	for x := 1; x <= sl; x++ {
		column[0] = x
		lastDiag := x - 1

		for y := 1; y <= ql; y++ {
			oldDiag := column[y]
			var cost int
			if q[y-1] != s[x-1] {
				cost = 1
			}
			column[y] = min(
				column[y]+1,
				column[y-1]+1,
				lastDiag+cost,
			)
			lastDiag = oldDiag
		}
	}

	return column[ql]
}

// input returns the query and the filter function.
func input(q string) (string, func(string) (string, bool)) {
	if q == "" {
		return "", func(s string) (string, bool) {
			return s, true
		}
	}

	f := make([]string, 0)
	b := &strings.Builder{}

	for w := range strings.SplitSeq(q, " ") {
		if strings.HasPrefix(w, "*") || strings.HasPrefix(w, "$") || strings.HasPrefix(w, "^") {
			f = append(f, w)
		} else {
			b.WriteString(w)
		}
	}

	upper := isUpper(b.String())
	if len(f) == 0 {
		return b.String(), func(s string) (string, bool) {
			if !upper {
				s = strings.ToLower(s)
			}

			return removeWhitespace(s), true
		}
	}

	return b.String(), func(s string) (string, bool) {
		if !upper {
			s = strings.ToLower(s)
		}

		found := true
		for _, fv := range f {
			if !found {
				return "", false
			}
			switch {
			case strings.HasPrefix(fv, "*"):
				b, a, fo := strings.Cut(s, fv[1:])
				s, found = b+a, fo
			case strings.HasPrefix(fv, "$"):
				s, found = strings.CutSuffix(s, fv[1:])
			case strings.HasPrefix(fv, "^"):
				s, found = strings.CutPrefix(s, fv[1:])
			}
		}

		return removeWhitespace(s), found
	}
}

// removeWhitespace removes the whitespace from the string.
func removeWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

// isUpper checks if the word is capitalized.
func isUpper(w string) bool {
	if len(w) == 0 {
		return false
	}
	for _, r := range w {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}
