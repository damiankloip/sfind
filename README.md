# SFind - A simple find tool

This is a test project mainly, providing a simple find command to recursively
search for files. Like find but without all the crazy functionality. It can
search with simple file glob patterns, or extended regex. By default sfind will
only match files but can be used to match directories, or both.

### Usage:
   sfind [options] [PATH] PATTERN

### Options:
```
   --count, -c		Return a count of matches
   --ext, -e		Use extended regex patterns
   --invert, -i		Only return items that don't match PATTERN
   --insensitive, -I	Case insensitive matches
   --include-dirs, -d	Include directories in matches. This has presidence over 'dirs-only'
   --dirs-only, -D	Only match directories
   --help, -h		show help
   --version, -v	print the version
```

### Examples:

Find all YAML files from the current directory

```
sfind '*.yml'
```

Find all YAML files from the ''/test' directory

```
sfind /test '*.yml'
```

Get a count of all Files from the current directory

```
sfind -c '*'
```
