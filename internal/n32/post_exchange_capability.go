package n32

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/dot-5g/sepp/internal/model"
)

type SecNegotiateReqData struct {
	Sender                     model.FQDN
	SupportedSecCapabilityList []model.SecurityCapability
}

type SecNegotiateRspData struct {
	Sender                model.FQDN
	SelectedSecCapability model.SecurityCapability
}

func HandlePostExchangeCapability(w http.ResponseWriter, r *http.Request, seppContext *model.SEPPContext) {
	reqData := new(SecNegotiateReqData)

	if err := json.NewDecoder(r.Body).Decode(reqData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Invalid request body: %v", err)
		return
	}

	if reqData.Sender == "" {
		http.Error(w, "Sender is required", http.StatusBadRequest)
		log.Printf("Sender is required")
		return
	}

	if len(reqData.SupportedSecCapabilityList) == 0 {
		http.Error(w, "SupportedSecCapabilityList is required", http.StatusBadRequest)
		log.Printf("SupportedSecCapabilityList is required")
		return
	}

	containsSupportedCapability := slices.Contains(reqData.SupportedSecCapabilityList, seppContext.SecurityCapability)
	if !containsSupportedCapability {
		http.Error(w, "Bad SecurityCapability", http.StatusBadRequest)
		log.Printf("Bad SecurityCapability - Only %s is supported", seppContext.SecurityCapability)
		return
	}

	rspData := SecNegotiateRspData{
		Sender:                seppContext.LocalFQDN,
		SelectedSecCapability: seppContext.SecurityCapability,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(rspData)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v", err)
		return
	}

	seppContext.Mu.Lock()
	seppContext.RemoteFQDN = reqData.Sender
	seppContext.Mu.Unlock()
}
