#!/usr/bin/env sh

set -e

/action/get-next-version \
  --repository /github/workspace \
  --format github-action
