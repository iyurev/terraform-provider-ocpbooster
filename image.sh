#!/usr/bin/env bash

BRANCH=$(git branch --format="%(refname:lstrip=2)")
COMMIT=""

IMAGE="terraform-openshift"
TAG="dev"

podman build -t ${IMAGE}:${TAG} -f image/Dockerfile .

