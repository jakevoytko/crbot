git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.7.0",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

go_rules_dependencies()
go_register_toolchains()

# Import Go dependencies.

go_repository(
    name = "com_github_bwmarrin_discordgo",
    importpath = "github.com/bwmarrin/discordgo",
    tag = "v0.17.0",
)

go_repository(
    name = "com_github_go_redis_redis",
    importpath = "github.com/go-redis/redis",
    tag = "v6.7.4"
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    commit = "b080dc9a8c480b08e698fb1219160d598526310f",
)

go_repository(
    name = "com_github_gorilla_websocket",
    importpath = "github.com/gorilla/websocket",
    commit = "23059f29570f0e13fca80ef6aea0f04c11daaa4d",
)

# Set up rules_docker
git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.3.0",
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
