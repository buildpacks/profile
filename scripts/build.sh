#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
GOOS="linux" go build -ldflags='-s -w' -o bin/main github.com/buildpacks/profile/cmd/main

if [ "${STRIP:-false}" != "false" ]; then
  strip bin/main
fi

if [ "${COMPRESS:-none}" != "none" ]; then
  $COMPRESS bin/main
fi

ln -fs main bin/build
ln -fs main bin/detect
cp ${SCRIPT_DIR}/execd_wrapper.sh bin/execd_wrapper.sh
