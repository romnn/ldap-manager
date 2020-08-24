## ldap-manager

[![Build Status](https://travis-ci.com/romnnn/ldap-manager.svg?branch=master)](https://travis-ci.com/romnnn/ldap-manager)
[![GitHub](https://img.shields.io/github/license/romnnn/ldap-manager)](https://github.com/romnnn/ldap-manager)
 [![Docker Pulls](https://img.shields.io/docker/pulls/romnn/ldap-manager)](https://hub.docker.com/r/romnn/ldap-manager) [![Test Coverage](https://codecov.io/gh/romnnn/ldap-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/romnnn/ldap-manager)
[![Release](https://img.shields.io/github/release/romnnn/ldap-manager)](https://github.com/romnnn/ldap-manager/releases/latest)

Your description goes here...


```bash
go get github.com/romnnn/ldap-manager

# run the HTTP server
go run github.com/romnnn/ldap-manager/cmd/ldap-manager serve --port 8090 http

# run the gRPC server
go run github.com/romnnn/ldap-manager/cmd/ldap-manager serve --port 8090 grpc
```

You can also download pre built binaries from the [releases page](https://github.com/romnnn/ldap-manager/releases), or use the `docker` image:

```bash
docker pull romnn/ldap-manager
```

For a list of options, run with `--help`.

TODO: Notes on static content and GRPC and CLI with references

#### Deployment (docker-compose)

TODO

#### Deployment (k8s via helm)

TODO

#### Serving the frontend externally

If you have a cluster environment and want to scale the `ldap-manager` container individually and use a more performant static servicer like `nginx`, you can disable serving static content using the `--no-static` (`NO_STATIC`) flag.

TODO: nginx example

#### Development

######  Prerequisites

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

If you want to (re-)generate the sample grpc service, make sure to install `protoc`, `protoc-gen-go` and `protoc-gen-go-grpc`.
You can then use the provided script:
```bash
apt install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
invoke compile-proto
```

#### Debug user config

To manually add a user via an LDIF file, use the `ldapadd` command
```bash
sudo apt install ldap-utils
# For debugging, it is a good idea to add the entries manually one after the other
ldapmodify -h localhost -p 389 -D cn=admin,dc=example,dc=org -w "admin" -f dev/pre-configured-users/1_add_ous.ldif
ldapmodify -h localhost -p 389 -D cn=admin,dc=example,dc=org -w "admin" -f dev/pre-configured-users/2_add_admin_group.ldif
```

#### Generate LDAP passwords

```bash
# This will use the default SSHA (Salted SHA1)
docker run --entrypoint slappasswd  mlan/openldap -s <my-password>
# You can generate SHA512 or others, see UNIX crypt(3) or PHP $crypt() for reference
docker run --entrypoint slappasswd  mlan/openldap -s 123456 -c '$6$%.16s'
```

#### Development deployment

```bash
docker-compose -f dev/docker-compose.yml up --build --force-recreate
```

#### Note

This project is still in the alpha stage and should not be considered production ready.

#### TODO

- Implement frontend
- Implement token based authentication
- Fully configure travis
- Restructure to allow usage via GRPC API?
- Restructure to allow CLI usage
    - new acc
    - change password
    - add group
    - add member to group
    - list users
    - verify?
- Publish helm chart via github pages
- Add images to the readme
- Rename users to accounts
- Implement missing password hashing algorithms
- Embed crypt(3) as vendored?