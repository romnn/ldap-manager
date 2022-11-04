import hashlib
import base64

"""
$ slappasswd
New password: tony
Re-enter new password: tony
{SSHA}rj5qElfD1JK7PCKSLZzxIldkREeFW2Dv
$ echo "rj5qElfD1JK7PCKSLZzxIldkREeFW2Dv" | base64 -D | hexdump -C
00000000  ae 3e 6a 12 57 c3 d4 92  bb 3c 22 92 2d 9c f1 22  |.>j.W....<".-.."|
00000010  57 64 44 47 85 5b 60 ef                           |WdDG.[|
00000018
"""

# {SSHA}jmzeBy7KhOsnV2dBOt6D7jUQZIRc7wz/
pw = "{SSHA}rj5qElfD1JK7PCKSLZzxIldkREeFW2Dv"  # "tony"

# base64 decode the LDAP's userPassword
pw_enc = pw[6:] # .encode("utf-8")
print(pw_enc)
pw_dec = base64.b64decode(pw_enc)
print(pw_dec)
# .decode("utf-8"))

# because this is SSHA: the last 4 bytes is the salt!
salt = pw_dec[-4:]
print(salt)

clear = "config"

# sprinkle the salt..
challenge_salted = clear + salt

# hash it ...
challenge_sha1 = hashlib.sha1(challenge_salted).digest()

# append with the original salt, and base64 encode it!
challenge_complete = "{SSHA}" + base64.encodestring(challenge_sha1 + salt_str)

print("%s vs %s" % (pw, challenge_complete))
