#!/usr/bin/env sh

set -e

/action/get-next-version \
  --repository "$REPOSITORY_PATH" \
  --target github-action \
  --prefix "$INPUT_PREFIX" \
  --feature-prefixes "$INPUT_FEATURE_PREFIXES" \
  --fix-prefixes "$INPUT_FIX_PREFIXES" \
  --chore-prefixes "$INPUT_CHORE_PREFIXES"
