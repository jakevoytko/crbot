load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

################################################################################
# rules_go

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "b27e55d2dcc9e6020e17614ae6e0374818a3e3ce6f2024036e688ada24110444",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/v0.21.0/rules_go-v0.21.0.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.21.0/rules_go-v0.21.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

################################################################################
# gazelle

http_archive(
    name = "bazel_gazelle",
    sha256 = "86c6d481b3f7aedc1d60c1c211c6f76da282ae197c3b3160f54bd3a8f847896f",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.1/bazel-gazelle-v0.19.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.1/bazel-gazelle-v0.19.1.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

################################################################################
# go dependencies

go_repository(
    name = "com_github_bwmarrin_discordgo",
    importpath = "github.com/bwmarrin/discordgo",
    tag = "v0.20.0",
)

go_repository(
    name = "com_github_go_redis_redis",
    importpath = "github.com/go-redis/redis",
    tag = "v6.15.6",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "530e935923ad688be97c15eeb8e5ee42ebf2b54a",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "c3e18be99d19e6b3e8f1559eea2c161a665c4b6b",
    importpath = "github.com/gorilla/websocket",
)

# Import my own projects.
go_repository(
    name = "com_github_jakevoytko_go_stringmap",
    commit = "f04a23f25e90807b0b32dbb1d141ec6b60ad4353",
    importpath = "github.com/jakevoytko/go-stringmap",
)

################################################################################
# rules_docker

# Download the rules_docker repository at release v0.14.1
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "dc97fccceacd4c6be14e800b2a00693d5e8d07f69ee187babfd04a80a9f8e250",
    strip_prefix = "rules_docker-0.14.1",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.1/rules_docker-v0.14.1.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()
