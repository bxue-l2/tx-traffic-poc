package main

import (
	"time"

	"github.com/urfave/cli"
)

type Config struct {
	NumInstances      uint
	RequestInterval   time.Duration
	SignerPrivateKeys []string
	Addresses         []string
	Hostname          string
	Timeout           time.Duration
	PadSize           uint
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	return &Config{
		NumInstances:      ctx.GlobalUint(NumInstancesFlag.Name),
		RequestInterval:   ctx.Duration(RequestIntervalFlag.Name),
		SignerPrivateKeys: ctx.StringSlice(SignerPrivateKeysFlag.Name),
		Addresses:         ctx.StringSlice(AddressesFlag.Name),
		Timeout:           ctx.Duration(TimeoutFlag.Name),
		Hostname:          ctx.String(HostnameFlag.Name),
		PadSize:           ctx.GlobalUint(PadSizeFlag.Name),
	}, nil
}
