## ldap-manager

[![Build Status](https://github.com/romnn/ldap-manager/workflows/test/badge.svg)](https://github.com/romnn/ldap-manager/actions)
[![GitHub](https://img.shields.io/github/license/romnn/ldap-manager)](https://github.com/romnn/ldap-manager)
[![Docker Pulls](https://img.shields.io/docker/pulls/romnn/ldap-manager)](https://hub.docker.com/r/romnn/ldap-manager)
[![Test Coverage](https://codecov.io/gh/romnn/ldap-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/romnn/ldap-manager)
[![Release](https://img.shields.io/github/release/romnn/ldap-manager)](https://github.com/romnn/ldap-manager/releases/latest)

<p align="center">
  <img width="200" src="website/icon/icon_lg.jpg">
</p>

LDAP Manager is the cloud-native LDAP web management interface. LDAP has been around for a long time and has become a popular choice for user and group management - however, this should not mean that it's management interface should be hard to deploy and look and feel like it was made in the last century.

LDAP Manager is written in Go and comes with a Vue/Typescript frontend in a single, self-contained docker container. It also exposes it's API over both REST and gRPC!

|                                                 |                                                 |
| :---------------------------------------------: | :---------------------------------------------: |
|      <img src="screenshots/home-of-user.png">      | <img src="screenshots/user-edit-by-admin.png"> |
| <img src="screenshots/users-list-for-admin.png"> |  <img src="screenshots/group-edit-by-admin.png">  |

Before you get started, make sure you have an OpenLDAP server like
[osixia/openldap](https://hub.docker.com/r/osixia/openldap/) running.
For more information on deployment and a full example,
see the [deployment guide](#Deployment).

```bash
go install github.com/romnn/ldap-manager/cmd/ldap-manager
ldap-manager serve --generate
go run github.com/romnn/ldap-manager/cmd/ldap-manager serve --generate --http-port 8090
```

You can also download pre-built binaries from the
[releases page](https://github.com/romnn/ldap-manager/releases),
or use the `docker` image:

```bash
docker run -p 8080:80 -p 9090:9090 romnn/ldap-manager --generate
```

For a list of options, run with `--help`. If you want to deploy OpenLDAP with LDAP Manager, read along.

### Deployment

```bash
helm dependency update deployment/helm/charts/ldapmanager/
```

##### docker-compose

```bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f deployment/docker-compose.yml up
```

##### Helm

TODO

##### Known bugs


##### Fixed bugs
- when creating new groups, errors are not properly passed through
- when creating new group, cannot add initial members
- removing groups in edit user does not do anything (no errors)
- errors when removing from group in user edit page break layout
- memberof does not respect the initial admin user for both users and admins groups

##### Considerations

- Serving the frontend externally
  If you have a cluster environment and want to scale the `ldap-manager` container individually or use a more performant static content server like `nginx`, you can disable serving static content using the `--no-static` (`NO_STATIC`) flag.

### Development

#### Tools

Before you get started, make sure you have installed the following tools:

```bash
$ python3 -m pip install pre-commit bump2version invoke
$ go install github.com/kyoh86/richgo@latest
$ go install golang.org/x/tools/cmd/goimports@latest
$ go install golang.org/x/lint/golint@latest
$ go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
```

Please always make sure code checks pass:

```bash
inv pre-commit
```

#### Compiling proto sources

If you want to (re-)compile the grpc service and gateway `.proto` source files,
you will need

- `protoc`
- `protoc-gen-go`
- `protoc-gen-go-grpc`.
- `protoc-gen-grpc-gateway`
- `protoc-gen-openapiv2`

```bash
apt install -y protobuf-compiler
brew install protobuf

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

To compile the protos, you can use the provided script:

```bash
inv compile-proto
```

#### Generate screenshots
```bash
cd deployment/screenshot
yarn install --dev
yarn run screenshot
```

#### TODO

- v2

  - refactor to use manual ldap search only where necessary
  - use an interface for the main functions of the manager in GRPC server
  - point out that the goal is user management only

  - documentation

- nice to have
  - test the grpc and http servers as well

  - Implement CLI interface
    - new acc
    - change password
    - add group
    - add member to group
    - list users
    - verify?

- done
  - refactor in general
  - add integration test with harbor
  - fix the docker container
  - fix the frontend
  - binds should always open new connections
  - function that opens a new connection
  - add a simple connection pool
  - fix nil pointer errors
  - decide what goes into pkg and what goes into internal
  - add tests for each file in pkg
  - Implement missing password hashing algorithms
  - Embed crypt(3) as vendored?
  - Fix flaky tests using fuzzy testing and check slappasswd source
  - add pagination
  - get rid of the password hashing mess
  - decide on a consistent naming (user vs account)
  - split into more files
  - update dependencies
  - fix issues and use new api for grpc and http without a base
  - add images to the readme
