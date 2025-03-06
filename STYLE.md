# Style Guide

The eigenda repo consists majorly of Go and Solidity code, so we break it up into two sections:
- [Go Style Guide](#go-style-guide)
- [Solidity Style Guide](#solidity-style-guide)

Our implicit rule is to defer to each language's respective style guides (see [resources](#resources) section). 
However, some disagreements keep coming up in PRs, so we document our conventions here to prevent too much bike-shedding.
This style guide is thus by definition perpetually a work in progress.

## Go Style Guide

This document outlines the coding standards and best practices for our Go projects.

### Code Organization

Follow the [Ben Johnson](https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091) repo structuring method:
- Write interfaces in a file at the root of each package, and
- Write every implementation of that interface in its own package in a subdirectory.

#### Function Organization
- Keep functions small and focused
- Limit function length to ~50 lines where possible
- Place the most important functions at the top of the file
- Public functions that aren't methods should be placed in files with `_utils.go` suffix

### Error Handling

- We follow Uber's [error-wrapping style](https://github.com/uber-go/guide/blob/master/style.md#error-wrapping).

### Comments and Documentation

- Document all exported functions, types, and variables.
- Use complete sentences with proper punctuation.
- Keep comments concise and to the point. Avoid redundant comments that just restate the code. This means, in general, avoid LLM-generated comments. 
- Use your brain to bring in external context that only you have while writing the code. Code is often written once and read many times. Be mindful of the reader.

Here is an example of what NOT to do:
```go
// service for users
type UserService struct {
    db *sql.DB
}

// CreateUser creates a new user in the system.
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // ... implementation
}
```
Here is the same example with much more useful and complete comments. 
Note that its fine to have comments that are specific to the current implementation:
```go
// UserService handles all user-related operations.
//
// Safe for concurrent use by multiple goroutines.
type UserService struct {
    // db is the database connection used by the service.
    // It is currently hardcoded with a connectino pool of size 10.
    // Note: we use SQL for flexibility for now, but may switch to NoSQL once our query patterns are fixed.
    db *sql.DB
}

// CreateUser creates a new user in the sql database.
// The passed ctx should already have a deadline set, CreateUser does not add one.
// It returns an error if the user already exists.
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // ... implementation
}
```

# Solidity Style Guide

# Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Go Style Guide](https://google.github.io/styleguide/go/)
- [Go Proverbs](https://go-proverbs.github.io/)