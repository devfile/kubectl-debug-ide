#!/usr/bin/env bash

pod=outyet
img=golang
new_pod=outyet-debug

kubectl debug -it $pod \
    --image $img \
    --share-processes \
    --copy-to=$new_pod
