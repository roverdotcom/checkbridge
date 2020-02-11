## Checkbridge

Command-line utility to allow creating arbitrary [GitHub
checks](https://developer.github.com/v3/checks/) from other command-line
utilities.

## Usage

If you're using GitHub actions, this project is probably overkill for your
needs. Instead, take a look at [Lint Action], which is both more full-featured
and doesn't require you to create a GitHub app and install it.

The use case this was designed for is when you have an existing CI system that
is **not** running via GitHub actions, and you'd like to use [GitHub checks].

GitHub checks allow you to post line-level annotations to files in Pull Requests
and commits, which is especially useful for linters and other code analysis
tools that you may want to run before code is merged.

[github checks]: https://developer.github.com/v3/checks/
[lint action]: https://github.com/samuelmeuli/lint-action

_Coming soon_

## Development

This section is intended for developers of this tool.

You will need to have the [Go] toolchain installed. If it's correctly installed,
you'll have the `go` binary available on your path. If you don't have `go`
available, install it via your system's package manager. For example, on macOS:

```
brew install go
```

This project is currently developed and tested using Go version 1.13, which is
the latest public release. You'll need at least Go 1.11 to build as this project
uses the [new module system].

[go]: https://golang.org/
[new module system]: https://blog.golang.org/using-go-modules

### Running tests

```sh
go test ./...
```

### Linting

Linting is enforced in CI via GitHub actions on this repository. If you get a
lint failure, you probably need to configure your editor for Go support.

Your editor should already be running `gofmt` (or `goimports`) for you on save.
If not, you can run it manually:

```sh
gofmt -w .
```

We also run `golint` to find common code issues. Run it with:

```sh
golint ./...
```

If you don't have `golint` available, install it with:

```sh
go get -u golang.org/x/lint/golint
```

### Related

- [Lint Action](https://github.com/samuelmeuli/lint-action)
- [Probot](https://probot.github.io/)
