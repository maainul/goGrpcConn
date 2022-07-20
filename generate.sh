#!/usr/bin/env bash

set -euxo pipefail

source ./devenv.sh

GOBIN=$api_gunk_dir/bin go install \
        github.com/gunk/gunk

$api_gunk_dir/bin/gunk generate ./...