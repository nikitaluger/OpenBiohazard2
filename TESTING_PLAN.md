# Testing Plan

## Goals
- Add useful tests without depending on copyrighted game data.
- Prefer hermetic tests that run in CI.
- Keep each test change small and reviewable.

## Strategy
1. Start with pure logic helpers:
   - `fileio/streamreader.go`
   - `fileio/bitreader.go`
2. Add hermetic parser tests using `*Stream` functions and in-memory bytes.
3. Keep file-based tests only as optional local smoke tests.
4. Add one parser or helper per change.

## Hermetic Test Pattern
```go
data := []byte{ /* minimal valid input */ }
r := bytes.NewReader(data)
out, err := LoadSomeStream(r, int64(len(data)))
```

## Rules
- Use `t.TempDir()` for any file-based tests.
- Skip only when external data is genuinely optional.
- Avoid global state and shared temp files.
- Run `go test ./...` before committing.

## Suggested Order
1. streamreader
2. bitreader
3. TIM
4. PLD
5. EMD / MD1 / ESP
