#!/usr/bin/env bash

BRANCH=$(git branch --format="%(refname:lstrip=2)")
COMMIT=$(git log --pretty=format:'%h' | head -n 1)

IMAGE="terraform-openshift"
DEFAULT_TAG="latest"

printf "Current branch name: %s, short commit number: %s\n"  $BRANCH $COMMIT

podman build -t ${IMAGE}:${DEFAULT_TAG} -f image/Dockerfile .

podman tag ${IMAGE}:${DEFAULT_TAG} ${IMAGE}:${COMMIT}