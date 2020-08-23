package hash

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

// based on: https://medium.com/@ferdinand.neman/ssha-password-hash-with-golang-7d79d792bd3d

// SSHAEncoder ...
type SSHAEncoder struct{}

// Encode encodes the []byte of raw password
func (enc SSHAEncoder) Encode(rawPassPhrase []byte) ([]byte, error) {
	hash := makeSSHAHash(rawPassPhrase, generateSalt(4))
	b64 := base64.StdEncoding.EncodeToString(hash)
	return []byte(fmt.Sprintf("{SSHA}%s", b64)), nil
}

// Matches matches the encoded password and the raw password
func (enc SSHAEncoder) Matches(encodedPassPhrase, rawPassPhrase []byte) bool {
	// strip the {SSHA} prefix
	eppS := string(encodedPassPhrase)[6:]
	hash, err := base64.StdEncoding.DecodeString(eppS)
	if err != nil {
		return false
	}
	salt := hash[len(hash)-4:]

	sha := sha1.New()
	sha.Write(rawPassPhrase)
	sha.Write(salt)
	sum := sha.Sum(nil)

	if bytes.Compare(sum, hash[:len(hash)-4]) != 0 {
		return false
	}
	return true
}

// makeSSHAHash make hasing using SHA-1 with salt.
// This is not the final output though. You need to append {SSHA} string with base64 of this hash.
func makeSSHAHash(passphrase, salt []byte) []byte {
	sha := sha1.New()
	sha.Write(passphrase)
	sha.Write(salt)

	h := sha.Sum(nil)
	return append(h, salt...)
}
