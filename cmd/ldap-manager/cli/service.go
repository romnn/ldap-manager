package cli

import (
	"github.com/urfave/cli/v2"
)

var (
	LogLevel = cli.GenericFlag{
		Name: "log",
		Value: &EnumValue{
			Enum:    []string{"info", "debug", "warn", "fatal", "trace", "error", "panic"},
			Default: "info",
		},
		Aliases: []string{"log-level"},
		EnvVars: []string{"LOG", "LOG_LEVEL"},
		Usage:   "Log level",
	}
	GrpcPort = cli.IntFlag{
		Name:    "grpc-port",
		Value:   9090,
		EnvVars: []string{"GRPC_PORT"},
		Usage:   "grpc service port",
	}
	HttpPort = cli.IntFlag{
		Name:    "http-port",
		Value:   80,
		Aliases: []string{"port"},
		EnvVars: []string{"HTTP_PORT", "PORT"},
		Usage:   "http service port",
	}
	NoStatic = cli.BoolFlag{
		Name:    "no-static",
		Value:   false,
		Aliases: []string{"disable-serve-static"},
		EnvVars: []string{"NO_STATIC", "DISABLE_SERVE_STATIC"},
		Usage:   "disable serving of the static frontend",
	}
	StaticRoot = cli.StringFlag{
		Name:    "static-root",
		Value:   "./frontend/dist",
		EnvVars: []string{"STATIC_DIR", "STATIC_ROOT"},
		Usage:   "root source directory of the static files to be served",
	}
	ServiceFlags = []cli.Flag{
		&LogLevel,
		&GrpcPort,
		&HttpPort,
		&NoStatic,
		&StaticRoot,
	}
)
