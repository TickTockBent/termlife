# Contributing to termlife

Thanks for your interest in contributing!

## Reporting Issues

Found a bug or have a problem? [Open an issue](https://github.com/ticktockbent/termlife/issues/new) with:

- Your OS and terminal emulator
- Steps to reproduce the problem
- What you expected vs what happened
- Any error messages

## Feature Requests

Have an idea? [Open an issue](https://github.com/ticktockbent/termlife/issues/new) describing:

- What you'd like to see
- Why it would be useful
- Any implementation ideas (optional)

## Pull Requests

1. Fork the repo
2. Create a branch: `git checkout -b my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit with a clear message
6. Push and open a PR

### Code Guidelines

- Run `go fmt` before committing
- Add tests for new functionality
- Keep changes focused - one feature/fix per PR

### Building

```bash
go build -o termlife .
go test ./...
```

### Testing GIF Export

```bash
./termlife --gif 10 --pattern glider --size 20x20 --gif-out test.gif
```

## Questions?

Open an issue - happy to help!
