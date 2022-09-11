# Change repo

[![Latest version](https://img.shields.io/github/v/tag/daishe/change-repo?label=latest%20version&sort=semver)](https://github.com/daishe/change-repo/releases)
[![Latest release status](https://img.shields.io/github/workflow/status/daishe/change-repo/Release?label=release%20build&logo=github&logoColor=fff)](https://github.com/daishe/change-repo/actions/workflows/release.yaml)

[![Go version](https://img.shields.io/github/go-mod/go-version/daishe/change-repo?label=version&logo=go&logoColor=fff)](https://golang.org/dl/)
[![License](https://img.shields.io/github/license/daishe/change-repo)](https://github.com/daishe/change-repo/blob/master/LICENSE)

A simple CLI utility to change the current working directory to one containing Git repository.

## Usage

Just run the CLI and it will scan the provided directory recursively and present a list of all found Git repositories, allowing you to choose from that list:

```sh
change-repo .
```

When you pick from the list it will open a new shell in the selected location. You can also configure a set of standard locations using `CHANGE_REPO_PATH` environment variable:

```sh
export CHANGE_REPO_PATH="/path/to/location/1:/path/to/other/location"
```

and then invoke application without any arguments:

```sh
change-repo
```

it will first allow you to pick the directory containing Git repositories from the list of locations in `CHANGE_REPO_PATH` environment variable and then will recursively scan the selected directory and present the list of all Git repositories in it, allowing you to chose from that list.

That's it!

## Help

To get the complete list of all flags, use

```sh
change-repo --help
```

Most likely you will also be interested in creating alias

```sh
alias cr="change-repo"
```

## Options

Change-repo has several options that control its behavior:

- `--help` - display help message
- `--maxdepth` - controls recursion depth when scanning for Git repositories (default: 20)
- `--version` - display version and copyright information

## Configuration

Change-repo can be configured using the following environment variables:

- `CHANGE_REPO_PATH` - list of paths to directories containing Git repositories. Used when invoking change-repo without any arguments. When empty or not set defaults to current working directory.
- `CHANGE_REPO_SHELL` - when set, its value will be used instead of `SHELL` environment variable.
- `SHELL` - shell to invoke when changing directory.

## License

Change-repo is open-sourced software licensed under the [GNU General Public License version 3 or later](https://www.gnu.org/licenses/).
