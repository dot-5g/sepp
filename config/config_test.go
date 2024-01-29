package config_test

import (
	"os"
	"testing"

	"github.com/dot-5g/sepp/config"
)

func TestConfig(t *testing.T) {

	file, err := os.Open("config_test.yaml")
	if err != nil {
		t.Fatalf("Failed to open config file: %s", err)
	}
	conf, err := config.ReadConfig(file)

	if err != nil {
		t.Fatalf("Failed to read config: %s", err)
	}

	if conf.SEPP.Local.N32.FQDN != "local-sepp.example.com" {
		t.Errorf("Expected FQDN 'local-sepp.example.com', got '%s'", conf.SEPP.Local.N32.FQDN)
	}

	if conf.SEPP.Local.N32.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", conf.SEPP.Local.N32.Host)
	}

	if conf.SEPP.Local.N32.Port != "1234" {
		t.Errorf("Expected port '1234', got '%s'", conf.SEPP.Local.N32.Port)
	}

	if conf.SEPP.Local.N32.TLS.Cert != "/etc/sepp/certs/server.crt" {
		t.Errorf("Expected TLS cert '/etc/sepp/certs/server.crt', got '%s'", conf.SEPP.Local.N32.TLS.Cert)
	}

	if conf.SEPP.Local.N32.TLS.Key != "/etc/sepp/certs/server.key" {
		t.Errorf("Expected TLS key '/etc/sepp/certs/server.key', got '%s'", conf.SEPP.Local.N32.TLS.Key)
	}

	if conf.SEPP.Local.N32.TLS.CA != "/etc/sepp/certs/ca.crt" {
		t.Errorf("Expected TLS CA '/etc/sepp/certs/ca.crt', got '%s'", conf.SEPP.Local.N32.TLS.CA)
	}

	if conf.SEPP.Local.SBI.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", conf.SEPP.Local.SBI.Host)
	}

	if conf.SEPP.Local.SBI.Port != "1235" {
		t.Errorf("Expected port '1235', got '%s'", conf.SEPP.Local.SBI.Port)
	}

	if conf.SEPP.Local.SBI.TLS.Cert != "/etc/sepp/certs/server.crt" {
		t.Errorf("Expected TLS cert '/etc/sepp/certs/server.crt', got '%s'", conf.SEPP.Local.SBI.TLS.Cert)
	}

	if conf.SEPP.Local.SBI.TLS.Key != "/etc/sepp/certs/server.key" {
		t.Errorf("Expected TLS key '/etc/sepp/certs/server.key', got '%s'", conf.SEPP.Local.SBI.TLS.Key)
	}

	if conf.SEPP.Local.SBI.TLS.CA != "/etc/sepp/certs/ca.crt" {
		t.Errorf("Expected TLS CA '/etc/sepp/certs/ca.crt', got '%s'", conf.SEPP.Local.SBI.TLS.CA)
	}

	if conf.SEPP.Remote.URL != "https://remote-sepp.example.com" {
		t.Errorf("Expected URL 'https://remote-sepp.example.com', got '%s'", conf.SEPP.Remote.URL)
	}

}
