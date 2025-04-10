![GitHub release](https://img.shields.io/github/v/release/Tagliapietra96/fuzzy)
![Build Status](https://github.com/Tagliapietra96/fuzzy/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/Tagliapietra96/fuzzy/path.svg)](https://pkg.go.dev/github.com/Tagliapietra96/fuzzy)
[![Go Report Card](https://goreportcard.com/badge/github.com/Tagliapietra96/fuzzy)](https://goreportcard.com/report/github.com/Tagliapietra96/fuzzy)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

# fuzzy

Welcome to fuzzy – a Go package designed to perform fuzzy searches in slices of strings, assigning a score to each match based on how well it fits. With fuzzy, you can quickly and intelligently search through collections of text, making it perfect for command-line tools, text editors, or any application that requires fast and flexible string matching. 😊

# Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage Examples](#usage-examples)
    * [Standard Fuzzy Search](#standard-fuzzy-search)
    * [Levenshtein-Based Fuzzy Search](#levenshtein-based-fuzzy-search)
5. [API and Data Structures](#api-and-data-structures)
    * [The Match Struct](#the-match-struct)
    * [Primary Functions](#primary-functions)
6. [How It Works](#how-it-works)
7. [Use Cases](#use-cases)
Inspiration
Support the Project
License

## Introduction

fuzzy is a Go package that brings the power of fuzzy matching to your applications. It analyzes a list of strings (typically representing lines, filenames, or commands) and returns the best matches along with a score—lower scores indicate a better match. With a design that balances simplicity and performance, fuzzy is suitable even for very large datasets! 🚀

## Features

* **Smart Fuzzy Search:**

    Quickly search through a slice of strings and obtain scores that reflect how closely each string matches your query.
* **Levenshtein Distance Option:**

    When you need a more forgiving search (especially useful when handling typos), leverage the Levenshtein-based search ensuring at least 60% of the query is present.
* **Customizable Query Filtering:**

    Enhance your query with filters directly within the query string:
    * `*filter` – requires the source to contain the filter.
    * `$filter` – requires the source to end with the filter.
    * `^filter` – requires the source to start with the filter.
* **Flexible Sorting:**

    Use SortMatches to arrange results first by match score and then by the position within the source.

## Installation

Install fuzzy using `go get`:

```bash
go get github.com/Tagliapietra96/fuzzy 
```

## Usage Examples

Below are some structured examples to illustrate how fuzzy can be used. Inline comments indicate the expected output for clarity.

### Standard Fuzzy Search

The `Find` function conducts a regular fuzzy search and computes a score based on factors such as substring presence and character distance.

```go
package main

import "github.com/Tagliapietra96/fuzzy" 

func main() {
    // List of strings to search through
    data := []string{"cart", "clap", "ca", "cat", "cow"}

    // String to use as query
    input := "ca"

    // Execute a fuzzy search using the query "ca"
    // Expected behavior:
    //  - "cart" -> Match{Score: 2, Position: 0} - The query is a substring of the line, the score is the differce of lenght
    //  - "clap" -> Match{Score: 3, Position: 1} - The query isn't a substring of the line, the score is the differce of lenght + the distance between every rune of the query founded in the line
    //  - "ca" -> Match{Sycore: 0, Position: 2} - The score is 0 (perfect match)
    //  - "cat" -> Match{Score:1, Position: 3}
    //  - "cow" -> Not a match! the line must have all the runes of the query 
    matches := fuzzy.Find(input, data)

    // Sort the result score
    //  - "ca" -> Match{Sycore: 0, Position: 2}
    //  - "cat" -> Match{Score:1, Position: 3}
    //  - "cart" -> Match{Score: 2, Position: 0}
    //  - "clap" -> Match{Score: 3, Position: 1}
    sortedMatches := fuzzy.SortMatches(matches)

    // Execute a fuzzy search using the query "ca" and filtering onli the lines that end with "t"
    // Expected behavior:
    //  - "cart" -> Match{Score: 2, Position: 0} - The query is a substring of the line, the score is the differce of lenght
    //  - "clap" -> Match{Score: 3, Position: 1} - Not a match! The line don't end with "t"
    //  - "ca" -> Match{Sycore: 0, Position: 2} - Not a match! The line don't end with "t"
    //  - "cat" -> Match{Score:1, Position: 3}
    //  - "cow" -> Not a match! the line must have all the runes of the query and end with "t"
    matchesWithFilter := fuzzy.Find("$t ca", data)
}
```

### Levenshtein-Based Fuzzy Search

The `LevenshteinFind` function applies the Levenshtein distance to evaluate match quality, tolerating minor typos or character errors.

```go
package main

import "github.com/Tagliapietra96/fuzzy" 

func main() {
    // List of strings to search through
    data := []string{"cart", "clap", "ca", "cat", "cow"}

    // String to use as query
    input := "caw"

    // Execute a fuzzy search using the query "caw"
    // Expected behavior:
    //  - "cart" -> Match{Score: 2, Position: 0}
    //  - "clap" -> Match{Score: 3, Position: 1}
    //  - "ca" -> Match{Sycore: 1, Position: 2}
    //  - "cat" -> Match{Score:1, Position: 3}
    //  - "cow" -> Match{Score: 1, Position: 4}
    matches := fuzzy.LevenshteinFind(input, data)

    // Sort the result score
    //  - "ca" -> Match{Sycore: 1, Position: 2}
    //  - "cat" -> Match{Score:1, Position: 3}
    //  - "cow" -> Match{Score: 1, Position: 4}
    //  - "cart" -> Match{Score: 2, Position: 0}
    //  - "clap" -> Match{Score: 3, Position: 1}
    sortedMatches := fuzzy.SortMatches(matches)

    // Execute a fuzzy search using the query "caw" and filtering onli the lines that end with "t"
    // Expected behavior:
    //  - "cart" -> Match{Score: 2, Position: 0} 
    //  - "clap" -> Match{Score: 3, Position: 1} - Not a match! The line don't end with "t"
    //  - "ca" -> Match{Sycore: 1, Position: 2} - Not a match! The line don't end with "t"
    //  - "cat" -> Match{Score:1, Position: 3}
    //  - "cow" -> Match{Score: 1, Position: 4} - Not a match! the line don't and end with "t"
    matchesWithFilter := fuzzy.LevenshteinFind("$t ca", data)
}
```

## API and Data Structures

### The `Match` Struct
The `Match` structure represents a search result:

```go
// Match contains the score of the match and its position within the source slice.
type Match struct {
    Score    int // Lower is better
    Position int // Index of the matching element in the source
}
```

### Primary Functions

* `Find(queryValue string, source []string) []Match`

    Searches for the query in the provided slice and returns all matching entries along with their scores.
    Tip: By default, the search is case-insensitive unless the query is capitalized.
* `LevenshteinFind(queryValue string, source []string) []Match`

    Uses the Levenshtein distance for a more flexible, approximate matching, ideal for handling typos. A match requires at least 60% similarity with the query.
* `SortMatches(m []Match) []Match`

    Orders the matches—first by score (ascending) and then by source position if scores are equal.

## How It Works

1. **Query Parsing & Filtering:**

    The package splits your query into the main string and any additional filters. Filters can enforce where the matching string should begin, end, or simply contain a given substring.

2. **Scoring:**

    * Exact or Substring Match: If the query exists as an exact match or substring, the score is determined by the length difference between the source string and the query.
    * Inexact Match: When characters are mismatched or out of order, extra penalties (i.e., distance) are added to the score.
    * Levenshtein Mode: Calculates the edit distance in a manner that is forgiving of typos, as long as at least 60% of the query is present.

3. **Sorting:**

    Once matches are identified, the SortMatches function arranges them so that the best matches (with the lowest scores) are shown first.

## Use Cases

**fuzzy** can be integrated into various applications:
* **CLI Tools:** Enhance command-line interfaces with real-time suggestions and autocompletion.
* **Text Editors and IDEs:** Implement smart search functionalities for files, commands, or code snippets.
* **Data Filtering Systems:** Quickly process and filter large datasets while handling typos or partial matches gracefully.
