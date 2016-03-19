package main

import (
  "os"
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
    cli.BoolFlag {
      Name: "full-path, f",
      Usage: "Match PATTERN again the full file (or directory) path. Ext option is implied.",
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
