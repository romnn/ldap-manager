package cli

import (
	"github.com/urfave/cli/v2"
	"time"
)

var (
	Key = cli.StringFlag{
		Name:    "key",
		Aliases: []string{"public-key", "signing-key"},
		EnvVars: []string{"PRIVATE_KEY", "KEY", "SIGNING_KEY"},
		Usage:   "private key to sign the tokens with",
	}
	Jwks = cli.StringFlag{
		Name:    "jwks",
		Aliases: []string{"jwks-json", "jwk-set"},
		EnvVars: []string{"JWKS", "JWK_SET", "JWKS_JSON"},
		Usage:   "json encoded jwk set containing the public keys",
	}
	KeyFile = cli.PathFlag{
		Name:    "key-file",
		Aliases: []string{"public-key-file", "signing-key-file"},
		EnvVars: []string{"PRIVATE_KEY_FILE", "KEY_FILE", "SIGNING_KEY_FILE"},
		Usage:   "file with private key to sign the tokens with",
	}
	JwksFile = cli.PathFlag{
		Name:    "jwks-file",
		Aliases: []string{"jwks-json-file", "jwk-set-file"},
		EnvVars: []string{"JWKS_FILE", "JWK_SET_FILE", "JWKS_JSON_FILE"},
		Usage:   "json file with the jwk set containing the public keys",
	}
	Generate = cli.BoolFlag{
		Name:    "generate",
		Value:   false,
		Aliases: []string{"gen", "create"},
		EnvVars: []string{"GENERATE", "CREATE", "GEN"},
		Usage:   "generate new keys if none were supplied",
	}
	ExpirationTime = cli.GenericFlag{
		Name: "expiration-time",
		Value: &DurationValue{
			Default: 24 * time.Hour,
		},
		Aliases: []string{"token-expire"},
		EnvVars: []string{"EXPIRATION_TIME", "EXPIRATION_TIME"},
		Usage:   "expiration times for JWT tokens (e.g. 5h30m40s)",
	}
	Issuer = cli.StringFlag{
		Name:    "issuer",
		Value:   "issuer@example.org",
		Aliases: []string{"jwt-issuer"},
		EnvVars: []string{"ISSUER"},
		Usage:   "jwt token issuer",
	}
	Audience = cli.StringFlag{
		Name:    "audience",
		Value:   "example.org",
		Aliases: []string{"jwt-audience"},
		EnvVars: []string{"AUDIENCE"},
		Usage:   "jwt token audience",
	}
	AuthFlags = []cli.Flag{
		&Key,
		&Jwks,
		&KeyFile,
		&JwksFile,
		&Generate,
		&ExpirationTime,
		&Issuer,
		&Audience,
	}
)
