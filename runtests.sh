#!/bin/bash

set -eo pipefail

bazel build //...:all
bazel test //...:all