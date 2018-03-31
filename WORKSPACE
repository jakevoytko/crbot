http_archive(
    name = "io_bazel_rules_go",
    sha256 = "4b2c61795ac2eefcb28f3eb8e1cb2d8fb3c2eafa0f6712473bc5f93728f38758",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.10.2/rules_go-0.10.2.tar.gz",
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
    commit = "12892e8c234f4fe6f6803f052061de9057903bb2",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "eb925808374e5ca90c83401a40d711dc08c0c0f6",
    importpath = "github.com/gorilla/websocket",
)

# Import my own projects.
git_repository(
    name = "com_github_jakevoytko_go_stringmap",
    commit = "96db7d019a36e4cca914cd8e343d3ec8f1741271",
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
