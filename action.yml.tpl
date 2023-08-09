name: 'get-next-version'
description: >
  Gets the next version for your repository according to
  semantic versioning based on conventional commits.
inputs:
  prefix:
    description: 'Sets the version prefix'
    required: false
    default: ''
outputs:
  version:
    description: 'Next version'
  hasNextVersion:
    description: 'Whether there is a next version'
runs:
  using: 'docker'
  image: '<docker-image>'
