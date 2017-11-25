#!/bin/bash

set -eo pipefail

# This will load the image and then fail due to a bazel bug in argument parsing.
bazel run --cpu k8 :crbot_image -- --filename "/secret.json" --localhost docker.for.mac.localhost

# This will run the image that was loaded in the last step.
docker run -i --rm bazel:crbot_image --filename "/secret.json" --localhost docker.for.mac.localhost
