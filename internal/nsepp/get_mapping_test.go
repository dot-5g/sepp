package nsepp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dot-5g/sepp/internal/nsepp"
)

func TestGivenNoParameterWhenGETMappingThenReturn400(t *testing.T) {
	req, err := http.NewRequest("GET", "/nsepp-telescopic/v1/mapping", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	nsepp.HandleGetMapping(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func GivenBothParametersWhenGETMappingThenReturn400(t *testing.T) {
	req, err := http.NewRequest("GET", "/nsepp-telescopic/v1/mapping?foreign-fqdn=sepp.local&telescopic-label=sepp.local", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	nsepp.HandleGetMapping(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGivenForeignFQDNWhenGETMappingThenReturn200(t *testing.T) {
	req, err := http.NewRequest("GET", "/nsepp-telescopic/v1/mapping?foreign-fqdn=sepp.local", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	nsepp.HandleGetMapping(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGivenTelescopicLabelWhenGETMappingThenReturn200(t *testing.T) {
	req, err := http.NewRequest("GET", "/nsepp-telescopic/v1/mapping?telescopic-label=sepp.local", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	nsepp.HandleGetMapping(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
