#!/usr/bin/env sh

set -e

NEXT_VERSION=$(/action/get-next-version /github/workspace)

git tag "${NEXT_VERSION}"
git push origin "${NEXT_VERSION}"