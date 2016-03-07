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
  invert bool
  insensitive bool
}

type RegexMatcher struct {
  BaseMatcher
}

func (m RegexMatcher) match(path string) (bool, error) {
  match, err := regexp.MatchString(m.pattern, path)
  return (match != m.invert), err
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
  return (match != m.invert), err
}
