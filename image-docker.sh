#!/usr/bin/env bash

BRANCH=$(git branch --format="%(refname:lstrip=2)")
COMMIT=$(git log --pretty=format:'%h' | head -n 1)

IMAGE="terraform-openshift"
DEFAULT_TAG="latest"

printf "Current branch name: %s, short commit number: %s\n"  $BRANCH $COMMIT

docker build -t ${IMAGE}:${DEFAULT_TAG} -f image/Dockerfile .

docker tag ${IMAGE}:${DEFAULT_TAG} ${IMAGE}:${COMMIT}
