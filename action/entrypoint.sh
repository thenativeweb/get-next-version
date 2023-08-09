#!/usr/bin/env sh

set -e

/action/get-next-version \
  --repository /github/workspace \
  --target github-action \
  --prefix "$INPUT_PREFIX"
