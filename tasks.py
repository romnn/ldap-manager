"""
Tasks for maintaining the project.
Execute 'invoke --list' for guidance on using Invoke
"""
import shutil
import pprint
import sys
import os

from invoke import task
import webbrowser

PKG = "github.com/romnnn/ldap-manager"
CMD_PKG = "github.com/romnnn/ldap-manager/cmd/ldap-manager"



ROOT_DIR = os.path.dirname(os.path.realpath(__file__))
BUILD_DIR = os.path.join(ROOT_DIR, "build")
TRAVIS_CONFIG_FILE = os.path.join(ROOT_DIR, ".travis.yml")


def _delete_file(file):
    try:
        file.unlink(missing_ok=True)
    except TypeError:
        # missing_ok argument added in 3.8
        try:
            file.unlink()
        except FileNotFoundError:
            pass


@task
def format(c):
    """Format code
    """
    c.run("pre-commit run go-fmt --all-files")
    c.run("pre-commit run go-imports --all-files")


@task
def test(c):
    """Run tests
    """
    c.run("env GO111MODULE=on go test -v -race ./...")


@task
def cyclo(c):
    """Check code complexity
    """
    c.run("pre-commit run go-cyclo --all-files")


@task
def lint(c):
    """Lint code
    """
    c.run("pre-commit run go-lint --all-files")
    c.run("pre-commit run go-vet --all-files")


@task
def install_hooks(c):
    """Install pre-commit hooks
    """
    c.run("pre-commit install")


@task
def pre_commit(c):
    """Run all pre-commit checks
    """
    c.run("pre-commit run --all-files")


@task(help=dict(publish="Publish the coverage result to codecov.io (default False)",),)
def coverage(c, publish=False):
    """Create coverage report
    """
    c.run(
        "env GO111MODULE=on go test -v -race -coverprofile=coverage.txt -coverpkg=all -covermode=atomic ./..."
    )
    if publish:
        # Publish the results via codecov
        c.run("bash <(curl -s https://codecov.io/bash)")


@task
def compile_proto(c):
    """Build the project
    """
    # Generate golang grpc server and client stubs
    c.run(str(" ").join([
        "protoc",
        "--proto_path=%s" % ROOT_DIR,
        "--go_out=grpc/ldap-manager",
        "--go-grpc_out=grpc/ldap-manager",
        "--grpc-gateway_out=grpc/ldap-manager",
        "--go_opt=paths=source_relative",
        "--go-grpc_opt=paths=source_relative",
        "--go-grpc_opt=paths=source_relative",
        "--grpc-gateway_opt=logtostderr=true,paths=source_relative",
        os.path.join(ROOT_DIR, "ldap_manager.proto"),
    ]))
    

@task
def cc(c):
    """Build the project for all architectures
    """
    
    output = "{{.Dir}}-{{.OS}}-{{.Arch}}"
    TRAVIS_TAG = os.environ.get("TRAVIS_TAG")
    BINARY = os.environ.get("BINARY")
    if TRAVIS_TAG and BINARY:
        output = "%s-%s-{{.OS}}-{{.Arch}}" % (BINARY, TRAVIS_TAG)
    
    # FIXME: compiling github.com/docker/docker/pkg/system on windows fails, so windows is disabled for now
    c.run(
        'gox -os="linux darwin" -arch="amd64" -output="build/%s" -ldflags "-X main.Rev=`git rev-parse --short HEAD`" -verbose %s'
        % (output, CMD_PKG)
    )



@task
def build(c):
    """Build the project
    """
    c.run("pre-commit run go-build --all-files")


@task
def lint_chart(c):
    """Lints the helm chart
    """
    c.run("pre-commit run lint-chart --all-files")


@task
def run(c):
    """Run the cmd target
    """
    options = sys.argv[3:]
    c.run("go run {} {}".format(CMD_PKG, " ".join(options)))



@task
def clean_build(c):
    """Clean up files from package building
    """
    c.run("rm -fr build/")


@task
def clean_coverage(c):
    """Clean up files from coverage measurement
    """
    c.run("find . -name 'coverage.txt' -exec rm -fr {} +")


@task(pre=[clean_build, clean_coverage])
def clean(c):
    """Runs all clean sub-tasks
    """
    pass


def _create(d, *keys):
    current = d
    for key in keys:
        try:
            current = current[key]
        except (TypeError, KeyError):
            current[key] = dict()
            current = current[key]


def _fix_token(config_file=None, force=False, verify=True):
    from ruamel.yaml import YAML
    yaml = YAML()
    config_file = config_file or TRAVIS_CONFIG_FILE
    with open(config_file, "r") as _file:
        try:
            travis_config = yaml.load(_file)
        except Exception:
            raise ValueError(
                "Failed to parse the travis configuration. "
                "Make sure the config only contains valid YAML and keys as specified by travis."
            )

        # Get the generated token from the top level deploy config added by the travis cli
        try:
            real_token = travis_config["deploy"]["api_key"]["secure"]
        except (TypeError, KeyError):
            raise AssertionError("Can't find any top level deployment tokens")

        try:
            # Find the build stage that deploys to releases
            releases_stages = [
                stage
                for stage in travis_config["jobs"]["include"]
                if stage.get("deploy", dict()).get("provider") == "releases"
            ]
            assert (
                len(releases_stages) > 0
            ), "Can't set the new token because there are no stages deploying to releases"
            assert (
                len(releases_stages) < 2
            ), "Can't set the new token because there are multiple stages deploying to releases"
        except (TypeError, KeyError):
            raise AssertionError(
                "Can't set the new token because there are no deployment stages")

        try:
            is_mock_token = releases_stages[0]["deploy"]["token"]["secure"] == "REPLACE_ME"
            is_same_token = releases_stages[0]["deploy"]["token"]["secure"] == real_token

            unmodified = is_mock_token or is_same_token
        except (TypeError, KeyError):
            unmodified = False

        # Set the new generated token as the stages deploy token
        _create(releases_stages[0], "deploy", "token", "secure")
        releases_stages[0]["deploy"]["token"]["secure"] = real_token

        # Make sure it is fine to overwrite the config file
        assert unmodified or force, (
            'The secure token in the "{}" stage has already been changed. '
            "Retry with --force if you are sure about replacing it.".format(
                releases_stages[0].get("stage", "releases deployment")
            )
        )

        # Remove the top level deploy config added by the travis cli
        travis_config.pop("deploy")

        if not unmodified and verify:
            pprint.pprint(travis_config)
            if (
                not input("Do you want to save this configuration? (y/n) ")
                .strip()
                .lower()
                == "y"
            ):
                return

    # Save the new travis config
    assert travis_config
    with open(config_file, "w") as _file:
        yaml.dump(travis_config, _file)
    print("Fixed!")


@task(help=dict(
    force="Force overriding the current travis configuration",
    verify="Verify config changes by asking for the user's approval"
))
def fix_token(c, force=False, verify=True):
    """
    Add the token generated by the travis cli script to the correct entry
    """
    _fix_token(force=force, verify=verify)
