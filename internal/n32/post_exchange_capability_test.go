package n32_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dot-5g/sepp/internal/n32"
)

func TestGivenTLSCapabilityWhenHandlePostExchangeCapabilityThenReturns200(t *testing.T) {
	seppFQDN := "sepp.local"
	reqBody, err := json.Marshal(n32.SecNegotiateReqData{
		Sender:                     "testSender",
		SupportedSecCapabilityList: []n32.SecurityCapability{n32.TLS},
	})
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/n32c-handshake/v1/exchange-capability", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := n32.N32C{
		FQDN: n32.FQDN(seppFQDN),
	}

	handler.HandlePostExchangeCapability(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := n32.SecNegotiateRspData{
		Sender:                n32.FQDN(seppFQDN),
		SelectedSecCapability: n32.TLS,
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

func TestGivenALSCapabilityWhenHandlePostExchangeCapabilityThenReturns4xx(t *testing.T) {
	seppFQDN := "sepp.local"
	reqBody, err := json.Marshal(n32.SecNegotiateReqData{
		Sender:                     "testSender",
		SupportedSecCapabilityList: []n32.SecurityCapability{n32.ALS},
	})
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/n32c-handshake/v1/exchange-capability", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := n32.N32C{
		FQDN: n32.FQDN(seppFQDN),
	}

	handler.HandlePostExchangeCapability(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
