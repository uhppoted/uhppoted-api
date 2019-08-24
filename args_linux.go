package main

import (
	"flag"
)

var dir = flag.String("dir", "/var/uhppoted", "Working directory")
var logfile = flag.String("logfile", "/var/logs/uhppoted.log", "uhppoted log file")
var logfilesize = flag.Int("logfilesize", 10, "uhppoted log file size")
var pidFile = flag.String("pid", "/var/uhppoted/uhppoted.pid", "uhppoted PID file")
var useSyslog = flag.Bool("syslog", false, "Use syslog for event logging")
