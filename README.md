# getignore

`getignore` is a command-line utility to fetch `.gitignore` files from [Github's Gitignore](https://www.github.com/github/gitignore) repository.

Start your new projects right, by including a `.gitignore` file from the very beginning!

## Usage

- Interactive: Run `getignore`, search for a file, and you're done!

![interactive-search](./media/interactive-search.gif)

- Non-interactive: Run `getignore --search <query>`. 
- `getignore` will automatically pick the best matching search result and append its contents to your `.gitignore` file

![non-interactive search](./media/getignore-non-interactive.png)

## Installation

- Install the [Go programming language](https://golang.org/)
- Run `go get github.com/haroldadmin/getignore`

## Contributions

`getignore` is a very small side-project, and I would continue to maintain it in my free time. If you would like to lend a hand by adding new features or fix bugs, feel free to open issues or pull requests.

## License

See [License](./LICENSE).