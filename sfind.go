package main

import (
  "os"
  "fmt"
  "path/filepath"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "SFind"
  app.Version = "0.0.1"
  app.Usage = "A simple find tool"
  app.UsageText = "sfind [options] [PATH] PATTERN"

  app.Flags = []cli.Flag {
    cli.BoolFlag {
      Name: "count, c",
      Usage: "Return a count of matches",
    },
    cli.BoolFlag {
      Name: "ext, e",
      Usage: "Use extended regex patterns",
    },
    cli.BoolFlag {
      Name: "invert, i",
      Usage: "Only return items that don't match PATTERN",
    },
    cli.BoolFlag {
      Name: "insensitive, I",
      Usage: "Case insensitive matches",
    },
    cli.BoolFlag {
      Name: "include-dirs, d",
      Usage: "Include directories in matches. This has presidence over 'dirs-only'",
    },
    cli.BoolFlag {
      Name: "dirs-only, D",
      Usage: "Only match directories",
    },
  }

  app.Action = func(c *cli.Context) {
    base_path, pattern := determineArgs(c)
    matcher := createMatcher(pattern, c)
    result := createResult(c)

    outputResults(base_path, result, matcher, c)
  }

  app.Run(os.Args)
}

func outputResults(base_path string, result Result, matcher FileMatcher, c *cli.Context) {
  include_dirs := c.Bool("include-dirs")
  dirs_only := c.Bool("dirs-only")
  // Small shortcut to skip dir checks early.
  no_dir_filters := !(include_dirs || dirs_only)

  // Allow the result to print something before files are walked.
  result.beforeResults()

  err := filepath.Walk(base_path, func (path string, fileInfo os.FileInfo, _ error) error {
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

    matched, _ := matcher.match(filepath.Base(path))

    if matched {
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

func determineArgs(c *cli.Context) (string, string) {
  var base_path, pattern string
  var err error
  args := c.Args()
  length := c.NArg();

  switch {
    case length == 0:
      fmt.Println("No arguments provided")
      os.Exit(1)
    case length == 1:
      // Assume the single argument is a pattern. Default root to cwd.
      base_path, err = os.Getwd()

      if err != nil {
          fmt.Println(err)
          os.Exit(1)
      }

      pattern = args.First()
    case length > 1:
      // Assume path will be the first arg, and pattern the second.
      base_path = args.First()
      pattern = args.Get(1)
  }

  return base_path, pattern
}

func createMatcher(pattern string, c *cli.Context) FileMatcher {
  invert := c.Bool("invert")
  insensitive := c.Bool("insensitive")

  if c.Bool("ext") {
    return newRegexMatcher(pattern, invert, insensitive)
  }

  return newFilepathMatcher(pattern, invert, insensitive)
}

func createResult(c *cli.Context) Result {
  if c.Bool("count") {
    return &CountResult{}
  }

  return &PrintResult{}
}
