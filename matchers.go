package main

import (
  "fmt"
  "os"
  "regexp"
  "path/filepath"
  "strings"
)

type FileMatcher interface {
  match(path string) (bool)
}

type BaseMatcher struct {
  pattern string
  insensitive bool
}

type RegexMatcher struct {
  regexp *regexp.Regexp
  BaseMatcher
}

func (m RegexMatcher) match(path string) (bool) {
  return m.regexp.MatchString(path)
}

func newRegexMatcher(pattern string, insensitive bool) RegexMatcher {
  // If this is case-insensitive, prepend '(?i)' to the regex. Likely NOT the
  // best way to do this!
  if insensitive {
    // This assignment is not ideal, but easy.
    prefix := "(?i)"
    pattern = prefix + pattern
  }

  regexp, err := regexp.Compile(pattern)

  if err != nil {
      fmt.Println(err)
      os.Exit(1)
  }

  return RegexMatcher{regexp, BaseMatcher{pattern, insensitive}}
}

type FilepathMatcher struct {
  BaseMatcher
}

func (m FilepathMatcher) match(path string) (bool) {
  // Poor man's case insensitivity.
  if m.insensitive {
    path = strings.ToLower(path)
  }

  match, _ := filepath.Match(m.pattern, path)
  return match
}

func newFilepathMatcher(pattern string, insensitive bool) FilepathMatcher {
  return FilepathMatcher{BaseMatcher{pattern, insensitive}}
}

// An empty matcher always returns true. This is selected when pattern is empty.
type EmptyMatcher struct {}

func (m EmptyMatcher) match(path string) (bool) {
  return true;
}

func newEmptyMatcher() EmptyMatcher {
  return EmptyMatcher{}
}
