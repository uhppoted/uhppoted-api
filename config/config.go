package config

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"uhppote/encoding/conf"
)

type DeviceMap map[uint32]*Device

type Device struct {
	Address *net.UDPAddr
	Door    []string
}

type MQTT struct {
	ServerID          string      `conf:"server.ID"`
	Broker            string      `conf:"broker"`
	BrokerCertificate string      `conf:"broker.certificate"`
	ClientCertificate string      `conf:"client.certificate"`
	ClientKey         string      `conf:"client.key"`
	Topics            Topics      `conf:"topic"`
	EventsKeyID       string      `conf:"events.key"`
	SystemKeyID       string      `conf:"system.key"`
	EventIDs          string      `conf:"events.index.filepath"`
	Permissions       Permissions `conf:"permissions"`
	HMAC              HMAC        `conf:"security.HMAC"`
	Authentication    string      `conf:"security.authentication"`
	HOTP              HOTP        `conf:"security.hotp"`
	RSA               RSA         `conf:"security.rsa"`
	Nonce             Nonce       `conf:"security.nonce"`
	SignOutgoing      bool        `conf:"security.outgoing.sign"`
	EncryptOutgoing   bool        `conf:"security.outgoing.encrypt"`
}

type Topics struct {
	Root     string `conf:"root"`
	Requests string `conf:"requests"`
	Replies  string `conf:"replies"`
	Events   string `conf:"events"`
	System   string `conf:"system"`
}

func (t *Topics) Resolve(subtopic string) string {
	if strings.HasPrefix(subtopic, "/") {
		return strings.ReplaceAll(strings.TrimPrefix(subtopic, "/"), " ", "")
	}

	if strings.HasPrefix(subtopic, "./") {
		return strings.ReplaceAll(t.Root+"/"+strings.TrimPrefix(subtopic, "./"), " ", "")
	}

	return strings.ReplaceAll(t.Root+"/"+subtopic, " ", "")
}

type HMAC struct {
	Required bool   `conf:"required"`
	Key      string `conf:"key"`
}

type HOTP struct {
	Range    uint64 `conf:"range"`
	Secrets  string `conf:"secrets"`
	Counters string `conf:"counters"`
}

type RSA struct {
	KeyDir string `conf:"keys"`
}

type Nonce struct {
	Required bool   `conf:"required"`
	Server   string `conf:"server"`
	Clients  string `conf:"clients"`
}

type Permissions struct {
	Enabled bool   `conf:"enabled"`
	Users   string `conf:"users"`
	Groups  string `conf:"groups"`
}

type Config struct {
	BindAddress         *net.UDPAddr  `conf:"bind.address"`
	BroadcastAddress    *net.UDPAddr  `conf:"broadcast.address"`
	ListenAddress       *net.UDPAddr  `conf:"listen.address"`
	HealthCheckInterval time.Duration `conf:"monitoring.healthcheck.interval"`
	WatchdogInterval    time.Duration `conf:"monitoring.watchdog.interval"`
	Devices             DeviceMap     `conf:"/^UT0311-L0x\\.([0-9]+)\\.(.*)/"`
	MQTT                `conf:"mqtt"`
}

func NewConfig() *Config {
	bind, broadcast, listen := DefaultIpAddresses()

	c := Config{
		BindAddress:      &bind,
		BroadcastAddress: &broadcast,
		ListenAddress:    &listen,
		MQTT: MQTT{
			ServerID:          "twystd-uhppoted",
			Broker:            "tcp://127.0.0.1:1883",
			BrokerCertificate: mqttBrokerCertificate,
			ClientCertificate: mqttClientCertificate,
			ClientKey:         mqttClientKey,
			Topics: Topics{
				Root:     "twystd/uhppoted/gateway",
				Requests: "./requests",
				Replies:  "./replies",
				Events:   "./events",
				System:   "./system",
			},
			EventsKeyID:     "events",
			SystemKeyID:     "system",
			SignOutgoing:    true,
			EncryptOutgoing: true,
			HMAC: HMAC{
				Required: false,
				Key:      "",
			},
			Authentication: "HOTP, RSA",
			HOTP: HOTP{
				Range:    8,
				Secrets:  hotpSecrets,
				Counters: hotpCounters,
			},
			RSA: RSA{
				KeyDir: rsaKeyDir,
			},
			Nonce: Nonce{
				Required: true,
				Server:   nonceServer,
				Clients:  nonceClients,
			},
			Permissions: Permissions{
				Enabled: false,
				Users:   users,
				Groups:  groups,
			},
			EventIDs: eventIDs,
		},

		HealthCheckInterval: 15 * time.Second,
		WatchdogInterval:    5 * time.Second,

		Devices: make(DeviceMap, 0),
	}

	return &c
}

func (c *Config) Load(path string) error {
	if path == "" {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return conf.Unmarshal(bytes, c)
}

// Ref. https://stackoverflow.com/questions/23529663/how-to-get-all-addresses-and-masks-from-local-interfaces-in-go
func DefaultIpAddresses() (net.UDPAddr, net.UDPAddr, net.UDPAddr) {
	bind := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 0,
		Zone: "",
	}

	broadcast := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 60000,
		Zone: "",
	}

	listen := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 60001,
		Zone: "",
	}

	copy(bind.IP, net.IPv4zero)
	copy(broadcast.IP, net.IPv4bcast)
	copy(listen.IP, net.IPv4zero)

	if ifaces, err := net.Interfaces(); err == nil {
	loop:
		for _, i := range ifaces {
			if addrs, err := i.Addrs(); err == nil {
				for _, a := range addrs {
					switch v := a.(type) {
					case *net.IPNet:
						if v.IP.To4() != nil && i.Flags&net.FlagLoopback == 0 {
							copy(bind.IP, v.IP.To4())
							copy(listen.IP, v.IP.To4())
							if i.Flags&net.FlagBroadcast != 0 {
								addr := v.IP.To4()
								mask := v.Mask
								binary.BigEndian.PutUint32(broadcast.IP, binary.BigEndian.Uint32(addr)|^binary.BigEndian.Uint32(mask))
							}
							break loop
						}
					}
				}
			}
		}
	}

	return bind, broadcast, listen
}

func (f DeviceMap) MarshalConf(tag string) ([]byte, error) {
	var s strings.Builder

	if len(f) > 0 {
		fmt.Fprintf(&s, "# DEVICES\n")
		for id, device := range f {
			fmt.Fprintf(&s, "UTO311-L0x.%d.address = %s\n", id, device.Address)
			for d, door := range device.Door {
				fmt.Fprintf(&s, "UTO311-L0x.%d.door.%d = %s\n", id, d+1, door)
			}
			fmt.Fprintf(&s, "\n")
		}
	}

	return []byte(s.String()), nil
}

func (f *DeviceMap) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	re := regexp.MustCompile(`^/(.*?)/$`)
	match := re.FindStringSubmatch(tag)
	if len(match) < 2 {
		return f, fmt.Errorf("Invalid 'conf' regular expression tag: %s", tag)
	}

	re, err := regexp.Compile(match[1])
	if err != nil {
		return f, err
	}

	for key, value := range values {
		match := re.FindStringSubmatch(key)
		if len(match) > 1 {
			id, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				return f, fmt.Errorf("Invalid 'testMap' key %s: %v", key, err)
			}

			d, ok := (*f)[uint32(id)]
			if !ok || d == nil {
				d = &Device{
					Door: make([]string, 4),
				}

				(*f)[uint32(id)] = d
			}

			switch match[2] {
			case "address":
				address, err := net.ResolveUDPAddr("udp", value)
				if err != nil {
					return f, fmt.Errorf("Device %v, invalid address '%s': %v", id, value, err)
				} else {
					d.Address = &net.UDPAddr{
						IP:   make(net.IP, net.IPv4len),
						Port: address.Port,
						Zone: "",
					}

					copy(d.Address.IP, address.IP.To4())
				}

			case "door.1":
				d.Door[0] = value

			case "door.2":
				d.Door[1] = value

			case "door.3":
				d.Door[2] = value

			case "door.4":
				d.Door[3] = value
			}
		}
	}

	return f, nil
}

func resolve(v string) (*net.UDPAddr, error) {
	address, err := net.ResolveUDPAddr("udp", v)
	if err != nil {
		return nil, err
	}

	addr := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: address.Port,
		Zone: "",
	}

	copy(addr.IP, address.IP.To4())

	return &addr, nil
}
