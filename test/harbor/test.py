import ldap

connect = ldap.initialize("ldap://localhost:49964")
connect.set_option(ldap.OPT_REFERRALS, 0)
base_dn = "dc=example,dc=org"
test = connect.simple_bind_s(f"cn=readonly,{base_dn}", "readonly")

print(test)
result = connect.search_s(
    base_dn,
    ldap.SCOPE_SUBTREE,
    "*",
    ["cn"],
)
print(result)
