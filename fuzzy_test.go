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
