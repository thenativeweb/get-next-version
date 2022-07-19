# get-next-version

get-next-version gets the next version for your repository according to SemVer (semantic versioning) based on conventional commits.

## Installation

Go to the [releases page](https://github.com/thenativeweb/get-next-version/releases), find the download url for your architecture and operating system, and copy it.

Then, run the following steps:

```shell
# Download the latest release (insert the url here)
$ curl -o get-next-version <URL>

# Ensure the binary is executable
$ chmod a+x get-next-version

# Move the binary to the application directory
$ sudo mv get-next-version /usr/local/bin
```

## Quick Start

Go to the repository and run `get-next-version`. The tool will analyse the history of your repository and output the next version for your release.

```shell
$ ./get-next-version
```
