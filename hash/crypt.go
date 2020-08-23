package hash

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

const (
	// MagicPrefix ...
	MagicPrefix = "$apr1$"
	// RandomSalt ...
	RandomSalt = ""
	// SaltLenMax ...
	SaltLenMax = 8
	// SaltLenMin ...
	SaltLenMin = 1 // Real minimum is 0, but that isn't useful.

)

// Hash64Chars is the character set used by the Hash64 encoding algorithm.
const Hash64Chars = "./0123456789" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz"

// Hash64 is a variant of Base64 encoding.  It is commonly used with
// password hashing algorithms to encode the result of their checksum
// output.
//
// The algorithm operates on up to 3 bytes at a time, encoding the
// following 6-bit sequences into up to 4 hash64 ASCII bytes.
//
//     1. Bottom 6 bits of the first byte
//     2. Top 2 bits of the first byte, and bottom 4 bits of the second byte.
//     3. Top 4 bits of the second byte, and bottom 2 bits of the third byte.
//     4. Top 6 bits of the third byte.
//
// This encoding method does not emit padding bytes as Base64 does.
func Hash64(src []byte) (hash []byte) {
	if len(src) == 0 {
		return []byte{}
	}

	hashSize := (len(src) * 8) / 6
	if (len(src) % 6) != 0 {
		hashSize++
	}
	hash = make([]byte, hashSize)

	dst := hash
	for len(src) > 0 {
		switch len(src) {
		default:
			dst[0] = Hash64Chars[src[0]&0x3f]
			dst[1] = Hash64Chars[((src[0]>>6)|(src[1]<<2))&0x3f]
			dst[2] = Hash64Chars[((src[1]>>4)|(src[2]<<4))&0x3f]
			dst[3] = Hash64Chars[(src[2]>>2)&0x3f]
			src = src[3:]
			dst = dst[4:]
		case 2:
			dst[0] = Hash64Chars[src[0]&0x3f]
			dst[1] = Hash64Chars[((src[0]>>6)|(src[1]<<2))&0x3f]
			dst[2] = Hash64Chars[(src[1]>>4)&0x3f]
			src = src[2:]
			dst = dst[3:]
		case 1:
			dst[0] = Hash64Chars[src[0]&0x3f]
			dst[1] = Hash64Chars[(src[0]>>6)&0x3f]
			src = src[1:]
			dst = dst[2:]

		}
	}

	return
}

// GenerateSalt generates a random salt parameter string of a given length.
//
// If the length is greater than SaltLenMax, a string of that length
// will be returned instead. Similarly, if length is less than
// SaltLenMin, a string of that length will be returned instead.
func GenerateSalt(length int) string {
	if length > SaltLenMax {
		length = SaltLenMax
	} else if length < SaltLenMin {
		length = SaltLenMin
	}
	rlen := (length * 6 / 8)
	if (length*6)%8 != 0 {
		rlen++
	}
	buf := make([]byte, rlen)
	rand.Read(buf)
	salt := Hash64(buf)
	return fmt.Sprintf("%s%s", MagicPrefix, salt[:length])
}

// Crypt takes key and salt strings and performs the MD5-Crypt hashing
// algorithm on them, returning a full hash string suitable for storage
// and later password verification.
//
// If the salt string is the value RandomSalt, a randomly-generated salt
// parameter string will be generated of length SaltLenMax.
func Crypt(keystr, saltstr string) string {
	var key, salt []byte
	var keyLen, saltLen int

	key = []byte(keystr)
	keyLen = len(key)

	if saltstr == "" {
		saltstr = GenerateSalt(SaltLenMax)
	}
	saltbytes := []byte(saltstr)
	if !bytes.HasPrefix(saltbytes, []byte(MagicPrefix)) {
		return "invalid magic prefix"
	}

	salttoks := bytes.Split(saltbytes, []byte{'$'})
	numtoks := len(salttoks)

	if numtoks < 3 {
		return "invalid salt format"
	}
	salt = salttoks[2]

	if len(salt) > 8 {
		salt = salt[0:8]
	}
	saltLen = len(salt)

	B := md5.New()
	B.Write(key)
	B.Write(salt)
	B.Write(key)
	Bsum := B.Sum(nil)

	A := md5.New()
	A.Write(key)
	A.Write([]byte(MagicPrefix))
	A.Write(salt)
	cnt := keyLen
	for ; cnt > 16; cnt -= 16 {
		A.Write(Bsum)
	}
	A.Write(Bsum[0:cnt])
	for cnt = keyLen; cnt > 0; cnt >>= 1 {
		if (cnt & 1) == 0 {
			A.Write(key[0:1])
		} else {
			A.Write([]byte{0})
		}
	}
	Asum := A.Sum(nil)

	Csum := Asum
	for round := 0; round < 1000; round++ {
		C := md5.New()

		if (round & 1) != 0 {
			C.Write(key)
		} else {
			C.Write(Csum)
		}

		if (round % 3) != 0 {
			C.Write(salt)
		}

		if (round % 7) != 0 {
			C.Write(key)
		}

		if (round & 1) == 0 {
			C.Write(key)
		} else {
			C.Write(Csum)
		}

		Csum = C.Sum(nil)
	}

	buf := make([]byte, 0, 23+len(MagicPrefix)+saltLen)
	buf = append(buf, MagicPrefix...)
	buf = append(buf, salt...)
	buf = append(buf, '$')
	buf = append(buf, Hash64([]byte{
		Csum[12], Csum[6], Csum[0],
		Csum[13], Csum[7], Csum[1],
		Csum[14], Csum[8], Csum[2],
		Csum[15], Csum[9], Csum[3],
		Csum[5], Csum[10], Csum[4],
		Csum[11],
	})...)

	return string(buf)
}

// Verify hashes a key using the same salt parameters as the given
// hash string, and if the results match, it returns true.
func Verify(key, hash string) bool {
	nhash := Crypt(key, hash)
	if hash == nhash {
		return true
	}
	return false
}

/*
import "encoding/binary"

func Hash64(buffer []byte, seed uint64) uint64 {
	const (
		k0 = 0xD6D018F5
		k1 = 0xA2AA033B
		k2 = 0x62992FC1
		k3 = 0x30BC5B29
	)

	ptr := buffer

	hash := (seed + k2) * k0

	if len(ptr) >= 32 {
		v := [4]uint64{hash, hash, hash, hash}

		for len(ptr) >= 32 {
			v[0] += binary.LittleEndian.Uint64(ptr) * k0
			ptr = ptr[8:]
			v[0] = rotateRight(v[0], 29) + v[2]
			v[1] += binary.LittleEndian.Uint64(ptr) * k1 // Check Needed because of line 27
			ptr = ptr[8:]
			v[1] = rotateRight(v[1], 29) + v[3]
			v[2] += binary.LittleEndian.Uint64(ptr) * k2 // Check Needed because of line 30
			ptr = ptr[8:]
			v[2] = rotateRight(v[2], 29) + v[0]
			v[3] += binary.LittleEndian.Uint64(ptr) * k3 // Check Needed because of line 33
			ptr = ptr[8:]
			v[3] = rotateRight(v[3], 29) + v[1]
		}
	}

	// Extra Code Excluded

	return hash
}

func rotateRight(v uint64, k uint) uint64 {
	return (v >> k) | (v << (64 - k))
}
*/
