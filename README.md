![Test Status](https://github.com/roverdotcom/checkbridge/workflows/Test/badge.svg)
## Checkbridge

Command-line utility to allow creating arbitrary [GitHub
checks](https://developer.github.com/v3/checks/) from other command-line
utilities.

![Screenshot](/.github/screenshots/go-lint-check-screenshot.png)

## About

If you're using GitHub actions, have a look at [Lint Action] instead, which
is both more full-featured and doesn't require you to create a GitHub app and
install it.

The use case this was designed for is when you have an existing CI system that
is **not** running via GitHub actions, and you'd like to use [GitHub checks].

GitHub checks allow you to post line-level annotations to files in Pull Requests
and commits, which is especially useful for linters and other code analysis
tools that you may want to run before code is merged.

In order to use GitHub checks on commits and pull requests, you need to have
a GitHub app provisioned and installed in your organization with `write` scope
on the `checks` [permission].

**TODO** show steps to create and install a GitHub app.

[github checks]: https://developer.github.com/v3/checks/
[lint action]: https://github.com/samuelmeuli/lint-action
[permission]: https://developer.github.com/v3/apps/permissions/#permission-on-checks

## Usage

_Coming soon_

## Configuration

All configuration can be passed as either command-line arguments or environment variables. All flags have shorthand values, run `checkbridge --help` to see them:

```
Flags:
  -a, --application-id int    GitHub application ID (numeric)
  -c, --commit-sha string     commit SHA to report status checks for
  -z, --exit-zero             exit zero even when tool reports issues
  -r, --github-repo string    GitHub repository (e.g. 'roverdotcom/checkbridge')
  -h, --help                  help for checkbridge
  -i, --installation-id int   GitHub installation ID (numeric)
  -p, --private-key string    GitHub application private key path or value
  -v, --verbose               verbose output
```

### Required flags

| Flag               | Environment Variable         |
| ------------------ | ---------------------------- |
| `--application-id` | `CHECKBRIDGE_APPLICATION_ID` |
| `--private-key`    | `CHECKBRIDGE_PRIVATE_KEY`    |

#### Optional flags

The following flags can be configured, but are not required by default.

| Flag                | Environment Variable          |
| ------------------- | ----------------------------- |
| `--installation-id` | `CHECKBRIDGE_INSTALLATION_ID` |
| `--commit-sha`      | `CHECKBRIDGE_COMMIT_SHA`      |
| `--github-repo`     | `CHECKBRIDGE_GITHUB_REPO`     |

### Defaults

`--installation-id` will be looked up dynamically if not provided, by doing a `GET` to
`/repos/:owner/:repo/installation` with the provided private key / application ID.

`--commit-sha` will be read from `$GITHUB_SHA`, `$BUILDKITE_COMMIT`, or `$(git rev-parse HEAD)`

`--github-repo` will be read from `$GITHUB_REPO` or `$BUILDKITE_REPOSITORY` if present

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
