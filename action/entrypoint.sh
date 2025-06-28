#!/usr/bin/env sh

set -e

/action/get-next-version \
  --repository /github/workspace \
  --target github-action \
  --prefix "$INPUT_PREFIX" \
  --feature-prefixes "$INPUT_FEATURE_PREFIXES" \
  --fix-prefixes "$INPUT_FIX_PREFIXES" \
  --chore-prefixes "$INPUT_CHORE_PREFIXES"
