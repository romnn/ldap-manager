#!/bin/sh
ldapsearch -LLL -o ldif-wrap=no -h localhost -p 389 \
    -b "ou=groups,dc=example,dc=org" \
    -D "cn=admin,dc=example,dc=org" \
    -w "admin" \
    '(cn=admins)' dn