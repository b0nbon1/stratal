# Contributing to Stratal

Thank you for your interest in contributing to Stratal! We welcome contributions from the community and are grateful for your support.

## Table of Contents

- [Contributing to Stratal](#contributing-to-stratal)
  - [Table of Contents](#table-of-contents)
  - [Code of Conduct](#code-of-conduct)
  - [Getting Started](#getting-started)
  - [How to Contribute](#how-to-contribute)
  - [Development Setup](#development-setup)
    - [Prerequisites](#prerequisites)
    - [Local Development](#local-development)
    - [Using Docker](#using-docker)
  - [Coding Standards](#coding-standards)
    - [Go Code](#go-code)
    - [TypeScript/React Code](#typescriptreact-code)
    - [Database](#database)
    - [Commit Messages](#commit-messages)
  - [Testing](#testing)
    - [Running Tests](#running-tests)
    - [Writing Tests](#writing-tests)
  - [Submitting Changes](#submitting-changes)
  - [Reporting Bugs](#reporting-bugs)
  - [Feature Requests](#feature-requests)
  - [Development Tips](#development-tips)
    - [Useful Make Commands](#useful-make-commands)
    - [Project Structure](#project-structure)
    - [Debugging](#debugging)
  - [Getting Help](#getting-help)
  - [Recognition](#recognition)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors. We expect:

- Respectful and constructive communication
- Welcoming diverse perspectives and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/stratal.git
   cd stratal
   ```
3. **Add the upstream repository**:
   ```bash
   git remote add upstream https://github.com/original-owner/stratal.git
   ```
4. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## How to Contribute

We accept the following types of contributions:

- **Bug fixes**: Fix issues reported in the issue tracker
- **New features**: Implement features from the roadmap or propose new ones
- **Documentation**: Improve README, code comments, or add examples
- **Tests**: Add or improve test coverage
- **Code quality**: Refactoring, optimization, or cleanup
- **Task types**: Implement new built-in task types
- **UI improvements**: Enhance the web interface

## Development Setup

### Prerequisites

- **Go** 1.25 or higher
- **Node.js** 22+ and npm/yarn
- **PostgreSQL** 16+
- **Redis** 6+
- **Docker** (optional, for containerized development)

### Local Development

1. **Install dependencies**:
   ```bash
   make deps
   ```

2. **Set up the database**:
   ```bash
   # Create PostgreSQL database
   createdb stratal

   # Run migrations
   make migrate-up
   ```

3. **Configure environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your local settings
   ```

4. **Start Redis**:
   ```bash
   redis-server
   ```

5. **Run the application**:
   ```bash
   # Start the server
   make run-server

   # In another terminal, start a worker
   make run-worker

   # For web UI development
   cd web && npm install && npm run dev
   ```

### Using Docker

```bash
docker-compose up -d
```

## Coding Standards

### Go Code

- Follow standard Go conventions and use `gofmt` for formatting
- Run `go vet` to check for common mistakes
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and concise (single responsibility)
- Handle errors explicitly - don't ignore them

### TypeScript/React Code

- Follow the existing code style in the `web` directory
- Use TypeScript types instead of `any` where possible
- Use functional components with hooks
- Keep components small and focused
- Add proper PropTypes or TypeScript interfaces

### Database

- All schema changes must include migrations
- Use SQLC for type-safe database operations
- Write both `up` and `down` migrations
- Test migrations in both directions

### Commit Messages

Write clear, descriptive commit messages:

```
<type>: <short summary>

<optional detailed description>

<optional footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Build process or tooling changes

Example:
```
feat: add retry mechanism for failed tasks

Implement exponential backoff retry logic for tasks that fail
due to temporary errors. Configurable via job definition.

Closes #123
```

## Testing

### Running Tests

```bash
# Run all Go tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./pkg/worker/...

# Run frontend tests
cd web && npm test
```

### Writing Tests

- Add unit tests for new functionality
- Aim for at least 80% code coverage for new code
- Include both positive and negative test cases
- Use table-driven tests for multiple scenarios
- Mock external dependencies (database, Redis, HTTP calls)

Example test structure:
```go
func TestTaskExecutor_Execute(t *testing.T) {
    tests := []struct {
        name    string
        task    *Task
        want    *Result
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Submitting Changes

1. **Update your branch** with the latest upstream changes:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push your changes**:
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request**:
   - Go to the repository on GitHub
   - Click "New Pull Request"
   - Select your branch
   - Fill out the PR template with:
     - Description of changes
     - Related issue numbers
     - Testing performed
     - Screenshots (for UI changes)

4. **Code Review**:
   - Address review feedback promptly
   - Keep discussions focused and professional
   - Update your PR based on suggestions
   - Request re-review after making changes

5. **Merge Requirements**:
   - All tests must pass
   - Code review approval required
   - No merge conflicts
   - Follows coding standards

## Reporting Bugs

When reporting bugs, please include:

- **Clear title** describing the issue
- **Steps to reproduce** the problem
- **Expected behavior** vs actual behavior
- **Environment details**: OS, Go version, database version
- **Logs or error messages** (if applicable)
- **Screenshots** (for UI issues)

Use the GitHub issue tracker and apply the `bug` label.

## Feature Requests

For feature requests, please:

- Check if the feature already exists or is planned (see Future Work in README)
- Describe the problem your feature would solve
- Explain your proposed solution
- Consider alternative approaches
- Indicate if you're willing to implement it

Use the GitHub issue tracker and apply the `enhancement` label.

## Development Tips

### Useful Make Commands

```bash
make help          # Show all available commands
make build         # Build all binaries
make test          # Run tests
make lint          # Run linters
make migrate-up    # Run database migrations
make migrate-down  # Rollback migrations
make sqlc          # Regenerate SQLC code
```

### Project Structure

- `cmd/` - Application entry points (server, worker, CLI)
- `pkg/` - Reusable packages
- `internal/` - Private application code
- `migrations/` - Database migration files
- `web/` - React frontend application
- `test/` - Integration tests

### Debugging

- Use the built-in Go debugger (delve): `dlv debug ./cmd/server`
- Enable debug logging: Set `LOG_LEVEL=debug` in your environment
- Check logs in the database: Query the `job_run_logs` table
- Use the Insomnia collection for API testing

## Getting Help

- **GitHub Discussions**: Ask questions and discuss ideas
- **GitHub Issues**: Report bugs and request features
- **Code Comments**: Check inline documentation
- **README**: Review architecture and feature documentation

## Recognition

Contributors will be recognized in:
- The project's contributors list on GitHub
- Release notes for significant contributions
- The README (for major features)

Thank you for contributing to Stratal! Your efforts help make this project better for everyone.
