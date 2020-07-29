package config

const (
	// DefaultConfig is the default file path for the uhppoted configuration file
	DefaultConfig = "/usr/local/etc/com.github.twystd.uhppoted/uhppoted.conf"

	restUsers  string = "/usr/local/etc/com.github.twystd.uhppoted/rest/users"
	restGroups string = "/usr/local/etc/com.github.twystd.uhppoted/rest/groups"

	mqttBrokerCertificate string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt/broker.cert"
	mqttClientCertificate string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt/client.cert"
	mqttClientKey         string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt/client.key"
	mqttUsers             string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt.permissions.users"
	mqttGroups            string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt.permissions.groups"
	hotpSecrets           string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt.hotp.secrets"
	rsaKeyDir             string = "/usr/local/etc/com.github.twystd.uhppoted/mqtt/rsa"

	eventIDs     string = "/usr/local/var/com.github.twystd.uhppoted/mqtt.events.retrieved"
	hotpCounters string = "/usr/local/var/com.github.twystd.uhppoted/mqtt.hotp.counters"
	nonceServer  string = "/usr/local/var/com.github.twystd.uhppoted/mqtt.nonce"
	nonceClients string = "/usr/local/var/com.github.twystd.uhppoted/mqtt.nonce.counters"

	httpdAuthDB string = "/usr/local/etc/com.github.twystd.uhppoted/httpd/auth.json"
)
