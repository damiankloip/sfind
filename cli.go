package main

import (
  "os"
  "fmt"
  "path/filepath"
  "sync"
  "github.com/codegangsta/cli"
)

type FileData struct {
  Path string
  Info os.FileInfo
}

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

  var numWorkers int = 8

  path_channel, err := func () (chan FileData, error) {
    channel := make(chan FileData)
    var err error

    go func() {
      err = filepath.Walk(base_path, func (path string, fileInfo os.FileInfo, file_err error) error {
        if file_err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        channel <- FileData{path, fileInfo}
        return file_err
      })

      close(channel)

      return
    }()

    return channel, err
  }()

  if err != nil {
    fmt.Println(err)
  }

  var wg sync.WaitGroup

  for workers := 0; workers < numWorkers; workers++ {
    wg.Add(1)

    go func() error {
      defer wg.Done()
      for fileData := range path_channel {
        path := fileData.Path
        is_dir := fileData.Info.IsDir()

        // Determine how to deal with dirs based on any flags.
        if (no_dir_filters && is_dir) {
          // If there are no dir filters and this is dir, return. Otherwise, this
          // will try to use the matcher below. By default, dirs are not included in
          // results.
          continue
        // This means there are dir filters, so check against them. Don't check
        // include_dirs as if it hasn't been caught above (no_dir_filters) we want
        // everything to go to below matching.
        } else if (!is_dir && dirs_only) {
          continue
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
      }

      return nil
    }()
  }

  wg.Wait()
  // Allow the result to print something after the files have been walked.
  result.afterResults()
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

  if c.Bool("ext") || c.Bool("full-path") {
    return newRegexMatcher(pattern, insensitive)
  }

  return newFilepathMatcher(pattern, insensitive)
}

// Choose an appropriate result.
// This is based on whether or not a count is needed.
func createResult(c *cli.Context) Result {
  if c.Bool("count") {
    return &CountResult{mutex: &sync.Mutex{}}
  }

  return &PrintResult{}
}
