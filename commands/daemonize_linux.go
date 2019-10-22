package commands

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"
	"uhppoted/config"
)

type Daemonize struct {
	user  string
	group string
}

type data struct {
	Description      string
	Documentation    string
	Executable       string
	PID              string
	User             string
	Group            string
	Uid              int
	Gid              int
	LogFiles         []string
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
}

const serviceTemplate = `[Unit]
Description={{.Description}}
Documentation={{.Documentation}}
After=syslog.target network.target

[Service]
Type=simple
ExecStart={{.Executable}}
PIDFile={{.PID}}
User={{.User}}
Group={{.Group}}

[Install]
WantedBy=multi-user.target
`

const logRotateTemplate = `{{range .LogFiles}}{{.}} {
    daily
    rotate 30
    compress
        compresscmd /bin/bzip2
        compressext .bz2
        dateext
    missingok
    notifempty
    su uhppoted uhppoted
    postrotate
       /usr/bin/killall -HUP uhppoted
    endscript
}{{end}}
`

const confTemplate = `bind.address = {{.BindAddress}}
broadcast.address = {{.BroadcastAddress}}

# Example configuration for UTO311-L04 with serial number 305419896
# UT0311-L0x.305419896.address = 192.168.1.100:60000
# UT0311-L0x.305419896.door.1 = Front Door
# UT0311-L0x.305419896.door.2 = Side Door
# UT0311-L0x.305419896.door.3 = Garage
# UT0311-L0x.305419896.door.4 = Workshop
`

func NewDaemonize() *Daemonize {
	return &Daemonize{
		user:  "uhppoted",
		group: "uhppoted",
	}
}

func (c *Daemonize) Parse(args []string) error {
	flags := flag.NewFlagSet("daemonize", flag.ExitOnError)
	ug := flags.String("user", "uhppoted:uhppoted", "user:group for uhppoted service. Defaults to uhppoted:uhppoted")
	err := flags.Parse(args)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(\w+?):(\w+)`)
	match := re.FindStringSubmatch(*ug)

	if match == nil {
		return fmt.Errorf("Invalid user:group '%s'", *ug)
	}

	c.user = match[1]
	c.group = match[2]

	return nil
}

func (c *Daemonize) Execute(ctx Context) error {
	fmt.Println("   ... daemonizing")

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	uid, gid, err := c.getUser()
	if err != nil {
		return err
	}

	bind, broadcast, err := config.DefaultIpAddresses()
	if err != nil {
		return err
	}

	if bind == nil || broadcast == nil {
		return errors.New("Unable to determine default bind and broadcast IP addresses")
	}

	d := data{
		Description:      "UHPPOTE UTO311-L0x access card controllers service/daemon ",
		Documentation:    "https://github.com/twystd/uhppote-go",
		Executable:       executable,
		PID:              "/var/uhppoted/uhppoted.pid",
		User:             "uhppoted",
		Group:            "uhppoted",
		Uid:              uid,
		Gid:              gid,
		LogFiles:         []string{"/var/log/uhppoted/uhppoted.log"},
		BindAddress:      bind,
		BroadcastAddress: broadcast,
	}

	if err := c.systemd(&d); err != nil {
		return err
	}

	if err := c.logrotate(&d); err != nil {
		return err
	}

	if err := c.mkdirs(&d); err != nil {
		return err
	}

	if err := c.conf(&d); err != nil {
		return err
	}

	fmt.Println("   ... uhppoted registered as a systemd service")
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Println("     > sudo systemctl start uhppoted")
	fmt.Println()
	fmt.Println("   The firewall may need additional rules to allow UDP broadcast e.g. for UFW:")
	fmt.Println()
	fmt.Printf("     > sudo ufw allow from %s to any port 60000 proto udp\n", d.BindAddress.IP)
	fmt.Println()

	return nil
}

func (c *Daemonize) systemd(d *data) error {
	path := filepath.Join("/etc/systemd/system", "uhppoted.service")
	t := template.Must(template.New("uhppoted.service").Parse(serviceTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, d)
}

func (c *Daemonize) logrotate(d *data) error {
	path := filepath.Join("/etc/logrotate.d", "uhppoted")
	t := template.Must(template.New("uhppoted.logrotate").Parse(logRotateTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, d)
}

func (c *Daemonize) conf(d *data) error {
	path := filepath.Join("/etc/uhppoted", "uhppoted.conf")
	t := template.Must(template.New("uhppoted.conf").Parse(confTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	err = t.Execute(f, d)
	if err != nil {
		return err
	}

	return os.Chown(path, d.Uid, d.Gid)
}

func (c *Daemonize) mkdirs(d *data) error {
	directories := []string{
		"/var/uhppoted",
		"/var/log/uhppoted",
		"/etc/uhppoted",
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}

		if err := os.Chown(dir, d.Uid, d.Gid); err != nil {
			return err
		}
	}

	return nil
}

func (c *Daemonize) getUser() (int, int, error) {
	u, err := user.Lookup(c.user)
	if err != nil {
		return 0, 0, err
	}

	g, err := user.LookupGroup(c.group)
	if err != nil {
		return 0, 0, err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, 0, err
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return 0, 0, err
	}

	return uid, gid, nil
}

func (c *Daemonize) Cmd() string {
	return "daemonize"
}

func (c *Daemonize) Description() string {
	return "Registers uhppoted as a service/daemon"
}

func (c *Daemonize) Usage() string {
	return ""
}

func (c *Daemonize) Help() {
	fmt.Println("Usage: uhppoted daemonize")
	fmt.Println()
	fmt.Println(" Registers uhppoted as a systemd service/daemon that runs on startup")
	fmt.Println()
}
