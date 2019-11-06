package commands

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"uhppoted/config"
	"uhppoted/eventlog"
)

type Run struct {
	configuration string
	dir           string
	pidFile       string
	logFile       string
	logFileSize   int
	console       bool
	debug         bool
}

var runCmd = Run{
	configuration: "/usr/local/etc/com.github.twystd.uhppoted/uhppoted.conf",
	dir:           "/usr/local/var/com.github.twystd.uhppoted",
	pidFile:       "/usr/local/var/com.github.twystd.uhppoted/uhppoted.pid",
	logFile:       "/usr/local/var/com.github.twystd.uhppoted/logs/uhppoted.log",
	logFileSize:   10,
	console:       false,
	debug:         false,
}

func (r *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.StringVar(&r.configuration, "config", r.configuration, "Sets the configuration file path")
	flagset.StringVar(&r.dir, "dir", r.dir, "Work directory")
	flagset.StringVar(&r.pidFile, "pid", r.pidFile, "Sets the service PID file path")
	flagset.StringVar(&r.logFile, "logfile", r.logFile, "Sets the log file path")
	flagset.IntVar(&r.logFileSize, "logfilesize", r.logFileSize, "Sets the log file size before forcing a log rotate")
	flagset.BoolVar(&r.console, "console", r.console, "Writes log entries to stdout")
	flagset.BoolVar(&r.debug, "debug", r.debug, "Displays vaguely useful internal information")

	return flagset
}

func (r *Run) Execute(ctx Context) error {
	log.Printf("uhppoted daemon %s - %s (PID %d)\n", VERSION, "MacOS", os.Getpid())

	f := func(c *config.Config) error {
		return r.exec(c)
	}

	return r.execute(ctx, f)
}

func (r *Run) exec(c *config.Config) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	if !r.console {
		events := eventlog.Ticker{Filename: r.logFile, MaxSize: r.logFileSize}
		logger = log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)
		rotate := make(chan os.Signal, 1)

		signal.Notify(rotate, syscall.SIGHUP)

		go func() {
			for {
				<-rotate
				log.Printf("Rotating uhppoted log file '%s'\n", r.logFile)
				events.Rotate()
			}
		}()
	}

	r.run(c, logger)

	return nil
}
