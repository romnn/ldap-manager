package hash

import (
	"fmt"
	// "encoding/base64"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"math/big"
)

// LDAPPasswordHashingAlgorithm ...
type LDAPPasswordHashingAlgorithm uint16

const (
	// Default ...
	Default LDAPPasswordHashingAlgorithm = iota
	// SHA512CRYPT ...
	SHA512CRYPT
	// SHA256CRYPT ...
	SHA256CRYPT
	// BLOWFISH ...
	BLOWFISH
	// EXTDES ...
	EXTDES
	// MD5CRYPT ...
	MD5CRYPT
	// SMD5 ...
	SMD5
	// MD5 ...
	MD5
	// SHA ...
	SHA
	// SSHA ...
	SSHA
	// CRYPT ...
	CRYPT
	// CLEAR ...
	CLEAR
)

// LDAPPasswordHashingAlgorithms ...
var LDAPPasswordHashingAlgorithms = map[string]LDAPPasswordHashingAlgorithm{
	"Default":     Default,
	"SHA512CRYPT": SHA512CRYPT,
	"SHA256CRYPT": SHA256CRYPT,
	"BLOWFISH":    BLOWFISH,
	"EXTDES":      EXTDES,
	"MD5CRYPT":    MD5CRYPT,
	"SMD5":        SMD5,
	"MD5":         MD5,
	"SHA":         SHA,
	"SSHA":        SSHA,
	"CRYPT":       CRYPT,
	"CLEAR":       CLEAR,
}

func generateSaltDeprecated(l int) (string, error) {
	permitted := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ./")
	var salt []rune
	var i int
	for {
		if len(salt) >= l {
			break
		}
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(permitted))))
		if err != nil {
			return "", err
		}
		salt[i] = permitted[int(n.Int64())]
		i++
	}
	return string(salt), nil
}

func encodeSSHA(pw string) string {
	return ""
}

func encodeSHA256(pw string) string {
	// '{CRYPT}' . crypt($password, '$5$' . generate_salt(8))
	h := sha256.New()
	h.Write([]byte(pw))
	h.Write([]byte("$5$"))
	h.Write(generateSalt(8))
	return fmt.Sprintf("{CRYPT}%s", h.Sum(nil))
}

func encodeSHA512(pw string) string {
	// '{CRYPT}' . crypt($password, '$6$' . generate_salt(8));
	salt := []byte("12345678") // generateSalt(8)
	fmt.Println(salt)
	fmt.Println(pw)
	h := sha512.New()
	h.Write([]byte(pw))
	h.Write([]byte("$6$"))
	h.Write(salt)
	return fmt.Sprintf("{CRYPT}%s", h.Sum(nil))
}

/*
func encodeMD5(pw string) string {
	// '{MD5}' . base64_encode(md5($password, TRUE));
	h := sha512.New()
	h.Write([]byte(pw))
	h.Write([]byte("$6$"))
	h.Write(generateSalt(8))
	b64 := base64.StdEncoding.EncodeToString(hash)
	return fmt.Sprintf("{MD5}%s", b64), nil
}
*/

// generateSalt generates a byte array containing random bytes
func generateSalt(l int) []byte {
	sbytes := make([]byte, l)
	rand.Read(sbytes)
	return sbytes
}

// Password ...
func Password(password string, algorithm LDAPPasswordHashingAlgorithm) (string, error) {
	switch algorithm {
	case SSHA, Default:
		return encodeSSHA(password), nil
	case SHA256CRYPT:
		return encodeSHA256(password), nil
	case SHA512CRYPT:
		return encodeSHA512(password), nil
	default:
		return encodeSSHA(password), nil
	}
}

/*
$hash_algo = 'SSHA';

	switch ($hash_algo) {

   # Blowfish & EXT_DES didn't work
   #  case 'BLOWFISH':
   #    $hashed_pwd = '{CRYPT}' . crypt($password, '$2a$12$' . generate_salt(13));
   #    break;

   #  case 'EXT_DES':
   #    $hashed_pwd = '{CRYPT}' . crypt($password, '_' . generate_salt(8));
   #    break;

	 case 'MD5CRYPT':
	   $hashed_pwd = '{CRYPT}' . crypt($password, '$1$' . generate_salt(9));
	   break;

	 case 'SMD5':
	   $salt = generate_salt(8);
	   $hashed_pwd = '{SMD5}' . base64_encode(md5($password . $salt, TRUE) . $salt);
	   break;

	 case 'SHA':
	   $hashed_pwd = '{SHA}' . base64_encode(sha1($password, TRUE));
	   break;

	 case 'SSHA':
	   $salt = generate_salt(8);
	   $hashed_pwd = '{SSHA}' . base64_encode(sha1($password . $salt, TRUE) . $salt);
	   break;

	 case 'CRYPT':
	   $salt = generate_salt(2);
	   $hashed_pwd = '{CRYPT}' . crypt($password, $salt);
	   break;

	 case 'CLEAR':
	   error_log("$log_prefix password hashing - WARNING - Saving password in cleartext. This is extremely bad practice and should never ever be done in a production environment.");
	   $hashed_pwd = $password;
	   break;


	}
*/

/*
function ldap_hashed_password($password) {

	global $PASSWORD_HASH, $log_prefix;

	$check_algos = array (
						  "SHA512CRYPT" => "CRYPT_SHA512",
						  "SHA256CRYPT" => "CRYPT_SHA256",
   #                       "BLOWFISH"    => "CRYPT_BLOWFISH",
   #                       "EXT_DES"     => "CRYPT_EXT_DES",
						  "MD5CRYPT"    => "CRYPT_MD5"
						 );

	$remaining_algos = array (
							   "SSHA",
							   "SHA",
							   "SMD5",
							   "MD5",
							   "CRYPT",
							   "CLEAR"
							 );

	$available_algos = array();

	foreach ($check_algos as $algo_name => $algo_function) {
	  if (defined($algo_function) and constant($algo_function) != 0) {
		array_push($available_algos, $algo_name);
	  }
	  else {
		error_log("$log_prefix password hashing - the system doesn't support ${algo_name}");
	  }
	}
	$available_algos = array_merge($available_algos, $remaining_algos);

	if (isset($PASSWORD_HASH)) {
	  if (!in_array($PASSWORD_HASH, $available_algos)) {
		$hash_algo = $available_algos[0];
		error_log("$log_prefix LDAP password: the chosen hash method ($PASSWORD_HASH) wasn't available");
	  }
	  else {
		$hash_algo = $PASSWORD_HASH;
	  }
	}
	else {
	  $hash_algo = $available_algos[0];
	}
	error_log("$log_prefix LDAP password: using '${hash_algo}' as the hashing method");

	$hash_algo = 'SSHA';

	switch ($hash_algo) {

	 case 'SHA512CRYPT':
	   $hashed_pwd = '{CRYPT}' . crypt($password, '$6$' . generate_salt(8));
	   break;

	 case 'SHA256CRYPT':
	   $hashed_pwd = '{CRYPT}' . crypt($password, '$5$' . generate_salt(8));
	   break;

   # Blowfish & EXT_DES didn't work
   #  case 'BLOWFISH':
   #    $hashed_pwd = '{CRYPT}' . crypt($password, '$2a$12$' . generate_salt(13));
   #    break;

   #  case 'EXT_DES':
   #    $hashed_pwd = '{CRYPT}' . crypt($password, '_' . generate_salt(8));
   #    break;

	 case 'MD5CRYPT':
	   $hashed_pwd = '{CRYPT}' . crypt($password, '$1$' . generate_salt(9));
	   break;

	 case 'SMD5':
	   $salt = generate_salt(8);
	   $hashed_pwd = '{SMD5}' . base64_encode(md5($password . $salt, TRUE) . $salt);
	   break;

	 case 'MD5':
	   $hashed_pwd = '{MD5}' . base64_encode(md5($password, TRUE));
	   break;

	 case 'SHA':
	   $hashed_pwd = '{SHA}' . base64_encode(sha1($password, TRUE));
	   break;

	 case 'SSHA':
	   $salt = generate_salt(8);
	   $hashed_pwd = '{SSHA}' . base64_encode(sha1($password . $salt, TRUE) . $salt);
	   break;

	 case 'CRYPT':
	   $salt = generate_salt(2);
	   $hashed_pwd = '{CRYPT}' . crypt($password, $salt);
	   break;

	 case 'CLEAR':
	   error_log("$log_prefix password hashing - WARNING - Saving password in cleartext. This is extremely bad practice and should never ever be done in a production environment.");
	   $hashed_pwd = $password;
	   break;


	}

	error_log("$log_prefix password update - algo $hash_algo | pwd $hashed_pwd");

	return $hashed_pwd;

   }
*/
