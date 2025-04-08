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
API and Data Structures
How It Works
Use Cases
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

### Levenshtein-Based Fuzzy Search
