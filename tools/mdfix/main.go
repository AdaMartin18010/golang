package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Config controls the fixer behavior.
type Config struct {
	RootDir          string
	DryRun           bool
	Verbose          bool
	Backup           bool
	IncludeHidden    bool
	MaxFileSizeBytes int64
}

var (
	headingRegexp       = regexp.MustCompile(`^(#{1,6})([^#\s].*)$`)
	trailingHashesRegex = regexp.MustCompile(`^(#{1,6} .+?)\s+#+\s*$`)
	listDashRegex       = regexp.MustCompile(`^\-([^\s\-].*)$`)
	listStarRegex       = regexp.MustCompile(`^\*([^\s\*].*)$`)
	listPlusRegex       = regexp.MustCompile(`^\+([^\s\+].*)$`)
	orderedListRegex    = regexp.MustCompile(`^(\d+)\.([^\s].*)$`)
	blockQuoteRegex     = regexp.MustCompile(`^>([^\s>].*)$`)
	trimTrailingSpaces  = regexp.MustCompile(`\s+$`)
)

func main() {
	cfg := Config{}
	flag.StringVar(&cfg.RootDir, "root", ".", "root directory to scan")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "only report files that would change")
	flag.BoolVar(&cfg.Verbose, "v", false, "verbose output")
	flag.BoolVar(&cfg.Backup, "backup", true, "write .bak alongside before modifying")
	flag.BoolVar(&cfg.IncludeHidden, "all", false, "include hidden files and directories")
	var maxSizeMB int
	flag.IntVar(&maxSizeMB, "max-mb", 5, "skip files larger than this size (MB)")
	flag.Parse()
	cfg.MaxFileSizeBytes = int64(maxSizeMB) * 1024 * 1024

	var changedFiles []string
	var totalMD, skippedLarge int

	walkErr := filepath.WalkDir(cfg.RootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if !cfg.IncludeHidden && strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(name), ".md") {
			return nil
		}
		totalMD++
		info, err := d.Info()
		if err == nil && info.Size() > cfg.MaxFileSizeBytes {
			skippedLarge++
			if cfg.Verbose {
				fmt.Printf("skip large file: %s (%d bytes)\n", path, info.Size())
			}
			return nil
		}
		orig, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fixed := fixMarkdown(orig)
		if !bytes.Equal(orig, fixed) {
			changedFiles = append(changedFiles, path)
			if cfg.DryRun {
				return nil
			}
			if cfg.Backup {
				_ = os.WriteFile(path+".bak", orig, 0644)
			}
			if err := os.WriteFile(path, fixed, 0644); err != nil {
				return err
			}
			if cfg.Verbose {
				fmt.Println("fixed:", path)
			}
		}
		return nil
	})
	if walkErr != nil {
		fmt.Fprintln(os.Stderr, walkErr)
		os.Exit(1)
	}

	// Summary
	fmt.Printf("scanned md files: %d, skipped large: %d, %s files: %d\n",
		totalMD, skippedLarge, ternary(cfg.DryRun, "would change", "changed"), len(changedFiles))
	if cfg.DryRun {
		for _, f := range changedFiles {
			fmt.Println("would change:", f)
		}
	}
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func fixMarkdown(input []byte) []byte {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)
	var out []string
	inCodeFence := false
	var fenceMarker string

	for scanner.Scan() {
		line := scanner.Text()

		// Detect fenced code blocks: ``` or ~~~
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			marker := line[:3]
			if inCodeFence {
				if marker == fenceMarker {
					inCodeFence = false
				}
				out = append(out, line)
				continue
			}
			inCodeFence = true
			fenceMarker = marker
			out = append(out, line)
			continue
		}

		if inCodeFence {
			out = append(out, line)
			continue
		}

		origLine := line

		// Normalize headings: "#Title" => "# Title" and remove trailing ###
		if m := headingRegexp.FindStringSubmatch(line); m != nil {
			line = m[1] + " " + strings.TrimSpace(m[2])
			line = trailingHashesRegex.ReplaceAllString(line, "$1")
			out = append(out, line)
			continue
		}
		line = trailingHashesRegex.ReplaceAllString(line, "$1")

		// Normalize list markers: ensure a space after -, *, +, and ordered lists
		if m := listDashRegex.FindStringSubmatch(line); m != nil {
			line = "- " + strings.TrimSpace(m[1])
		} else if m := listStarRegex.FindStringSubmatch(line); m != nil {
			line = "* " + strings.TrimSpace(m[1])
		} else if m := listPlusRegex.FindStringSubmatch(line); m != nil {
			line = "+ " + strings.TrimSpace(m[1])
		} else if m := orderedListRegex.FindStringSubmatch(line); m != nil {
			line = m[1] + ". " + strings.TrimSpace(m[2])
		}

		// Normalize blockquotes: ensure "> "
		if m := blockQuoteRegex.FindStringSubmatch(line); m != nil {
			line = "> " + strings.TrimSpace(m[1])
		}

		// Trim trailing spaces, but preserve two-space hard break
		if strings.HasSuffix(line, "   ") {
			// more than two spaces: reduce to exactly two
			line = strings.TrimRight(line, " ") + "  "
		} else if strings.HasSuffix(line, "  ") {
			// keep as-is
		} else {
			line = trimTrailingSpaces.ReplaceAllString(line, "")
		}

		out = append(out, line)

		_ = origLine // reserved for future heuristics
	}

	// Ensure a blank line after headings
	out = ensureBlankAfterHeading(out)

	// Ensure file ends with a single newline
	res := strings.Join(out, "\n")
	if !strings.HasSuffix(res, "\n") {
		res += "\n"
	}
	return []byte(res)
}

func ensureBlankAfterHeading(lines []string) []string {
	var result []string
	for i := 0; i < len(lines); i++ {
		result = append(result, lines[i])
		if isHeadingLine(lines[i]) {
			// if next line exists and is not blank and not fence start, insert blank
			if i+1 < len(lines) && len(strings.TrimSpace(lines[i+1])) != 0 && !isFenceStart(lines[i+1]) {
				result = append(result, "")
			}
		}
	}
	return result
}

func isHeadingLine(s string) bool {
	for level := 1; level <= 6; level++ {
		prefix := strings.Repeat("#", level) + " "
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func isFenceStart(s string) bool {
	return strings.HasPrefix(s, "```") || strings.HasPrefix(s, "~~~")
}
