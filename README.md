# Git History to LaTeX Converter

## Directory Structure

```
git2latex/
├── main.go              # Main application source code
├── Makefile             # Build automation (cross-platform)
├── README.md            # This file
├── go.mod               # Go module dependencies (generated)
├── go.sum               # Go module checksums (generated)
```

## Installation

### Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/git2latex.git
cd git2latex

# Install dependencies and build
make

# Run the application
make run
```

### Manual Installation

```bash
# Initialize Go modules
go mod init git2latex

# Install dependencies
go get github.com/sqweek/dialog

# Build the application
go build -o git2latex main.go

# Run it
./git2latex          # Unix/macOS
git2latex.exe        # Windows
```

## Usage

### Using Make (Recommended)

```bash
# Build the project
make

# Run with GUI folder selection
make run

# After generating .tex file, compile to PDF
make pdf

# Clean all generated files
make clean

# Show detected operating system
make info
```

### Using the Application

1. **Run the program:**
   ```bash
   make run
   ```

2. **Select Git Repository:**
   - A file dialog will open
   - Navigate to and select your git repository folder

3. **Select Output Directory:**
   - Another dialog will open
   - Choose where to save the LaTeX file

4. **Configure Options:**
   - Enter git reference (e.g., `main`, `HEAD~10`, or press Enter for all commits)
   - Enter output filename (or press Enter for default: `git-history.tex`)

5. **Compile to PDF:**
   ```bash
   make pdf
   ```
   Or manually:
   ```bash
   pdflatex -interaction=nonstopmode git-history.tex
   pdflatex -interaction=nonstopmode git-history.tex  # Run twice for TOC
   ```