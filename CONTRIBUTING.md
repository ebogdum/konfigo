# Contributing to Konfigo

First off, thank you for considering contributing to Konfigo! It's people like you that make Konfigo such a great tool.

This document provides guidelines for contributing to Konfigo. Please read it carefully to ensure that your contributions can be accepted.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
  - [Pull Requests](#pull-requests)
- [Development Setup](#development-setup)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running Tests](#running-tests)
- [Styleguides](#styleguides)
  - [Go Code](#go-code)
  - [Git Commit Messages](#git-commit-messages)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by the [Konfigo Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [PROJECT_EMAIL_OR_CONTACT_LINK].

*(Note: You will need to create a `CODE_OF_CONDUCT.md` file. You can use a template like the [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct.html)).*

## How Can I Contribute?

### Reporting Bugs

Bugs are tracked as [GitHub issues](https://github.com/ebogdum/konfigo/issues). Before submitting a bug report, please perform the following steps:

1.  **Check the [project documentation](https://ebogdum.github.io/konfigo/)** to see if the behavior is intended.
2.  **Search [existing issues](https://github.com/ebogdum/konfigo/issues)** to see if the bug has already been reported. If it has, add a comment to the existing issue instead of opening a new one.
3.  If you cannot find an existing issue that describes your problem, **create a new issue**.

When filing a bug report, please include as much detail as possible. Fill out the required template, the information it asks for helps us resolve issues faster. Key details include:

*   A clear and descriptive title.
*   Steps to reproduce the bug.
*   What you expected to happen.
*   What actually happened.
*   Your environment (Go version, OS, Konfigo version).
*   Any relevant logs or screenshots.

### Suggesting Enhancements

Enhancement suggestions are also tracked as [GitHub issues](https://github.com/ebogdum/konfigo/issues).

1.  **Use a clear and descriptive title** for the issue to identify the suggestion.
2.  **Provide a step-by-step description of the suggested enhancement** in as many details as possible.
3.  **Explain why this enhancement would be useful** to most Konfigo users.
4.  **Provide code examples** if possible, to illustrate the use cases.

### Your First Code Contribution

Unsure where to begin contributing to Konfigo? You can start by looking through these `good first issue` and `help wanted` issues:

*   [Good first issues](https://github.com/ebogdum/konfigo/labels/good%20first%20issue) - issues which should only require a few lines of code, and a test or two.
*   [Help wanted issues](https://github.com/ebogdum/konfigo/labels/help%20wanted) - issues which should be a bit more involved than `good first issue` issues.

*(Note: You'll need to create these labels in your GitHub issues if they don't exist.)*

### Pull Requests

When you're ready to contribute code, please follow these steps:

1.  **Fork the repository** and create your branch from `main` (or the default branch).
2.  **Set up your development environment** (see [Development Setup](#development-setup)).
3.  **Make your changes.**
    *   Ensure your code adheres to the [Styleguides](#styleguides).
    *   Add new tests for any new features or bug fixes.
    *   Ensure all tests pass.
    *   Update documentation if you are changing behavior or adding features.
4.  **Commit your changes** using a descriptive commit message (see [Git Commit Messages](#git-commit-messages)).
5.  **Push your branch** to your fork.
6.  **Open a pull request** to the `main` branch of the `ebogdum/konfigo` repository.
    *   Fill in the pull request template with a clear description of your changes.
    *   Link to any relevant issues.
7.  **Address any review comments** and make necessary changes.

## Development Setup

This section describes how to set up your development environment to contribute to Konfigo.

### Prerequisites

*   Go (version X.Y.Z or higher - *please specify the version you use/require*)
*   Git
*   *(Any other specific tools or libraries needed, e.g., a particular linter, Docker, etc.)*

### Installation

1.  Fork the `ebogdum/konfigo` repository on GitHub.
2.  Clone your fork locally:
    ```sh
    git clone https://github.com/YOUR_USERNAME/konfigo.git
    cd konfigo
    ```
3.  Add the upstream repository:
    ```sh
    git remote add upstream https://github.com/ebogdum/konfigo.git
    ```
4.  Install dependencies:
    ```sh
    go mod tidy
    # or
    go get -d ./...
    ```

## Styleguides

### Go Code

*   All Go code should be formatted with `gofmt`. Most IDEs do this automatically.
*   Run `go vet ./...` to catch common errors.
*   Consider using `golangci-lint` or a similar linter for more comprehensive checks. *(Specify if you use a particular linter and its configuration)*.
*   Follow standard Go best practices (e.g., effective Go, error handling).
*   Write clear and concise comments where necessary.
*   *(Add any project-specific coding conventions here.)*

### Git Commit Messages

*   Use the present tense ("Add feature" not "Added feature").
*   Use the imperative mood ("Move cursor to..." not "Moves cursor to...").
*   Limit the first line to 72 characters or less.
*   Reference issues and pull requests liberally after the first line.
*   Consider using [Conventional Commits](https://www.conventionalcommits.org/) if you prefer a more structured approach.

Example:
```
feat: Add support for YAML configuration files

This commit introduces the ability to load configurations from YAML files,
in addition to the existing JSON and TOML formats.

Resolves #123
```

## Community

*   **Project Documentation:** [https://ebogdum.github.io/konfigo/](https://ebogdum.github.io/konfigo/)
*   **Issue Tracker:** [https://github.com/ebogdum/konfigo/issues](https://github.com/ebogdum/konfigo/issues)

---

We look forward to your contributions!