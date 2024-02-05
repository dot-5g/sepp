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
		log.Printf("N32 server - invalid request body: %v", err)
		return
	}

	if reqData.Sender == "" {
		http.Error(w, "Sender is required", http.StatusBadRequest)
		log.Printf("N32 server - sender is required")
		return
	}

	if len(reqData.SupportedSecCapabilityList) == 0 {
		http.Error(w, "SupportedSecCapabilityList is required", http.StatusBadRequest)
		log.Printf("N32 server - supportedSecCapabilityList is required")
		return
	}

	containsSupportedCapability := slices.Contains(reqData.SupportedSecCapabilityList, seppContext.SecurityCapability)
	if !containsSupportedCapability {
		http.Error(w, "Bad SecurityCapability", http.StatusBadRequest)
		log.Printf("N32 server - bad SecurityCapability - Only %s is supported", seppContext.SecurityCapability)
		return
	}

	rspData := SecNegotiateRspData{
		Sender:                seppContext.LocalN32FQDN,
		SelectedSecCapability: seppContext.SecurityCapability,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(rspData)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("N32 server - failed to encode response: %v", err)
		return
	}

	seppContext.Mu.Lock()
	seppContext.RemoteN32FQDN = reqData.Sender
	seppContext.Mu.Unlock()
	log.Printf("N32 server - successfully exchanged capability %s with remote SEPP %s", rspData.SelectedSecCapability, reqData.Sender)
}
