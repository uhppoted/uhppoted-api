package uhppoted

import (
	"context"
	"flag"
	"fmt"
	"github.com/uhppoted/uhppote-core/uhppote"
)

var VERSION = Version{
	version: uhppote.VERSION,
}

type Version struct {
	version string
}

func (c *Version) Name() string {
	return "version"
}

func (c *Version) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("version", flag.ExitOnError)
}

func (c *Version) Execute(ctx context.Context) error {
	fmt.Printf("%v\n", c.version)

	return nil
}

func (c *Version) Description() string {
	return "Displays the current version"
}

func (c *Version) Usage() string {
	return ""
}

func (c *Version) Help() {
	fmt.Println()
	fmt.Println("  Displays the uhppoted-rest version in the format v<major>.<minor>.<build> e.g. v1.00.10")
	fmt.Println()
}
