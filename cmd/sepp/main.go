package main

import (
	"flag"
	"net/http"
	"os"

	n32c "github.com/dot-5g/sepp/internal/n32"

	"github.com/dot-5g/sepp/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "config.yaml", "Path to the config file")
}

func main() {
	flag.Parse()
	config, err := loadConfiguration(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}
	server := initializeServer(config)
	startServer(server, config)
}

func loadConfiguration(filePath string) (*config.Config, error) {
	conf, err := config.ReadConfig(filePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func initializeServer(conf *config.Config) *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetOutput(os.Stdout)
	e.Use(middleware.Logger())
	n32c := n32c.N32C{FQDN: n32c.FQDN(conf.SEPP.FQDN)}
	n32cGroup := e.Group("/n32c-handshake/v1")
	n32cGroup.POST("/exchange-capability", n32c.HandlePostExchangeCapability)
	return e
}

func startServer(e *echo.Echo, config *config.Config) {
	address := ":" + config.SEPP.Port
	if config.SEPP.TLS.Enabled {
		if err := e.StartTLS(address, config.SEPP.TLS.Cert, config.SEPP.TLS.Key); err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	} else {
		e.Logger.Warn("TLS is disabled")
		if err := e.Start(address); err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}
}
