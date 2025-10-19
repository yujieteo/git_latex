# Detect operating system
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    EXE_EXT := .exe
    RM := del /Q
    RMDIR := rmdir /S /Q
    NULL := 2>NUL
else
    DETECTED_OS := $(shell uname -s)
    EXE_EXT :=
    RM := rm -f
    RMDIR := rm -rf
    NULL := 2>/dev/null
endif

BINARY_NAME := git2latex$(EXE_EXT)

.PHONY: all clean build pdf deps run

# Default target
all: deps build

# Install Go dependencies
deps:
	-go mod init git2latex $(NULL)
	go get github.com/sqweek/dialog
	go mod tidy

# Build the Go program
build:
	go build -o $(BINARY_NAME) main.go

# Run the program (will open GUI dialogs)
run: build
ifeq ($(DETECTED_OS),Windows)
	$(BINARY_NAME)
else
	./$(BINARY_NAME)
endif

# Compile LaTeX to PDF (run twice for TOC)
pdf: git-history.tex
	pdflatex -interaction=nonstopmode git-history.tex
	pdflatex -interaction=nonstopmode git-history.tex

# Generate LaTeX from git history (requires running git2latex first)
git-history.tex:
	@echo "Please run 'make run' to generate the .tex file with GUI folder selection"
	@exit 1

# Clean generated files
clean:
ifeq ($(DETECTED_OS),Windows)
	-$(RM) $(BINARY_NAME) git-history.tex git-history.pdf git-history.aux git-history.log git-history.out git-history.toc go.mod go.sum $(NULL)
else
	$(RM) $(BINARY_NAME) git-history.tex git-history.pdf git-history.aux git-history.log git-history.out git-history.toc go.mod go.sum
endif

# Show detected OS
info:
	@echo "Detected OS: $(DETECTED_OS)"
	@echo "Binary name: $(BINARY_NAME)"