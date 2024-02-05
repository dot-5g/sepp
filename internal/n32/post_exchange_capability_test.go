package n32_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/dot-5g/sepp/internal/model"
	"github.com/dot-5g/sepp/internal/n32"
)

func TestGivenSupportedCapabilityWhenHandlePostExchangeCapabilityThenReturns200(t *testing.T) {
	localFQDN := "local-sepp.example.com"
	seppContext := &model.SEPPContext{
		Mu:                          sync.Mutex{},
		LocalN32FQDN:                model.FQDN(localFQDN),
		RemoteN32FQDN:               model.FQDN(""),
		SupportedSecurityCapability: model.SecurityCapability("TLS"),
	}

	reqBody, err := json.Marshal(n32.SecNegotiateReqData{
		Sender:                     "testSender",
		SupportedSecCapabilityList: []model.SecurityCapability{model.TLS},
	})
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/n32c-handshake/v1/exchange-capability", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	n32.HandlePostExchangeCapability(rr, req, seppContext)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := n32.SecNegotiateRspData{
		Sender:                model.FQDN(localFQDN),
		SelectedSecCapability: model.TLS,
	}

	var actualResponse n32.SecNegotiateRspData
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if actualResponse != expectedResponse {
		t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", actualResponse, expectedResponse)
	}
}

func TestGivenSupportedCapabilityWhenHandlePostExchangeCapabilityThenRemoteFQDNIsStored(t *testing.T) {
	localFQDN := "local-sepp.example.com"
	remoteFQDN := "remote-sepp.example.com"
	seppContext := &model.SEPPContext{
		Mu:                          sync.Mutex{},
		LocalN32FQDN:                model.FQDN(localFQDN),
		RemoteN32FQDN:               model.FQDN(""),
		SupportedSecurityCapability: model.SecurityCapability("TLS"),
	}

	reqBody, err := json.Marshal(n32.SecNegotiateReqData{
		Sender:                     model.FQDN(remoteFQDN),
		SupportedSecCapabilityList: []model.SecurityCapability{model.TLS},
	})
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/n32c-handshake/v1/exchange-capability", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	n32.HandlePostExchangeCapability(rr, req, seppContext)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if seppContext.RemoteN32FQDN != model.FQDN(remoteFQDN) {
		t.Errorf("RemoteFQDN not stored: got %v want %v", seppContext.RemoteN32FQDN, remoteFQDN)
	}
}

func TestGivenUnsupportedCapabilityWhenHandlePostExchangeCapabilityThenReturns4xx(t *testing.T) {
	localFQDN := "local-sepp.example.com"
	seppContext := &model.SEPPContext{
		Mu:                          sync.Mutex{},
		LocalN32FQDN:                model.FQDN(localFQDN),
		RemoteN32FQDN:               model.FQDN(""),
		SupportedSecurityCapability: model.SecurityCapability("TLS"),
	}
	reqBody, err := json.Marshal(n32.SecNegotiateReqData{
		Sender:                     "testSender",
		SupportedSecCapabilityList: []model.SecurityCapability{model.ALS},
	})
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/n32c-handshake/v1/exchange-capability", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	n32.HandlePostExchangeCapability(rr, req, seppContext)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGivenUnsupportedCapabilityWhenHandlePostExchangeCapabilityThenRemoteFQDNNotStored(t *testing.T) {
	localFQDN := "local-sepp.example.com"
	seppContext := &model.SEPPContext{
		Mu:                          sync.Mutex{},
		LocalN32FQDN:                model.FQDN(localFQDN),
		RemoteN32FQDN:               model.FQDN(""),
		SupportedSecurityCapability: model.SecurityCapability("TLS"),
	}
	reqBody, err := json.Marshal(n32.SecNegotiateReqData{
		Sender:                     "testSender",
		SupportedSecCapabilityList: []model.SecurityCapability{model.ALS},
	})
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/n32c-handshake/v1/exchange-capability", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	n32.HandlePostExchangeCapability(rr, req, seppContext)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if seppContext.RemoteN32FQDN != model.FQDN("") {
		t.Errorf("RemoteFQDN stored: got %v want %v", seppContext.RemoteN32FQDN, "")
	}
}
