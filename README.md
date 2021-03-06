# getignore

[![Build and Test](https://github.com/haroldadmin/getignore/actions/workflows/build-test.yml/badge.svg)](https://github.com/haroldadmin/getignore/actions/workflows/build-test.yml)

`getignore` is a command-line utility to fetch `.gitignore` files from [Github's Gitignore](https://www.github.com/github/gitignore) repository.

Start your new projects right, by including a `.gitignore` file from the very beginning!

## Usage

### Interactive

Run `getignore search`:

![interactive-search](./media/getignore.gif)

### Non-interactive

- Run `getignore get <gitignore-filename>` (eg. `getignore get Node.gitignore`).
- `getignore` will find the file with the matching name and append its contents to your `.gitignore` file

## Installation

### macOS

Use Homebrew:

1. Add the `getignore` tap to your Homebrew installation:

```shell
brew tap haroldadmin/getignore
```

2. Install `getignore`:

```shell
brew install getignore
```

### Other

You can find pre-compiled binaries for popular platforms on the [Releases](https://github.com/haroldadmin/getignore/releases) page.

### Building from source

- Install the [Go programming language](https://golang.org/)
- Run `go get github.com/haroldadmin/getignore`

### Automated Release Publication

There are open issues for making `getignore` available as [pre-compiled binaries](https://github.com/haroldadmin/getignore/issues/3) and also with [various](https://github.com/haroldadmin/getignore/issues/1) [package managers](https://github.com/haroldadmin/getignore/issues/2) automatically on every release.

## Contributions

`getignore` is a very small side-project, and I would continue to maintain it in my free time. If you would like to lend a hand by adding new features or fix bugs, feel free to open issues or pull requests.

## License

See [License](./LICENSE).
