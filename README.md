## ldap-manager

[![Build Status](https://travis-ci.com/romnnn/ldap-manager.svg?branch=master)](https://travis-ci.com/romnnn/ldap-manager)
[![GitHub](https://img.shields.io/github/license/romnnn/ldap-manager)](https://github.com/romnnn/ldap-manager)
 [![Docker Pulls](https://img.shields.io/docker/pulls/romnn/ldap-manager)](https://hub.docker.com/r/romnn/ldap-manager) [![Test Coverage](https://codecov.io/gh/romnnn/ldap-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/romnnn/ldap-manager)
[![Release](https://img.shields.io/github/release/romnnn/ldap-manager)](https://github.com/romnnn/ldap-manager/releases/latest)

Your description goes here...


```bash
go get github.com/romnnn/ldap-manager
```


You can also download pre built binaries from the [releases page](https://github.com/romnnn/ldap-manager/releases), or use the `docker` image:

```bash
docker pull romnn/ldap-manager
```

For a list of options, run with `--help`.




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
docker run --entrypoint slappasswd  mlan/openldap -s <my-password>
```

#### Development deployment

```bash
docker-compose -f dev/docker-compose.yml up --build --force-recreate
```

#### Note

This project is still in the alpha stage and should not be considered production ready.

#### TODO

- Resilient creation of groups (setup)
- Serve static content
- Update dockerfile
- Implement frontend
- Write tests using testcontainers
- Restructure to allow usage via GRPC API?
- Restructure to allow CLI usage