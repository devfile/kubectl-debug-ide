#!/usr/bin/env bash

pod=outyet
img=ghcr.io/l0rd/outyet-dev:latest
new_pod=outyet-debug
repo=https://github.com/l0rd/outyet.git

kubectl debug-ide $pod \
                  --image $img \
                  --share-processes \
                  --copy-to=$new_pod \
                  --git-repository $repo
