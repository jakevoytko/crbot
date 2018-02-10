git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.9.0",
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
    tag = "v6.8.3",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "5119cf507ed5294cc409c092980c7497ee5d6fd2",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "4ac909741dfa57448bfadfdbca0cf7eeaa68f0e2",
    importpath = "github.com/gorilla/websocket",
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
