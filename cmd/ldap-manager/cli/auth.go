package cli

import (
	"time"

	"github.com/urfave/cli/v2"
)

var (
	// Key to sign the tokens with
	Key = cli.StringFlag{
		Name:    "key",
		Aliases: []string{"public-key", "signing-key"},
		EnvVars: []string{"PRIVATE_KEY", "KEY", "SIGNING_KEY"},
		Usage:   "private key to sign the tokens with",
	}
	// Jwks set containing the public keys
	Jwks = cli.StringFlag{
		Name:    "jwks",
		Aliases: []string{"jwks-json", "jwk-set"},
		EnvVars: []string{"JWKS", "JWK_SET", "JWKS_JSON"},
		Usage:   "json encoded jwk set containing the public keys",
	}
	// KeyFile with private key to sign the tokens with
	KeyFile = cli.PathFlag{
		Name:    "key-file",
		Aliases: []string{"public-key-file", "signing-key-file"},
		EnvVars: []string{"PRIVATE_KEY_FILE", "KEY_FILE", "SIGNING_KEY_FILE"},
		Usage:   "file with private key to sign the tokens with",
	}
	// JwksFile with the jwk set containing the public keys
	JwksFile = cli.PathFlag{
		Name:    "jwks-file",
		Aliases: []string{"jwks-json-file", "jwk-set-file"},
		EnvVars: []string{"JWKS_FILE", "JWK_SET_FILE", "JWKS_JSON_FILE"},
		Usage:   "json file with the jwk set containing the public keys",
	}
	// Generate configures if keys should be generated if not supplied
	Generate = cli.BoolFlag{
		Name:    "generate",
		Value:   false,
		Aliases: []string{"gen", "create"},
		EnvVars: []string{"GENERATE", "CREATE", "GEN"},
		Usage:   "generate new keys if none were supplied",
	}
	// ExpirationTime of JWT tokens
	ExpirationTime = cli.GenericFlag{
		Name: "expiration-time",
		Value: &DurationValue{
			Default: 24 * time.Hour,
		},
		Aliases: []string{"token-expire"},
		EnvVars: []string{"EXPIRATION_TIME", "EXPIRATION_TIME"},
		Usage:   "expiration times for JWT tokens (e.g. 5h30m40s). Default is 24h",
	}
	// Issuer for the JWT tokens
	Issuer = cli.StringFlag{
		Name:    "issuer",
		Value:   "issuer@example.org",
		Aliases: []string{"jwt-issuer"},
		EnvVars: []string{"ISSUER"},
		Usage:   "JWT token issuer",
	}
	// Audience for the JWT tokens
	Audience = cli.StringFlag{
		Name:    "audience",
		Value:   "example.org",
		Aliases: []string{"jwt-audience"},
		EnvVars: []string{"AUDIENCE"},
		Usage:   "JWT token audience",
	}
	// AuthFlags is a set of all CLI flags
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
