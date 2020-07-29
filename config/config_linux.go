package config

const (
	// DefaultConfig is the default file path for the uhppoted configuration file
	DefaultConfig = "/etc/uhppoted/uhppoted.conf"

	restUsers  string = "/etc/uhppoted/rest/users"
	restGroups string = "/etc/uhppoted/rest/groups"

	mqttBrokerCertificate string = "/etc/uhppoted/mqtt/broker.cert"
	mqttClientCertificate string = "/etc/uhppoted/mqtt/client.cert"
	mqttClientKey         string = "/etc/uhppoted/mqtt/client.key"
	mqttUsers             string = "/etc/uhppoted/mqtt.permissions.users"
	mqttGroups            string = "/etc/uhppoted/mqtt.permissions.groups"
	hotpSecrets           string = "/etc/uhppoted/mqtt.hotp.secrets"
	rsaKeyDir             string = "/etc/uhppoted/mqtt/rsa"

	eventIDs     string = "/var/uhppoted/mqtt.events.retrieved"
	hotpCounters string = "/var/uhppoted/mqtt.hotp.counters"
	nonceServer  string = "/var/uhppoted/mqtt.nonce"
	nonceClients string = "/var/uhppoted/mqtt.nonce.counters"

	httpdAuthDB string = "/etc/uhppoted/httpd/auth.json"
)
