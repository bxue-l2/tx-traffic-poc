package main

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = "traffic-generator"
	envPrefix  = "TRAFFIC_GENERATOR"
)

var (
	/* Required Flags */

	HostnameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-hostname"),
		Usage:    "Hostname at which disperser service is available",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "HOSTNAME"),
	}
	TimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "timeout"),
		Usage:    "Amount of time to wait for GPRC",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "TIMEOUT"),
		Value:    10 * time.Second,
	}
	NumInstancesFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-instances"),
		Usage:    "Number of generator instances to run in parallel",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_INSTANCES"),
	}
	PadSizeFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "data-size"),
		Usage:    "Data size in the extra field",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "PAD_SIZE"),
	}
	RequestIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "request-interval"),
		Usage:    "Duration between requests",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "REQUEST_INTERVAL"),
		Value:    2 * time.Second,
	}
	SignerPrivateKeysFlag = cli.StringSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signer-private-keys-hex"),
		Usage:    "List of Private key to use for signing requests",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "SIGNER_PRIVATE_KEYS_HEX"),
	}
	AddressesFlag = cli.StringSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signer-addresses-hex"),
		Usage:    "List of Address to use for signing requests",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "SIGNER_ADDRESSES"),
	}
)

var requiredFlags = []cli.Flag{
	NumInstancesFlag,
	RequestIntervalFlag,
	SignerPrivateKeysFlag,
	AddressesFlag,
	PadSizeFlag,
}

var optionalFlags = []cli.Flag{
	HostnameFlag,
	TimeoutFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
