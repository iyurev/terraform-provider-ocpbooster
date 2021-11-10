#!/usr/bin/env bash

OCP_VERSION="4.9.5"
CACHE_DIR=/tmp/openshift-installer-dist/
INSTALLER_ASSET="${PWD}/pkg/booster/assets/"

printf "Create cache directory if it doesn't exist\n"
mkdir -p $CACHE_DIR

for platform in linux mac; do

if [[ ! -f ${CACHE_DIR}/openshift-install-${platform}.tar.gz ]]; then
INSTALLER_URL="https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/openshift-install-${platform}-${OCP_VERSION}.tar.gz"
printf "Download openshift-install binary from %s, for %s, to the cache directory: %s\n" $INSTALLER_URL  $platform  $CACHE_DIR
curl --fail -L  -o ${CACHE_DIR}/openshift-install-${platform}.tar.gz  ${INSTALLER_URL}
if [[ $? != 0 ]];then
  printf "Failed for download installer binary file\n"
  exit 1
  fi
fi
done

printf "Create assets directory if it doesn't exist\n"
mkdir -p ${PWD}/pkg/booster/assets/
printf "OK\n"

printf "Copy openshift-install archives to asset directory\n"
cp ${CACHE_DIR}/*.tar.gz  $INSTALLER_ASSET


