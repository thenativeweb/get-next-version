#!/usr/bin/env sh

apt update
apt install git

git status
git log

/action/get-next-version /github/workspace