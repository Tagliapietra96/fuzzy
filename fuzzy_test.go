package fuzzy

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestSortMatches(t *testing.T) {
	testCases := []struct {
		name     string
		input    []Match
		expected []Match
	}{
		{
			name: "Sort matches with different scores",
			input: []Match{
				{Score: 3, Position: 1},
				{Score: 1, Position: 0},
				{Score: 2, Position: 2},
			},
			expected: []Match{
				{Score: 1, Position: 0},
				{Score: 2, Position: 2},
				{Score: 3, Position: 1},
			},
		},
		{
			name: "Sort matches with equal scores",
			input: []Match{
				{Score: 2, Position: 3},
				{Score: 2, Position: 1},
				{Score: 2, Position: 2},
			},
			expected: []Match{
				{Score: 2, Position: 1},
				{Score: 2, Position: 2},
				{Score: 2, Position: 3},
			},
		},
		{
			name:     "Empty slice",
			input:    []Match{},
			expected: []Match{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SortMatches(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestInput(t *testing.T) {
	testCases := []struct {
		name           string
		query          string
		expectedQuery  string
		line           string
		expectedResult string
		expectFound    bool
	}{
		{
			name:           "Empty query",
			query:          "",
			expectedQuery:  "",
			line:           "hello",
			expectedResult: "hello",
			expectFound:    true,
		},
		{
			name:           "Simple lowercase query",
			query:          "hello",
			expectedQuery:  "hello",
			line:           "hello",
			expectedResult: "hello",
			expectFound:    true,
		},
		{
			name:           "Query with contains modifier",
			query:          "*world",
			expectedQuery:  "",
			line:           "helloworld test",
			expectedResult: "hellotest",
			expectFound:    true,
		},
		{
			name:           "Query with ends with modifier",
			query:          "$test",
			expectedQuery:  "",
			line:           "hello test",
			expectedResult: "hello",
			expectFound:    true,
		},
		{
			name:           "Query with starts with modifier",
			query:          "^hello",
			expectedQuery:  "",
			line:           "hello world",
			expectedResult: "world",
			expectFound:    true,
		},
		{
			name:           "Multiple modifiers",
			query:          "^hello *world $test",
			expectedQuery:  "",
			line:           "hello big world test",
			expectedResult: "big",
			expectFound:    true,
		},
		{
			name:           "Contains modifier not found",
			query:          "*xyz",
			expectedQuery:  "",
			line:           "hello",
			expectedResult: "hello",
			expectFound:    false,
		},
		{
			name:           "Regex modifier match",
			query:          "?\\w+@\\w+\\.\\w+",
			expectedQuery:  "",
			line:           "contact me at test@example.com please",
			expectedResult: "contactmeattest@example.complease",
			expectFound:    true,
		},
		{
			name:           "Regex modifier no match",
			query:          "?\\d{3}-\\d{2}-\\d{4}",
			expectedQuery:  "",
			line:           "no ssn here",
			expectedResult: "nossnhere",
			expectFound:    false,
		},
		{
			name:           "Negated contains modifier",
			query:          "!*xyz",
			expectedQuery:  "",
			line:           "hello world",
			expectedResult: "helloworld",
			expectFound:    true,
		},
		{
			name:           "Negated starts with modifier",
			query:          "!^hello",
			expectedQuery:  "",
			line:           "world hello",
			expectedResult: "worldhello",
			expectFound:    true,
		},
		{
			name:           "Negated ends with modifier",
			query:          "!$test",
			expectedQuery:  "",
			line:           "test hello",
			expectedResult: "testhello",
			expectFound:    true,
		},
		{
			name:           "Negated regex modifier",
			query:          "!?\\d+",
			expectedQuery:  "",
			line:           "only text here",
			expectedResult: "onlytexthere",
			expectFound:    true,
		},
		{
			name:           "Complex query with negation and regex",
			query:          "search !*avoid ?\\w+@\\w+\\.\\w+",
			expectedQuery:  "search",
			line:           "search for email@example.com",
			expectedResult: "searchforemail@example.com",
			expectFound:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, filterFunc := input(tc.query)

			// Verify the query part
			if query != tc.expectedQuery {
				t.Errorf("Expected query '%s', got '%s'", tc.expectedQuery, query)
			}

			// Test the filter function
			result, found := filterFunc(tc.line)
			if found != tc.expectFound {
				t.Errorf("Expected found=%v, got found=%v", tc.expectFound, found)
			}

			if result != removeWhitespace(tc.expectedResult) {
				t.Errorf("Expected result '%s', got '%s'",
					removeWhitespace(tc.expectedResult), result)
			}
		})
	}
}

func TestRemoveWhitespace(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"hello world", "helloworld"},
		{"  spaces  everywhere  ", "spaceseverywhere"},
		{"no spaces", "nospaces"},
		{"", ""},
		{"   ", ""},
		{"multiple   spaces between", "multiplespacesbetween"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := removeWhitespace(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestIsUpper(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"hello", false},
		{"Hello", true},
		{"hEllo", true},
		{"HELLO", true},
		{"123", false},
		{"Hello123", true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := isUpper(tc.input)
			if result != tc.expected {
				t.Errorf("For input '%s', expected %v, got %v",
					tc.input, tc.expected, result)
			}
		})
	}
}

func TestFind(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		source   []string
		expected []Match
	}{
		{
			name:     "Empty query",
			query:    "",
			source:   []string{"test", "example"},
			expected: []Match{{Score: 0, Position: 0}, {Score: 0, Position: 1}},
		},
		{
			name:     "Simple substring match",
			query:    "test",
			source:   []string{"testing", "example"},
			expected: []Match{{Score: 3, Position: 0}},
		},
		{
			name:     "Case insensitive match",
			query:    "TEST",
			source:   []string{"testing", "Example"},
			expected: []Match{},
		},
		{
			name:     "Filter with contains",
			query:    "*this test",
			source:   []string{"this is a test", "another test"},
			expected: []Match{{Score: 3, Position: 0}},
		},
		{
			name:     "Filter with regex",
			query:    "?\\w+@\\w+\\.\\w+ email",
			source:   []string{"contact: email@example.com", "no email here"},
			expected: []Match{{Score: 20, Position: 0}},
		},
		{
			name:     "Filter with negation",
			query:    "!*another test",
			source:   []string{"this is a test", "another example"},
			expected: []Match{{Score: 7, Position: 0}},
		},
		{
			name:     "Complex query with negation and regex",
			query:    "test !*another ?\\w+",
			source:   []string{"test something", "another test"},
			expected: []Match{{Score: 9, Position: 0}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Find(tc.query, tc.source)

			sort.Slice(result, func(i, j int) bool {
				if result[i].Score != result[j].Score {
					return result[i].Score < result[j].Score
				}
				return result[i].Position < result[j].Position
			})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestMatchScore(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		source   string
		expected int
	}{
		{
			name:     "Empty query",
			query:    "",
			source:   "test",
			expected: 0,
		},
		{
			name:     "Exact match",
			query:    "test",
			source:   "test",
			expected: 0,
		},
		{
			name:     "Substring match",
			query:    "test",
			source:   "testing",
			expected: 3,
		},
		{
			name:     "Fuzzy match",
			query:    "tst",
			source:   "test",
			expected: 2,
		},
		{
			name:     "Case insensitive match",
			query:    "test",
			source:   "TEST",
			expected: 0,
		},
		{
			name:     "Case sensitive match",
			query:    "Test",
			source:   "test",
			expected: -1,
		},
		{
			name:     "No match",
			query:    "xyz",
			source:   "test",
			expected: -1,
		},
		{
			name:     "Query with filter - success",
			query:    "*est i",
			source:   "testing",
			expected: 3,
		},
		{
			name:     "Query with filter - fail",
			query:    "*xyz",
			source:   "testing",
			expected: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MatchScore(tc.query, tc.source)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestLevenshteinScore(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		source   string
		expected int
	}{
		{
			name:     "Empty query",
			query:    "",
			source:   "test",
			expected: 0,
		},
		{
			name:     "Exact match",
			query:    "test",
			source:   "test",
			expected: 0,
		},
		{
			name:     "One character difference",
			query:    "tset",
			source:   "test",
			expected: 2,
		},
		{
			name:     "Two character difference",
			query:    "tast",
			source:   "tent",
			expected: -1,
		},
		{
			name:     "Case insensitive match",
			query:    "test",
			source:   "TEST",
			expected: 0,
		},
		{
			name:     "Substring presence",
			query:    "test",
			source:   "testing",
			expected: 3,
		},
		{
			name:     "Query with filter - success",
			query:    "*est i",
			source:   "testing",
			expected: 3,
		},
		{
			name:     "Query with filter - fail",
			query:    "*xyz",
			source:   "testing",
			expected: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LevenshteinScore(tc.query, tc.source)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestLevenshteinFind(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		source   []string
		expected []Match
	}{
		{
			name:     "Misspelled query",
			query:    "tset",
			source:   []string{"test", "example"},
			expected: []Match{{Score: 2, Position: 0}},
		},
		{
			name:     "Partial match",
			query:    "tes",
			source:   []string{"testing", "example"},
			expected: []Match{{Score: 4, Position: 0}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LevenshteinFind(tc.query, tc.source)

			sort.Slice(result, func(i, j int) bool {
				if result[i].Score != result[j].Score {
					return result[i].Score < result[j].Score
				}
				return result[i].Position < result[j].Position
			})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func BenchmarkChunkFind(b *testing.B) {
	source := make([]string, 10000)
	for i := range source {
		source[i] = fmt.Sprintf("test%d", i)
	}

	for b.Loop() {
		ChunkFind("test", source)
	}
}

func BenchmarkChunkLevenshteinFind(b *testing.B) {
	source := make([]string, 10000)
	for i := range source {
		source[i] = fmt.Sprintf("test%d", i)
	}

	for b.Loop() {
		ChunkLevenshteinFind("tset", source)
	}
}

func BenchmarkFind(b *testing.B) {
	source := make([]string, 10000)
	for i := range source {
		source[i] = fmt.Sprintf("test%d", i)
	}

	for b.Loop() {
		Find("test", source)
	}
}

func BenchmarkLevenshteinFind(b *testing.B) {
	source := make([]string, 10000)
	for i := range source {
		source[i] = fmt.Sprintf("test%d", i)
	}

	for b.Loop() {
		LevenshteinFind("tset", source)
	}
}

func BenchmarkMatchScore(b *testing.B) {
	testCases := []struct {
		name   string
		query  string
		source string
	}{
		{"Exact match", "test", "test"},
		{"Substring match", "test", "testing"},
		{"Fuzzy match", "tst", "testing"},
		{"Query with filter", "*est i", "testing"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				MatchScore(tc.query, tc.source)
			}
		})
	}
}

func BenchmarkLevenshteinScore(b *testing.B) {
	testCases := []struct {
		name   string
		query  string
		source string
	}{
		{"Exact match", "test", "test"},
		{"One character difference", "tast", "test"},
		{"Fuzzy match", "tst", "testing"},
		{"Query with filter", "*est i", "testing"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				LevenshteinScore(tc.query, tc.source)
			}
		})
	}
}

func BenchmarkRemoveWhitespace(b *testing.B) {
	testString := "  multiple   spaces   everywhere  "

	for b.Loop() {
		removeWhitespace(testString)
	}
}

func BenchmarkIsUpper(b *testing.B) {
	testString := "HelloWorld123"

	for b.Loop() {
		isUpper(testString)
	}
}
