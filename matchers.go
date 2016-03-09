package main

import (
  "regexp"
  "path/filepath"
  "strings"
)

type FileMatcher interface {
  match(path string) (bool, error)
}

type BaseMatcher struct {
  pattern string
  insensitive bool
}

type RegexMatcher struct {
  BaseMatcher
}

func (m RegexMatcher) match(path string) (bool, error) {
  match, err := regexp.MatchString(m.pattern, path)
  return match, err
}

func newRegexMatcher(pattern string, insensitive bool) RegexMatcher {
  // If the pattern is empty, assume everything.
  if pattern == "" {
    pattern = "\\.*"
  }

  // If this is case-insensitive, prepend '(?i)' to the regex. Likely NOT the
  // best way to do this!
  if insensitive {
    // This assignment is not ideal, but easy.
    prefix := "(?i)"
    pattern = prefix + pattern
  }

  return RegexMatcher{BaseMatcher{pattern, insensitive}}
}

type FilepathMatcher struct {
  BaseMatcher
}

func (m FilepathMatcher) match(path string) (bool, error) {
  // Poor man's case insensitivity.
  if m.insensitive {
    path = strings.ToLower(path)
  }

  match, err := filepath.Match(m.pattern, path)
  return match, err
}

func newFilepathMatcher(pattern string, insensitive bool) FilepathMatcher {
  // If the pattern is empty, assume everything.
  if pattern == "" {
    pattern = "*"
  }

  return FilepathMatcher{BaseMatcher{pattern, insensitive}}
}
