git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.5.3",
)
load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "go_repository")

go_repositories()

# Import Go dependencies.

go_repository(
    name = "com_github_bwmarrin_discordgo",
    importpath = "github.com/bwmarrin/discordgo",
    commit = "0993a94b4e1c3291bed2047f583f34792269355c",
)

go_repository(
    name = "in_gopkg_redis_v5",
    importpath = "gopkg.in/redis.v5",
    commit = "a16aeec10ff407b1e7be6dd35797ccf5426ef0f0",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    commit = "88d0005bf4c3ec17306ecaca4281a8d8efd73e91",
)

go_repository(
    name = "com_github_gorilla_websocket",
    importpath = "github.com/gorilla/websocket",
    commit = "ea4d1f681babbce9545c9c5f3d5194a789c89f5b",
)
