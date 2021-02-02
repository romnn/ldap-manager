## ldap-manager

[![Build Status](https://github.com/romnn/ldap-manager/workflows/test/badge.svg)](https://github.com/romnn/ldap-manager/actions)
[![GitHub](https://img.shields.io/github/license/romnn/ldap-manager)](https://github.com/romnn/ldap-manager)
 [![Docker Pulls](https://img.shields.io/docker/pulls/romnn/ldap-manager)](https://hub.docker.com/r/romnn/ldap-manager) [![Test Coverage](https://codecov.io/gh/romnn/ldap-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/romnn/ldap-manager)
[![Release](https://img.shields.io/github/release/romnn/ldap-manager)](https://github.com/romnn/ldap-manager/releases/latest)

<p align="center">
  <img width="200" src="public/icon/icon_lg.jpg">
</p>

LDAP Manager is the cloud-native LDAP web management interface. LDAP has been around for a long time and has become a popular choice for user and group management - however, this should not mean that it's management interface should be hard to deploy and look and feel like it was made in the last century.

LDAP Manager is written in Go and comes with a Vue/Typescript frontend in a single, self-contained docker container. It also exposes it's API over both REST and gRPC!

| | |
|:-------------------------:|:-------------------------:|
| <img src="screenshots/home-user.png"> | <img src="screenshots/accounts-edit-admin.png"> |
| <img src="screenshots/accounts-list-admin.png"> | <img src="screenshots/groups-edit-admin.png"> |

Before you get started, make sure you have an OpenLDAP server like [osixia/openldap](https://hub.docker.com/r/osixia/openldap/) running. For more information on deployment and a full example, see the [deployment guide](#Deployment).

```bash
go run github.com/romnn/ldap-manager/cmd/ldap-manager serve \
    --http-port 8080 \
    --grpc-port 9090 \
    --generate
```

You can also download pre-built binaries from the [releases page](https://github.com/romnn/ldap-manager/releases), or use the `docker` image:

```bash
docker run -p 8080:80 -p 9090:9090 romnn/ldap-manager --generate
```

For a list of options, run with `--help`. If you want to deploy OpenLDAP with LDAP Manager, read along.

### Deployment

##### docker-compose

```bash
docker-compose -f deployment/docker-compose.yml up
```

##### k8s via helm

TODO

##### Considerations

- Serving the frontend externally
    If you have a cluster environment and want to scale the `ldap-manager` container individually or use a more performant static content server like `nginx`, you can disable serving static content using the `--no-static` (`NO_STATIC`) flag.

### Development

#####  Prerequisites

Before you get started, make sure you have installed the following tools::

    $ python3 -m pip install -U cookiecutter>=1.4.0
    $ python3 -m pip install pre-commit bump2version invoke ruamel.yaml halo
    $ go get -u golang.org/x/tools/cmd/goimports
    $ go get -u golang.org/x/lint/golint
    $ go get -u github.com/fzipp/gocyclo
    $ go get -u github.com/mitchellh/gox  # if you want to test building on different architectures

**Remember**: To be able to excecute the tools downloaded with `go get`, 
make sure to include `$GOPATH/bin` in your `$PATH`.
If `echo $GOPATH` does not give you a path make sure to run
(`export GOPATH="$HOME/go"` to set it). In order for your changes to persist, 
do not forget to add these to your shells `.bashrc`.

With the tools in place, it is strongly advised to install the git commit hooks to make sure checks are passing in CI:
```bash
invoke install-hooks
```

You can check if all checks pass at any time:
```bash
invoke pre-commit
```

Note for Maintainers: After merging changes, tag your commits with a new version and push to GitHub to create a release:
```bash
bump2version (major | minor | patch)
git push --follow-tags
```

If you want to (re-)generate the grpc service and gateway source files, make sure to install `protoc`, `protoc-gen-go` and `protoc-gen-go-grpc`.
You can then use the provided script:
```bash
apt install -y protobuf-compiler
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
invoke compile-proto
```

##### Deployment for development

```bash
docker-compose -f dev/docker-compose.yml up --build --force-recreate
```

To quickly work around CORS during development, you could use [proxybootstrap](https://github.com/romnn/proxybootstrap):
```bash
pip install proxybootstrap
proxybootstrap --port 5000 /api@http://127.0.0.1:8090 /@http://127.0.0.1:8080
```

In this example, 8090 is the HTTP service and 8080 is the frontend served via npm.
You can then access the website at [localhost:5000](http://localhost:5000).

#### Note

This project is still in the alpha stage and should not be considered production ready.

#### TODO

- v2
    - documentation
    - add images to the readme
    - Fix flaky tests using fuzzy testing and check slappasswd source
    - Implement missing password hashing algorithms
    - Embed crypt(3) as vendored?

- nice to have
    - Implement CLI interface
        - new acc
        - change password
        - add group
        - add member to group
        - list users
        - verify?
    - decide on a consistent naming (user vs account)
