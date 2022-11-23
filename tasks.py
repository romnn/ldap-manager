from pathlib import Path
from invoke import task

PKG = "github.com/romnn/ldap-manager"
CMD_PKG = "github.com/romnn/ldap-manager/cmd/ldap-manager"


ROOT_DIR = Path(__file__).parent
BUILD_DIR = ROOT_DIR / "build"
WEB_DIR = ROOT_DIR / "web"


@task
def format(c):
    """Format code"""
    c.run("pre-commit run go-fmt --all-files")
    c.run("pre-commit run go-imports --all-files")


@task
def embed(c):
    """Embeds the examples"""
    c.run(f"npx embedme {ROOT_DIR / 'README.md'}")


@task
def test(c):
    """Run tests"""
    cmd = [
        "richgo",
        "test",
        "-race",
        "-coverpkg=all",
        "-coverprofile=coverage.txt",
        "-covermode=atomic",
        "./...",
    ]
    c.run(" ".join(cmd))


@task
def cyclo(c):
    """Check code complexity"""
    c.run("pre-commit run go-cyclo --all-files")


@task
def lint(c):
    """Lint code"""
    c.run("pre-commit run go-lint --all-files")
    c.run("pre-commit run go-vet --all-files")


@task
def install_hooks(c):
    """Install pre-commit hooks"""
    c.run("pre-commit install")


@task
def install_hooks(c):
    """Install pre-commit hooks"""
    c.run("pre-commit install")


@task
def pre_commit(c):
    """Run all pre-commit checks"""
    c.run("pre-commit run --all-files")


@task
def compile_go_protos(c):
    """Compile golang proto files"""
    import shutil
    from pprint import pprint

    out_dir = ROOT_DIR / "pkg" / "grpc" / "gen"
    try:
        shutil.rmtree(out_dir)
    except FileNotFoundError:
        pass
    out_dir.mkdir(parents=True, exist_ok=True)

    services = [
        ROOT_DIR / "proto" / "ldap_manager.proto",
    ]
    for service in services:
        proto_path = service.parent
        print(
            f"compiling {service.relative_to(ROOT_DIR)} "
            + f"to {out_dir.relative_to(ROOT_DIR)}"
        )
        package = (
            f"{service.relative_to(proto_path)}="
            + f"github.com/{out_dir.relative_to(ROOT_DIR.parent)}"
        )
        cmd = [
            "protoc",
            f"--proto_path={proto_path}",
            f"--go_opt=M{package}",
            f"--go-grpc_opt=M{package}",
            f"--grpc-gateway_opt=M{package}",
            f"--go_out={out_dir}",
            f"--go-grpc_out={out_dir}",
            f"--grpc-gateway_out={out_dir}",
            "--go_opt=paths=source_relative",
            "--go-grpc_opt=paths=source_relative",
            "--grpc-gateway_opt=logtostderr=true,paths=source_relative",
            str(service),
        ]
        if False:
            pprint(cmd)
        c.run(" ".join(cmd))


@task
def compile_ts_protos(c):
    """Compile typescript proto files"""
    import shutil
    from pprint import pprint

    out_dir = WEB_DIR / "generated" / "src"
    print(out_dir)
    try:
        shutil.rmtree(out_dir)
    except FileNotFoundError:
        pass
    out_dir.mkdir(parents=True, exist_ok=True)

    services = [
        ROOT_DIR / "proto" / "ldap_manager.proto",
    ]

    for service in services:
        proto_path = service.parent
        print(
            f"compiling {service.relative_to(ROOT_DIR)} "
            + f"to {out_dir.relative_to(ROOT_DIR)}"
        )
        plugin_path = WEB_DIR / "node_modules" / ".bin" / "protoc-gen-ts_proto"
        cmd = [
            "protoc",
            f"--plugin={plugin_path}",
            f"--proto_path={proto_path}",
            f"--ts_proto_out={out_dir}",
            str(service),
        ]
        if True:
            pprint(cmd)
        c.run(" ".join(cmd))

    if False:
        # rebuild the ldap manager package
        # this should no longer be required
        package_dir = WEB_DIR / "generated"
        c.run(f"cd {package_dir} && yarn build")
        c.run(f"cd {WEB_DIR} && yarn upgrade ldap-manager --force --latest")


@task(pre=[compile_go_protos, compile_ts_protos])
def compile_protos(c):
    """Compiles protos"""
    pass


# @task
# def cc(c):
#     """Build the project for all architectures"""

#     output = "{{.Dir}}-{{.OS}}-{{.Arch}}"
#     TRAVIS_TAG = os.environ.get("TRAVIS_TAG")
#     BINARY = os.environ.get("BINARY")
#     if TRAVIS_TAG and BINARY:
#         output = "%s-%s-{{.OS}}-{{.Arch}}" % (BINARY, TRAVIS_TAG)

#     # FIXME: compiling github.com/docker/docker/pkg/system
# on windows fails, so windows is disabled for now
#     c.run(
#         'gox -os="linux darwin" -arch="amd64" -output="build/%s"'
# + '-ldflags "-X main.Rev=`git rev-parse --short HEAD`" -verbose %s'
#         % (output, CMD_PKG)
#     )


@task
def build(c):
    """Build the project"""
    c.run("pre-commit run go-build --all-files")


@task
def lint_chart(c):
    """Lints the helm chart"""
    c.run("pre-commit run lint-chart --all-files")


@task
def run(c):
    """Run the cmd"""
    import sys

    options = sys.argv[3:]
    c.run("go run {} {}".format(CMD_PKG, " ".join(options)))


@task
def clean_build(c):
    """Clean up files from package building """
    c.run("rm -fr build/")


@task
def clean_coverage(c):
    """Clean up files from coverage measurement """
    c.run("find . -name 'coverage.txt' -exec rm -fr {} +")


@task(pre=[clean_build, clean_coverage])
def clean(c):
    """Runs all clean sub-tasks """
    pass
