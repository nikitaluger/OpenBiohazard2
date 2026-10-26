# Contributing to OpenBiohazard2

## What belongs in the repository
- Go source code and tests
- Documentation and setup notes
- Build/configuration files
- Small, hermetic test fixtures

## What stays local
- Copyrighted game data (`data/`)
- Personal notes and debug dumps
- LLM/agent tooling and prompts
- Secrets and environment files

## Basic workflow
1. Fork the repository.
2. Create a small, focused branch.
3. Run `gofmt`, `go vet ./...`, and `go test ./...`.
4. Open a pull request referencing the related issue.

## Testing guidelines
- Prefer hermetic tests that do not require game data.
- Use `t.TempDir()` for file-based tests.
- Skip gracefully only when external data is genuinely optional.
