#!/bin/bash

set -eo pipefail

bazel run --cpu k8 :crbot_image -- --filename "/secret.json" --localhost docker.for.mac.localhost
