#!/usr/bin/env sh

set -e

NEXT_VERSION=$(/action/get-next-version -r /github/workspace)

git tag "${NEXT_VERSION}"
git push origin "${NEXT_VERSION}"

echo "::set-output name=version::${NEXT_VERSION}"
