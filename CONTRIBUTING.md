# Contributing to Camp

Thanks for your interest in contributing to Camp.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/camp.git`
3. Install Go 1.25 or higher
4. Run `go mod download` to install dependencies
5. Run `make run` to start the server

## Making Changes

1. Create a new branch: `git checkout -b your-branch-name`
2. Make your changes
3. Test your changes: `go test ./...`
4. Commit your changes with a clear message
5. Push to your fork: `git push origin your-branch-name`
6. Open a pull request

## Code Guidelines

- Write clear, readable code
- Add tests for new features
- Keep functions small and focused
- Handle errors properly
- Follow Go's standard formatting (`go fmt`)

## Testing

Run tests before submitting:
```bash
go test ./...
```

Add tests for any new code in `*_test.go` files.

## Pull Requests

- Describe what your changes do
- Reference any related issues
- Make sure tests pass
- Keep PRs focused on a single change

## Questions?

Open an issue if you need help or have questions.
