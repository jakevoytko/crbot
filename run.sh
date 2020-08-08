#!/bin/bash

set -eo pipefail

bazel run --cpu k8 :crbot_image -- --norun
bazel run //deploy:crbot.replace
