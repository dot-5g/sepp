package model

import (
	"sync"
)

type FQDN string

type SecurityCapability string

const TLS = SecurityCapability("TLS")
const ALS = SecurityCapability("ALS")

type SEPPContext struct {
	LocalFQDN          FQDN
	RemoteFQDN         FQDN
	SecurityCapability SecurityCapability
	Mu                 sync.Mutex
}
