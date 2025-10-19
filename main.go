package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/sqweek/dialog"
)

func main() {
	fmt.Println("=== Git History to LaTeX Converter ===")
	fmt.Println()

	// Select git repository directory with GUI
	gitDir, err := dialog.Directory().Title("Select Git Repository Folder").Browse()
	if err != nil {
		fmt.Println("No git repository selected. Exiting.")
		return
	}
	fmt.Printf("Selected git repository: %s\n", gitDir)

	// Select output directory with GUI
	outputDir, err := dialog.Directory().Title("Select Output Folder for LaTeX File").Browse()
	if err != nil {
		fmt.Println("No output directory selected. Exiting.")
		return
	}
	fmt.Printf("Selected output directory: %s\n", outputDir)

	// Get git reference via CLI (simple text input)
	fmt.Print("\nEnter git reference to show history FROM (e.g., main, HEAD~10, or commit hash).\n")
	fmt.Print("Leave empty to show ALL commits: ")
	gitRef := readLine()

	// Get number of commits or full history
	var gitArgs []string
	if gitRef == "" {
		// Show all commits in the branch
		gitArgs = []string{"-C", gitDir, "log", "-p", "--unified=3"}
	} else {
		// Show commits from the specified reference
		gitArgs = []string{"-C", gitDir, "log", "-p", "--unified=3", gitRef + "..HEAD"}
	}

	// Get output filename
	fmt.Print("Enter output filename (default: git-history.tex): ")
	outputFile := readLine()
	if outputFile == "" {
		outputFile = "git-history.tex"
	}

	fmt.Println("\nProcessing git history...")

	// Get git log with patches (full history)
	cmd := exec.Command("git", gitArgs...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running git log: %v\n", err)
		fmt.Fprintf(os.Stderr, "Make sure the directory is a valid git repository\n")
		dialog.Message("Error: %v\nMake sure the directory is a valid git repository", err).Title("Git Error").Error()
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		dialog.Message("Error creating output directory: %v", err).Title("Error").Error()
		os.Exit(1)
	}

	// Create full output path
	outputPath := outputDir + string(os.PathSeparator) + outputFile

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		dialog.Message("Error creating output file: %v", err).Title("Error").Error()
		os.Exit(1)
	}
	defer f.Close()

	// Write LaTeX document
	if err := writeLatexDocument(f, string(output)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing LaTeX: %v\n", err)
		dialog.Message("Error writing LaTeX: %v", err).Title("Error").Error()
		os.Exit(1)
	}

	fmt.Printf("\n✓ Successfully generated %s\n", outputPath)
	dialog.Message("Successfully generated:\n%s", outputPath).Title("Success").Info()
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func writeLatexDocument(w io.Writer, gitDiff string) error {
	// Write LaTeX preamble
	fmt.Fprintln(w, `\documentclass[11pt,a4paper]{article}
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{listings}
\usepackage{xcolor}
\usepackage{fancyvrb}
\usepackage[margin=1in]{geometry}
\usepackage{hyperref}

% Define colors for diff output
\definecolor{diffadd}{RGB}{0,128,0}
\definecolor{diffrem}{RGB}{128,0,0}
\definecolor{diffinfo}{RGB}{0,0,128}
\definecolor{difffile}{RGB}{128,0,128}
\definecolor{commitcolor}{RGB}{139,0,139}

% Custom listing style for diffs
\lstdefinestyle{diffstyle}{
    basicstyle=\ttfamily\small,
    breaklines=true,
    columns=fullflexible,
    keepspaces=true,
    showspaces=false,
    showstringspaces=false,
    breakatwhitespace=false,
    tabsize=4,
}

\title{Git History}
\author{Generated from Git Repository}
\date{\today}

\begin{document}
\maketitle
\tableofcontents
\newpage`)

	// Parse and write diff sections
	scanner := bufio.NewScanner(strings.NewReader(gitDiff))
	var currentFile string
	var commitCount int
	inDiff := false
	inCommit := false

	for scanner.Scan() {
		line := scanner.Text()
		escapedLine := escapeLaTeX(line)

		if strings.HasPrefix(line, "commit ") {
			// Close previous diff section if any
			if inDiff {
				fmt.Fprintln(w, "\\end{Verbatim}")
				inDiff = false
			}

			if inCommit {
				fmt.Fprintln(w, "\\newpage")
			}

			commitCount++
			commitHash := strings.TrimPrefix(line, "commit ")
			if len(commitHash) > 10 {
				commitHash = commitHash[:10]
			}
			fmt.Fprintf(w, "\\section{Commit %d: %s}\n\n", commitCount, escapeLaTeX(commitHash))
			fmt.Fprintln(w, "\\begin{Verbatim}[commandchars=\\\\\\{\\},codes={\\catcode`$=3}]")
			fmt.Fprintf(w, "\\textcolor{commitcolor}{%s}\n", escapedLine)
			inCommit = true
			inDiff = true

		} else if strings.HasPrefix(line, "Author: ") || strings.HasPrefix(line, "Date: ") {
			fmt.Fprintf(w, "\\textcolor{commitcolor}{%s}\n", escapedLine)

		} else if strings.HasPrefix(line, "diff --git") {
			// Close previous verbatim if needed
			if inDiff {
				fmt.Fprintln(w, "\\end{Verbatim}")
			}

			// Extract file names
			parts := strings.Split(line, " ")
			if len(parts) >= 4 {
				currentFile = strings.TrimPrefix(parts[2], "a/")
			}

			// Start new subsection for file
			fmt.Fprintf(w, "\\subsection{%s}\n\n", escapeLaTeX(currentFile))
			fmt.Fprintln(w, "\\begin{Verbatim}[commandchars=\\\\\\{\\},codes={\\catcode`$=3}]")
			fmt.Fprintf(w, "\\textcolor{difffile}{%s}\n", escapedLine)
			inDiff = true

		} else if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
			if inDiff {
				fmt.Fprintf(w, "\\textcolor{difffile}{%s}\n", escapedLine)
			}

		} else if strings.HasPrefix(line, "@@") {
			if inDiff {
				fmt.Fprintf(w, "\\textcolor{diffinfo}{%s}\n", escapedLine)
			}

		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			if inDiff {
				fmt.Fprintf(w, "\\textcolor{diffadd}{%s}\n", escapedLine)
			}

		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			if inDiff {
				fmt.Fprintf(w, "\\textcolor{diffrem}{%s}\n", escapedLine)
			}

		} else {
			if inDiff {
				fmt.Fprintln(w, escapedLine)
			}
		}
	}

	// Close last diff section
	if inDiff {
		fmt.Fprintln(w, "\\end{Verbatim}")
	}

	// Close document
	fmt.Fprintln(w, "\\end{document}")

	return scanner.Err()
}

func escapeLaTeX(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\textbackslash{}`,
		`{`, `\{`,
		`}`, `\}`,
		`$`, `\$`,
		`&`, `\&`,
		`%`, `\%`,
		`#`, `\#`,
		`_`, `\_`,
		`~`, `\textasciitilde{}`,
		`^`, `\textasciicircum{}`,
	)
	return replacer.Replace(s)
}
