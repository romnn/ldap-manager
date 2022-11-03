package cli

import (
	"github.com/urfave/cli/v2"
)

var (
	// LogLevel configures the logging level of the service
	LogLevel = cli.GenericFlag{
		Name: "log",
		Value: &EnumValue{
			Enum: []string{
				"info",
				"debug",
				"warn",
				"fatal",
				"trace",
				"error",
				"panic",
			},
			Default: "info",
		},
		Aliases: []string{"log-level"},
		EnvVars: []string{"LOG", "LOG_LEVEL"},
		Usage:   "Log level",
	}
	// GRPCPort configures the port to serve GRPC
	GRPCPort = cli.IntFlag{
		Name:    "grpc-port",
		Value:   9090,
		EnvVars: []string{"GRPC_PORT"},
		Usage:   "GRPC service port",
	}
	// HTTPPort configures the port to serve HTTP
	HTTPPort = cli.IntFlag{
		Name:    "http-port",
		Value:   8080,
		Aliases: []string{"port"},
		EnvVars: []string{"HTTP_PORT", "PORT"},
		Usage:   "HTTP service port",
	}
	// NoStatic configures if static assets should not be served
	NoStatic = cli.BoolFlag{
		Name:    "no-static",
		Value:   false,
		Aliases: []string{"disable-serve-static"},
		EnvVars: []string{"NO_STATIC", "DISABLE_SERVE_STATIC"},
		Usage:   "disable serving of the static frontend",
	}
	// StaticRoot configures the static file root dir
	StaticRoot = cli.StringFlag{
		Name:    "static-root",
		Value:   "./web/dist",
		EnvVars: []string{"STATIC_DIR", "STATIC_ROOT"},
		Usage:   "root source directory of the static files to be served",
	}
	// ServiceFlags is the set of all service CLI flags
	ServiceFlags = []cli.Flag{
		&LogLevel,
		&GRPCPort,
		&HTTPPort,
		&NoStatic,
		&StaticRoot,
	}
)
