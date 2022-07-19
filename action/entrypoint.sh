#!/usr/bin/env sh

set -e

NEXT_VERSION=$(/action/get-next-version -r /github/workspace)

echo "::set-output name=version::${NEXT_VERSION}"
