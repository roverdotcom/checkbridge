## Checkbridge ![build status](https://github.com/roverdotcom/checkbridge/workflows/Test/badge.svg)

**Project Status**: Alpha

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

[github checks]: https://developer.github.com/v3/checks/
[lint action]: https://github.com/samuelmeuli/lint-action
[permission]: https://developer.github.com/v3/apps/permissions/#permission-on-checks

## Installing

Precompiled binaries are available for all releases [on GitHub]. Because they are
static binaries, you can simply download them and run them. For example, on Linux:

```bash
curl -L https://github.com/roverdotcom/checkbridge/releases/download/v0.1.0/checkbridge-0.1.0.linux-amd64.tar.gz \
    | tar zxf - -C /usr/local/bin
```

Would install `checkbridge` release `v0.1.0` into `/usr/local/bin/`

You can also install from source if you have Go 1.11+ installed:

```bash
go get github.com/roverdotcom/checkbridge
```

This will install the tip of `master` (not recommended for production usage) into `$GOPATH/bin/`

[on github]: https://github.com/roverdotcom/checkbridge/releases

## Usage

`checkbridge` requires GitHub credentials to report checks. Read through the [configuration]
and [authentication] sections to see get started. Once configured, you can create checks by
piping your desired tool into a `checkbridge` subcommand. For example:

```bash
golint ./... | checkbridge golint
```

[configuration]: #configuration
[authentication]: #authentication

## Configuration

Most configuration options can be passed as either command-line arguments or
environment variables. All flags have shorthand values, run `checkbridge --help`
to see them:

```
Flags:
  -o, --annotate-only         only leave annotations, never mark check as failed
  -a, --application-id int    GitHub application ID (numeric)
  -c, --commit-sha string     commit SHA to report status checks for
  -d, --details-url string    details URL to send for check
  -z, --exit-zero             exit zero even when tool reports issues
  -r, --github-repo string    GitHub repository (e.g. 'roverdotcom/checkbridge')
  -h, --help                  help for checkbridge
  -i, --installation-id int   GitHub installation ID (numeric)
  -m, --mark-in-progress      mark check as in progress before parsing
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
| `--github-token`    | `CHECKBRIDGE_GITHUB_TOKEN`    |

### Defaults

`--installation-id` will be looked up dynamically if not provided, by doing a `GET` to
`/repos/:owner/:repo/installation` with the provided private key / application ID.

`--commit-sha` will be read from `$GITHUB_SHA`, `$BUILDKITE_COMMIT`, or `$(git rev-parse HEAD)`

`--github-repo` will be read from `$GITHUB_REPO` or `$BUILDKITE_REPOSITORY` if present

`--github-token` will be read from `$GITHUB_TOKEN` if present (i.e. when run via GitHub actions)

## Authentication

Using the GitHub checks API requires a GitHub app to be created and installed, with `checks`
permission on the repo you're running against.

GitHub application tokens, unlike GitHub user tokens, are short-lived (1 hour) and thus must
usually be fetched for each run. When running via GitHub actions, you automatically have a
token (available under `secrets.github_token`) available for the "GitHub actions" application
for your use. You can see an example of this in this repo's [lint.yml](./.github/workflows/lint.yml)
workflow.

If you're running via GitHub actions, simply provide the token (as `--github-token`,
`$GITHUB_TOKEN`, or `$CHECKBRIDGE_GITHUB_TOKEN`) and you're good to go. If you're _not_ running via
GitHub actions (for example, you're using [BuildKite](https://buildkite.com/)), read on.

### Creating a GitHub app

First, you'll need to [create a GitHub app].

You'll start with [registering your app]. Give it a descriptive name and a homepage (these are
required fields). We won't be doing any OAuth (user authorization) so you can leave most of the
rest blank. Under the "Permissions" section, you'll need to select "Read & write" access
to the "Checks" permission.

If you're going to use your application on an organization (as opposed to just your own repos),
select "Any account" under "Where can this GitHub App be installed?"

Once created, you'll need to grab two pieces of information:

1. The application ID, which will be at the very top of the page ("App ID")
2. A private key, at the very bottom of the page. You can generate a new one when your
   app is first created. This will prompt you to download a `.pem` file containing the private key.

[create a github app]: https://developer.github.com/apps/building-github-apps/creating-a-github-app/
[registering your app]: https://github.com/settings/apps/new

After creating the app and saving the app's ID and private key, you'll need to install it. On the
lefthand side of the app's detail page, click "Install App" and select the organization (or your
own account) you'd like to use. This will prompt you to accept the new application. Verify that
it looks correct, and click "Install".

Installing an application generates an installation ID, which will be the final part of the URL on
the page you're sent to (i.e. `https://github.com/settings/installations/1234` where `1234` is your)
installation ID. Installation IDs represent an instance of an application being installed in an
organization or user account.

While saving and using the installation ID is optional, it saves an extra API call every time
`checkbridge` is run.

The final step is to make the application ID, private key, and optionally the installation ID,
available to `checkbridge`. Do _not_ commit the private key to your repository.

You should verify your credentials are correct by running `checkbridge check-auth`:
For example:

```bash
# These could be configured in your CI environment
export CHECKBRIDGE_APPLICATION_ID="456"
export CHECKBRIDGE_PRIVATE_KEY="/tmp/private_key.pm"
export CHECKBRIDGE_INSTALLATION_ID="1234"  # application 456 installed in "myorg"

checkbridge check-auth --github-repo=myorg/myrepo
```

Or specifying as command-line flags:

```bash
# --installation-id left off, will look it up dynamically
checkbridge check-auth \
  --github-repo=myorg/myrepo \
  --application-id=456 \
  --private-key=/tmp/private_key.pem
```

If it returns an error, validate you've passed the correct configuration values. If it returns
success, you're ready to use `checkbridge`.

## Available parsers

Currently, `checkbridge` has builtin support for [golint] and [mypy]. In addition, it has a generic
`regex` command, which allows you to specify a regular expression. For example, running the
following would create an annotation on `example.go` line `1`, with the message `this is a message`.

```bash
echo "example.go:1: this is a message" | checkbridge regex \
  --regex "(.*):(.*): (.*)" \
  --name "my custom linter" \
  --path-pos 1 \
  --line-pos 2 \
  --message-pos 3
```

Note that the positions start at 1, as per convention, where the 0th element is the whole string
match.

Run `checkbridge regex --help` to see all the available configuration options.

[golint]: https://github.com/golang/lint
[mypy]: https://mypy.readthedocs.io/

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
