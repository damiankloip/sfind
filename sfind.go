package main

import (
  "os"
  "fmt"
  "path/filepath"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "Simple Find"
  app.Version = "0.0.1"

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
  }

  app.Action = func(c *cli.Context) {
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

    var matcher FileMatcher
    var result Result
    invert := c.Bool("invert")
    insensitive := c.Bool("insensitive")

    if c.Bool("ext") {
      matcher = newRegexMatcher(pattern, invert, insensitive)
    } else {
      matcher = FilepathMatcher{BaseMatcher{pattern, invert, insensitive}}
    }

    if c.Bool("count") {
      result = &CountResult{}
    } else {
      result = &PrintResult{}
    }

    outputResults(base_path, result, matcher)
  }

  app.Run(os.Args)
}

func outputResults(base_path string, result Result, matcher FileMatcher) {
  // Allow the result to print something before files are walked.
  result.beforeResults()

  err := filepath.Walk(base_path, func (path string, fileInfo os.FileInfo, _ error) error {
    // Skip dirs. Could make this configurable.
    if fileInfo.IsDir() {
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
