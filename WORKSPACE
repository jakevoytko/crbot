git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.6.0",
)
load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "go_repository")

go_repositories()

# Import Go dependencies.

go_repository(
    name = "com_github_bwmarrin_discordgo",
    importpath = "github.com/bwmarrin/discordgo",
    tag = "v0.17.0",
)

go_repository(
    name = "com_github_go_redis_redis",
    importpath = "github.com/go-redis/redis",
    tag = "v6.7.2"
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    commit = "9419663f5a44be8b34ca85f08abc5fe1be11f8a3",
)

go_repository(
    name = "com_github_gorilla_websocket",
    importpath = "github.com/gorilla/websocket",
    commit = "71fa72d4842364bc5f74185f4161e0099ea3624a",
)
