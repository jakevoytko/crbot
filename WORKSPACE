http_archive(
    name = "io_bazel_rules_go",
    sha256 = "1868ff68d6079e31b2f09b828b58d62e57ca8e9636edff699247c9108518570b",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.11.1/rules_go-0.11.1.tar.gz",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

go_rules_dependencies()

go_register_toolchains()

# Import Go dependencies.

go_repository(
    name = "com_github_bwmarrin_discordgo",
    importpath = "github.com/bwmarrin/discordgo",
    tag = "v0.18.0",
)

go_repository(
    name = "com_github_go_redis_redis",
    importpath = "github.com/go-redis/redis",
    tag = "v6.10.2",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "4ec37c66abab2c7e02ae775328b2ff001c3f025a",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "21ab95fa12b9bdd8fecf5fa3586aad941cc98785",
    importpath = "github.com/gorilla/websocket",
)

# Import my own projects.
git_repository(
    name = "com_github_jakevoytko_go_stringmap",
    commit = "73c758d8e93882c7203fb171e4e9986e0514b575",
    remote = "https://github.com/jakevoytko/go-stringmap.git",
)

# Set up rules_docker
git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.4.0",
)

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()
