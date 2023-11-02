# get-next-version

get-next-version gets the next version for your repository according to semantic versioning based on [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#specification).

## Installation

Go to the [releases page](https://github.com/thenativeweb/get-next-version/releases), find the download url for your architecture and operating system, and copy it.

Then, run the following steps:

```shell
# Download the latest release (insert the url here)
$ curl -L -o get-next-version <URL>

# Ensure the binary is executable
$ chmod a+x get-next-version

# Move the binary to the application directory
$ sudo mv get-next-version /usr/local/bin
```

## Quick Start

Go to the repository and run `get-next-version`. The tool will analyze the history of your repository and output the next version for your release.

```shell
$ get-next-version
```

Optionally, you may hand over the `--repository` (or short `-r`) flag to specify the path to the repository you want to analyze, if it is not in the current working directory.

```shell
$ get-next-version --repository <PATH>
```

If you need to prefix the version, you can use the `--prefix` (or short `-p`) flag. Note that the prefix must be a valid tag name on its own.

By default, output will be printed to the console in a human-readable format. If you want to print the output in a machine-readable format, you can use the `--target` (or short `-t`) flag:

```shell
# Print output in JSON format
$ get-next-version --target json

# Write output to the GITHUB_OUTPUT file in GitHub Action format (see https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-output-parameter)
$ get-next-version --target github-action
```

## Using the GitHub Action

For convenience, you may use the GitHub Action when running `get-next-version` inside a workflow on GitHub.

**⚠️ When cloning the repository, make sure to set the `fetch-depth` option to `0`, otherwise `get-next-version` will not be able to analyze the history of the repository!**

**⚠️ The action uses the parameter `target=github-action` by default, which will not print any human-readable output, but only write the output to the GITHUB_OUTPUT file.**

An example workflow that makes use of the GitHub Action is shown below:

```yaml
name: Example workflow

on: pull_request

jobs:
  example:
    name: Example
    runs-on: ubuntu-latest

    steps:
    - name: Clone repository
      uses: actions/checkout@v3
      with:
        fetch-depth: 0
        ref: ${{ github.event.pull_request.head.sha }}
    - name: Get next version
      id: get_next_version
      uses: thenativeweb/get-next-version@main
      with:
        prefix: 'v' # optional, defaults to ''
    - name: Show the next version
      run: |
        echo ${{ steps.get_next_version.outputs.version }}
        echo ${{ steps.get_next_version.outputs.hasNextVersion }}
```

## Using commit messages

In case you are not familiar with conventional commits (as mentioned above), here is a short summary. Basically, you should prefix your commit messages with one of the following keywords:

- `chore` – used for maintenance, does not result in a new version
- `fix` – used for bug fixes, results in a new patch version (e.g. from `1.2.3` to `1.2.4`)
- `feat` – used for introducing new features, results in a new minor version (e.g. from `1.2.3` to `1.3.0`)
- `feat!` – used for breaking changes, results in a new major version (e.g. from `1.2.3` to `2.0.0`)

Some examples for commit messages are shown below:

- `chore: Initial commit`
- `fix: Correct typo`
- `feat: Add support for Node.js 18`
- `feat!: Change API from v1 to v2`

Please note that `!` indicates breaking changes, and will always result in a new major version, independent of the type of change.
