package commands

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uhppote"
	"uhppoted/config"
	"uhppoted/rest"
)

var debug = flag.Bool("debug", false, "Displays vaguely useful internal information")

type Run struct {
	commands []Command
}

func (c *Run) Parse(args []string) error {
	return nil
}

func (c *Run) Cmd() string {
	return "run"
}

func (c *Run) Description() string {
	return "Runs the uhppoted daemon/service until terminated by the system service manager"
}

func (c *Run) Usage() string {
	return "uhppoted [--debug] [--config <file>] [--logfile <file>] [--logfilesize <bytes>] [--pid <file>]"
}

func (c *Run) Help() {
	fmt.Println("Usage: uhppoted <options>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --config      Configuration file path")
	fmt.Println("    --dir         Work directory")
	fmt.Println("    --logfile     Sets the log file path")
	fmt.Println("    --logfilesize Sets the log file size before forcing a log rotate")
	fmt.Println("    --pid         Sets the PID file path")
	fmt.Println("    --debug       Displays vaguely useful internal information")
	fmt.Println("    --console     (Windows only) Runs as command-line application")
	fmt.Println()
}

func execute(ctx Context) error {
	conf := config.NewConfig()
	if err := conf.Load(*configuration); err != nil {
		log.Printf("\n   WARN:  Could not load configuration (%v)\n\n", err)
	}

	if err := os.MkdirAll(*dir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("Unable to create working directory '%v': %v", *dir, err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	if err := ioutil.WriteFile(*pidFile, []byte(pid), 0644); err != nil {
		return fmt.Errorf("Unable to create pid file: %v\n", err)
	}

	defer func() {
		os.Remove(*pidFile)
	}()

	start(conf, *logfile, *logfilesize)

	return nil
}

func run(c *config.Config, logger *log.Logger) {
	// ... syscall SIG handlers

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// ... listen forever

	for {
		err := listen(c, logger, interrupt)

		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		log.Printf("exit\n")
		break
	}
}

func listen(c *config.Config, logger *log.Logger, interrupt chan os.Signal) error {
	// ... listen

	u := uhppote.UHPPOTE{
		BindAddress:      c.BindAddress,
		BroadcastAddress: c.BroadcastAddress,
		Devices:          make(map[uint32]*net.UDPAddr),
		Debug:            *debug,
	}

	for id, d := range c.Devices {
		if d.Address != nil {
			u.Devices[id] = d.Address
		}
	}

	restd := rest.RestD{
		HttpEnabled:        c.REST.HttpEnabled,
		HttpPort:           c.REST.HttpPort,
		HttpsEnabled:       c.REST.HttpsEnabled,
		HttpsPort:          c.REST.HttpsPort,
		TLSKeyFile:         c.REST.TLSKeyFile,
		TLSCertificateFile: c.REST.TLSCertificateFile,
		CACertificateFile:  c.REST.CACertificateFile,
		CORSEnabled:        c.REST.CORSEnabled,
		OpenApi: rest.OpenApi{
			Enabled:   c.OpenApi.Enabled,
			Directory: c.OpenApi.Directory,
		},
	}

	go func() {
		restd.Run(&u, logger)
	}()

	defer rest.Close()

	touched := time.Now()
	closed := make(chan struct{})

	// ... wait until interrupted/closed

	k := time.NewTicker(15 * time.Second)
	tick := time.NewTicker(5 * time.Second)

	defer k.Stop()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if err := watchdog(touched); err != nil {
				return err
			}

		case <-k.C:
			logger.Printf("... keep-alive")
			keepalive()

		case <-interrupt:
			logger.Printf("... interrupt")
			return nil

		case <-closed:
			logger.Printf("... closed")
			return errors.New("Server error")
		}
	}

	logger.Printf("... exit")
	return nil
}

func keepalive() {
	log.Printf("keep-alive")
}

func watchdog(touched time.Time) error {
	// dt := time.Since(touched)
	// now := time.Now()
	// timeout := touched.Add(IDLE)

	// if now.After(timeout) {
	// 	return errors.New(fmt.Sprintf("Channel idle for %v", dt))
	// }

	return nil
}
