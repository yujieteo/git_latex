.PHONY: all clean build pdf

# Default target
all: pdf

# Build the Go program
build:
	go build -o git2latex main.go

# Generate LaTeX from git diff
git-diff.tex: build
	./git2latex -o git-diff.tex

# Compile LaTeX to PDF
pdf: git-diff.tex
	pdflatex git-diff.tex
	pdflatex git-diff.tex

# Clean generated files
clean:
	rm -f git2latex git-diff.tex git-diff.pdf git-diff.aux git-diff.log git-diff.out git-diff.toc

# Custom reference (e.g., make diff REF=HEAD~5)
diff:
	./git2latex -ref $(REF) -o git-diff.tex
	pdflatex git-diff.tex
	pdflatex git-diff.tex