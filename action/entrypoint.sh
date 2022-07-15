#!/usr/bin/env sh

apt update
apt install -y git

git status
git log

/action/get-next-version /github/workspace