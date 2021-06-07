#!/usr/bin/env bash
IMAGE="terraform-openshift"
TAG="dev"

podman build -t ${IMAGE}:${TAG} -f image/Dockerfile .

