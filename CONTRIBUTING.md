# Contributing to Parse DMARC

Thank you for your interest in contributing to Parse DMARC! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- Bun
- Just
- Git

### Getting Started

1. Fork and clone the repository:

```bash
git clone https://github.com/meysam81/parse-dmarc.git
cd parse-dmarc
```

2. Install dependencies:

```bash
just install-deps
```

3. Build the project:

```bash
just build
```

## Project Structure

```
parse-dmarc/
├── main.go                # Main application entry point
├── internal/
│   ├── api/               # REST API and web server
│   ├── config/            # Configuration management
│   ├── imap/              # IMAP client for fetching emails
│   ├── parser/            # DMARC XML parser
│   └── storage/           # SQLite database layer
├── src/                   # Vue.js 3 dashboard
│   ├── assets/
│   ├── components/
│   ├── lib/
│   ├── stores/
│   ├── App.vue
│   └── main.js
├── package.json
├── Justfile
├── Dockerfile
└── README.md
```

## Development Workflow

Run the application in development mode:

```bash
just dev
```

Run tests:

```bash
just test
```

### Making Changes

1. Create a new branch:

```bash
git checkout -b feature/your-feature-name
```

2. Make your changes and commit:

```bash
git add .
git commit -m "Description of changes"
```

3. Push and create a pull request:

```bash
git push origin feature/your-feature-name
```

## Code Style

### Go Code

- Follow standard Go formatting (`gofmt`)
- Add comments for exported functions and types
- Keep functions small and focused
- Use meaningful variable names

### Vue.js Code

- Use Vue 3 Composition API
- Follow Vue style guide
- Keep components modular and reusable

## Testing

### Go Tests

Add tests for all new functionality:

```bash
just test
```

### Manual Testing

1. Generate a config file:

```bash
./bin/parse-dmarc -gen-config
```

2. Edit config.json with test credentials

3. Run in serve-only mode for UI testing:

```bash
./bin/parse-dmarc -serve-only
```

## Pull Request Guidelines

- Ensure all tests pass
- Update documentation if needed
- Keep PRs focused on a single feature/fix
- Write clear commit messages
- Reference any related issues

## Areas for Contribution

- **Forensic Reports**: Add support for DMARC forensic reports (RUF)
- **OAuth2**: Implement OAuth2 for IMAP authentication
- **Export**: Add CSV/JSON export functionality
- **Alerts**: Email alerts for compliance issues
- **Analytics**: Historical trend analysis
- **Documentation**: Improve docs and examples
- **Tests**: Increase test coverage

## Questions?

Open an issue for questions or discussions.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
