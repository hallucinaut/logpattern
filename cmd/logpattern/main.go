package main

import (
	"os/signal"
	"syscall"
	"context"
	"os/signal"
	"syscall"
	"context"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

type LogPattern struct {
	Pattern     string
	Count       int
	Percentage  float64
	Samples     []string
	FirstSeen   time.Time
	LastSeen    time.Time
}

type LogLine struct {
	Timestamp time.Time
	Message   string
	Level     string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(color.CyanString("logpattern - Log Pattern Detector"))
		fmt.Println()
		fmt.Println("Usage: logpattern <logfile>")
		os.Exit(1)
	}

	filename := os.Args[1]
	logs, err := parseLogFile(filename)
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	patterns := detectPatterns(logs)
	displayPatterns(patterns)
}

func parseLogFile(filename string) ([]LogLine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logs []LogLine
	scanner := bufio.NewScanner(file)
	levelRegex := regexp.MustCompile(`\b(DEBUG|INFO|WARN|ERROR|FATAL)\b`)

	for scanner.Scan() {
		line := scanner.Text()
		level := "UNKNOWN"
		if levelRegex.MatchString(line) {
			level = levelRegex.FindString(line)
		}

		logs = append(logs, LogLine{
			Timestamp: time.Now(),
			Message:   line,
			Level:     level,
		})
	}

	return logs, scanner.Err()
}

func detectPatterns(logs []LogLine) []LogPattern {
	patternMap := make(map[string]*LogPattern)

	for _, log := range logs {
		// Normalize the log message
		normalized := normalizeMessage(log.Message)

		pattern, exists := patternMap[normalized]
		if !exists {
			pattern = &LogPattern{
				Pattern: normalized,
				Samples: make([]string, 0, 3),
			}
			patternMap[normalized] = pattern
		}

		pattern.Count++
		if len(pattern.Samples) < 3 {
			pattern.Samples = append(pattern.Samples, log.Message)
		}
	}

	var patterns []LogPattern
	for _, p := range patternMap {
		patterns = append(patterns, *p)
	}

	// Calculate percentages
	total := len(logs)
	for i := range patterns {
		patterns[i].Percentage = float64(patterns[i].Count) / float64(total) * 100
	}

	// Sort by count
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Count > patterns[j].Count
	})

	return patterns[:min(20, len(patterns))]
}

func normalizeMessage(msg string) string {
	// Replace numbers with placeholders
	re := regexp.MustCompile(`\b\d+\b`)
	msg = re.ReplaceAllString(msg, "<NUM>")

	// Replace UUIDs
	re = regexp.MustCompile(`\b[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}\b`)
	msg = re.ReplaceAllString(msg, "<UUID>")

	// Replace IP addresses
	re = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	msg = re.ReplaceAllString(msg, "<IP>")

	// Replace timestamps
	re = regexp.MustCompile(`\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}`)
	msg = re.ReplaceAllString(msg, "<TIMESTAMP>")

	return strings.TrimSpace(msg)
}

func displayPatterns(patterns []LogPattern) {
	fmt.Println(color.CyanString("\n=== LOG PATTERN ANALYSIS ===\n"))

	for _, p := range patterns {
		countColor := color.GreenString
		if p.Count < 10 {
			countColor = color.YellowString
		} else if p.Count < 100 {
			countColor = color.HiYellowString
		}

		fmt.Printf("%-3d | %-10s | %s\n",
			p.Count,
			countColor(fmt.Sprintf("%.2f%%", p.Percentage)),
			truncatePattern(p.Pattern, 60),
		)

		if p.Count > 5 {
			fmt.Println(color.HiWhiteString("  Samples:"))
			for _, sample := range p.Samples {
				fmt.Printf("    %s\n", truncatePattern(sample, 80))
			}
		}
		fmt.Println()
	}
}

func truncatePattern(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}