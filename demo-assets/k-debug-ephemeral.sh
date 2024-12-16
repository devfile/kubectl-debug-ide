#!/usr/bin/env bash

pod=outyet
img=ghcr.io/l0rd/outyet-dev:latest
container=outyet

k debug -it $pod \
    --image $img \
    --target $container
