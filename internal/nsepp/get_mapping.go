package nsepp

import (
	"encoding/json"
	"log"
	"net/http"
)

type FQDN string

type TelescopicMapping struct {
	TelescopicLabel string
	SeppDomain      FQDN
	ForeignFqdn     FQDN
}

// Retrieve the mapping between the FQDN in a foreign PLMN and a telescopic FQDN, or vice versa.
func HandleGetMapping(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	foreignFqdn := queryParams.Get("foreign-fqdn")
	telescopicLabel := queryParams.Get("telescopic-label")

	if (foreignFqdn != "" && telescopicLabel != "") || (foreignFqdn == "" && telescopicLabel == "") {
		http.Error(w, "Either 'foreign-fqdn' or 'telescopic-label' must be provided, but not both.", http.StatusBadRequest)
		return
	}

	var rspData TelescopicMapping

	if telescopicLabel != "" {
		rspData = TelescopicMapping{
			ForeignFqdn: FQDN("foreignFqdn"),
		}
	} else if foreignFqdn != "" {
		rspData = TelescopicMapping{
			TelescopicLabel: "telescopicLabel",
			SeppDomain:      "seppDomain",
		}
	} else {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(rspData)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v", err)
		return
	}
}
