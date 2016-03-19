package main

import (
  "os"
  "fmt"
  "path/filepath"
  "github.com/codegangsta/cli"
)

// Main worker for iterating files.
func outputResults(base_path string, result Result, matcher FileMatcher, c *cli.Context) {
  invert := c.Bool("invert")
  full_path_match := c.Bool("full-path")
  var match_path string

  include_dirs := c.Bool("include-dirs")
  dirs_only := c.Bool("dirs-only")
  // Small shortcut to skip dir checks early.
  no_dir_filters := !(include_dirs || dirs_only)

  // Allow the result to print something before files are walked.
  result.beforeResults()

  err := filepath.Walk(base_path, func (path string, fileInfo os.FileInfo, err error) error {
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    is_dir := fileInfo.IsDir();

    // Determine how to deal with dirs based on any flags.
    if (no_dir_filters && is_dir) {
      // If there are no dir filters and this is dir, return. Otherwise, this
      // will try to use the matcher below. By default, dirs are not included in
      // results.
      return nil
    // This means there are dir filters, so check against them. Don't check
    // include_dirs as if it hasn't been caught above (no_dir_filters) we want
    // everything to go to below matching.
    } else if (!is_dir && dirs_only) {
      return nil
    }

    // If this is a full path match use the path as-is. Otherwise use the file
    // name only (default).
    if full_path_match {
      match_path = path
    } else {
      match_path = filepath.Base(path);
    }

    matched := matcher.match(match_path)

    if matched != invert {
      result.eachResult(path)
    }

    return nil
  })

  // Allow the result to print something after the files have been walked.
  result.afterResults()

  if err != nil {
    fmt.Println(err)
  }
}

// Determine the base path and pattern for searching based on args.
func determineArgs(c *cli.Context) (string, string) {
  var base_path, pattern string
  var err error

  args := c.Args()
  length := c.NArg();

  switch {
    case length == 1:
      pattern = args.First()
    case length > 1:
      // Assume path will be the first arg, and pattern the second.
      base_path = args.First()
      pattern = args.Get(1)
  }

  // If the base path is still empty, get the cwd.
  if base_path == "" {
    // Default root to cwd.
    base_path, err = os.Getwd()

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
  }

  return base_path, pattern
}

// Choose an appropriate matcher.
// Based on whether the pattern is empty, or extended regex is needed.
func createMatcher(pattern string, c *cli.Context) FileMatcher {
  if pattern == "" {
    return newEmptyMatcher()
  }

  insensitive := c.Bool("insensitive")

  if c.Bool("ext") {
    return newRegexMatcher(pattern, insensitive)
  }

  return newFilepathMatcher(pattern, insensitive)
}

// Choose an appropriate result.
// This is based on whether or not a count is needed.
func createResult(c *cli.Context) Result {
  if c.Bool("count") {
    return &CountResult{}
  }

  return &PrintResult{}
}
