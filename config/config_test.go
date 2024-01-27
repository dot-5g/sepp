package config_test

import (
	"strings"
	"testing"

	"github.com/dot-5g/sepp/config"
)

func TestConfig(t *testing.T) {
	testConfig := `
sepp:
  local:
    fqdn: "sepp.local"
    host: "localhost"
    port: "8080"
    tls:
      cert: "/path/to/cert"
      key: "/path/to/key"
      ca: "/path/to/ca"
  remote:
    url: "https://remote-sepp.example.com"
`

	reader := strings.NewReader(testConfig)
	conf, err := config.ReadConfig(reader)

	if err != nil {
		t.Fatalf("Failed to read config: %s", err)
	}

	if conf.SEPP.Local.N32.FQDN != "sepp.local" {
		t.Errorf("Expected FQDN 'sepp.local', got '%s'", conf.SEPP.Local.N32.FQDN)
	}

}
