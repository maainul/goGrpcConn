#!/usr/bin/env bash

readonly api_gunk_dir=$(git rev-parse --show-toplevel)/api/gunk

PATH=$api_gunk_dir/bin:$PATH

export PATH
