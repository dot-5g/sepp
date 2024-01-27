package main

import (
	"flag"
	"log"

	n32c "github.com/dot-5g/sepp/internal/n32"

	"github.com/dot-5g/sepp/config"

	"github.com/labstack/echo/v4"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "config.yaml", "Path to the config file")
}

func main() {
	flag.Parse()
	config, err := config.ReadConfig(configFilePath)
	if err != nil {
		log.Printf("Failed to read config file: %s\n", err)
		return
	}

	n32c := n32c.N32C{
		FQDN: n32c.FQDN(config.SEPP.FQDN),
	}

	echoServer := echo.New()

	n32cHandshakeGroup := echoServer.Group("/n32c-handshake/v1")
	n32cHandshakeGroup.POST("/exchange-capability", n32c.HandlePostExchangeCapability)

	echoServer.Logger.Fatal(echoServer.Start(":1323"))
}
