package uhppoted

import (
	"flag"
	"fmt"
)

type VersionV struct {
	Application string
	Version     string
}

func (cmd *VersionV) Name() string {
	return "version"
}

func (cmd *VersionV) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("version", flag.ExitOnError)
}

func (cmd *VersionV) Execute(args ...interface{}) error {
	fmt.Printf("%v\n", cmd.Version)

	return nil
}

func (cmd *VersionV) Description() string {
	return "Displays the current version"
}

func (cmd *VersionV) Usage() string {
	return ""
}

func (cmd *VersionV) Help() {
	fmt.Println()
	fmt.Printf("  Displays the %s version in the format v<major>.<minor>.<build> e.g. v1.00.10\n", cmd.Application)
	fmt.Println()
}
