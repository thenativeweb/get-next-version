#!/usr/bin/env sh

echo "[DEBUG] ls -lah ."
ls -lah .

echo "[DEBUG] ls -lah /github/workspace"
ls -lah .

echo "[DEBUG] ./get-next-version /github/workspace"
./get-next-version /github/workspace