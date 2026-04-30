# Contributing to Parse DMARC

Thank you for your interest in contributing to Parse DMARC! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- Git

### Getting Started

1. Fork and clone the repository:

```bash
git clone https://github.com/meysam81/parse-dmarc.git
cd parse-dmarc
```

2. Install dependencies:

```bash
make install-deps
```

3. Build the project:

```bash
make build
```

## Project Structure

```
parse-dmarc/
├── cmd/parse-dmarc/       # Main application entry point
├── internal/
│   ├── api/               # REST API and web server
│   ├── config/            # Configuration management
│   ├── imap/              # IMAP client for fetching emails
│   ├── parser/            # DMARC XML parser
│   └── storage/           # SQLite database layer
├── frontend/              # Vue.js 3 dashboard
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   └── App.vue
│   └── package.json
├── Makefile
├── Dockerfile
└── README.md
```

## Development Workflow

### Backend Development

Run the application in development mode:

```bash
make dev
```

Run tests:

```bash
go test ./...
```

Add tests for new features in `*_test.go` files.

### Frontend Development

Start the frontend dev server with hot reload:

```bash
cd frontend
npm run dev
```

Build the frontend:

```bash
cd frontend
npm run build
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
go test -v ./...
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
- **OAuth2 (Microsoft 365)**: Add Microsoft 365 / Entra ID provider to the
  IMAP XOAUTH2 flow (Google is supported via `--oauth-login --provider google`)
- **Export**: Add CSV/JSON export functionality
- **Alerts**: Email alerts for compliance issues
- **Analytics**: Historical trend analysis
- **Documentation**: Improve docs and examples
- **Tests**: Increase test coverage

## Questions?

Open an issue for questions or discussions.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
