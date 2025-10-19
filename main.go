package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	outputFile := flag.String("o", "git-diff.tex", "Output LaTeX file")
	gitRef := flag.String("ref", "HEAD", "Git reference to diff against (e.g., HEAD~1, main, commit hash)")
	flag.Parse()

	// Get git diff with full tree
	cmd := exec.Command("git", "diff", *gitRef, "--patch", "--unified=3")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running git diff: %v\n", err)
		os.Exit(1)
	}

	// Create output file
	f, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Write LaTeX document
	if err := writeLatexDocument(f, string(output)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing LaTeX: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s\n", *outputFile)
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

\title{Git Diff Output}
\author{Generated from Git Repository}
\date{\today}

\begin{document}
\maketitle
\tableofcontents
\newpage`)

	// Parse and write diff sections
	scanner := bufio.NewScanner(strings.NewReader(gitDiff))
	var currentFile string
	inDiff := false

	for scanner.Scan() {
		line := scanner.Text()
		escapedLine := escapeLaTeX(line)

		if strings.HasPrefix(line, "diff --git") {
			// Close previous diff section if any
			if inDiff {
				fmt.Fprintln(w, "\\end{Verbatim}")
				fmt.Fprintln(w)
			}

			// Extract file names
			parts := strings.Split(line, " ")
			if len(parts) >= 4 {
				currentFile = strings.TrimPrefix(parts[2], "a/")
			}

			// Start new section
			fmt.Fprintf(w, "\\section{%s}\n\n", escapeLaTeX(currentFile))
			fmt.Fprintln(w, "\\begin{Verbatim}[commandchars=\\\\\\{\\},codes={\\catcode`$=3}]")
			fmt.Fprintf(w, "\\textcolor{difffile}{%s}\n", escapedLine)
			inDiff = true

		} else if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
			fmt.Fprintf(w, "\\textcolor{difffile}{%s}\n", escapedLine)

		} else if strings.HasPrefix(line, "@@") {
			fmt.Fprintf(w, "\\textcolor{diffinfo}{%s}\n", escapedLine)

		} else if strings.HasPrefix(line, "+") {
			fmt.Fprintf(w, "\\textcolor{diffadd}{%s}\n", escapedLine)

		} else if strings.HasPrefix(line, "-") {
			fmt.Fprintf(w, "\\textcolor{diffrem}{%s}\n", escapedLine)

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
