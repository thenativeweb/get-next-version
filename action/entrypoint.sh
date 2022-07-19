#!/usr/bin/env sh

set -e

set +e
NEXT_VERSION=$(/action/get-next-version -r /github/workspace)
GET_NEXT_VERSION_EXIT_CODE=$?
set -e

if [ ${GET_NEXT_VERSION_EXIT_CODE} -ne 0 ]; then
  if [ ${GET_NEXT_VERSION_EXIT_CODE} -eq 2 ]; then
    echo "::set-output name=version::${NEXT_VERSION}"
    echo "::set-output name=hasNextVersion::false"
    exit 0
  fi

  exit ${GET_NEXT_VERSION_EXIT_CODE}
fi

echo "::set-output name=version::${NEXT_VERSION}"
echo "::set-output name=hasNextVersion::true"
