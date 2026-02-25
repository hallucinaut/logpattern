# logpattern

Log Pattern Detector - analyzes log files to identify recurring patterns.

## Purpose

Parse log files and extract common patterns by normalizing variable content such as numbers, UUIDs, IP addresses, and timestamps.

## Installation

```bash
go build -o logpattern ./cmd/logpattern
```

## Usage

```bash
logpattern <logfile>
```

### Examples

```bash
# Analyze application log
logpattern /var/log/app.log

# Analyze error logs
logpattern /var/log/error.log

# Analyze from stdin
cat app.log | logpattern -
```

## Output

```
=== LOG PATTERN ANALYSIS ===

   15 |  37.50% | Error: Connection to <IP> failed after <NUM> retries
  Samples:
    Error: Connection to 192.168.1.100 failed after 3 retries
    Error: Connection to 192.168.1.101 failed after 2 retries

    8 |  20.00% | User <UUID> logged in successfully
  Samples:
    User a3f5b8c2-d9e1-f4a6-b7c8-d9e0f1a2b3c4 logged in successfully

   5 |  12.50% | Request completed in <NUM>ms
```

## Normalization

The tool normalizes the following patterns:
- Numbers replaced with `<NUM>`
- UUIDs replaced with `<UUID>`
- IP addresses replaced with `<IP>`
- Timestamps replaced with `<TIMESTAMP>`

## Dependencies

- Go 1.21+
- github.com/fatih/color

## Build and Run

```bash
# Build
go build -o logpattern ./cmd/logpattern

# Run
go run ./cmd/logpattern /path/to/logfile.log
```

## License

MIT