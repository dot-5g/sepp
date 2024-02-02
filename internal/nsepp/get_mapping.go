package nsepp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dot-5g/sepp/internal/model"
)

type TelescopicMapping struct {
	TelescopicLabel string
	SeppDomain      model.FQDN
	ForeignFqdn     model.FQDN
}

// <Label representing FQDN from other PLMN>.<FQDN of the SEPP in the request initiating PLMN>
func generateTelescopicMapping(remoteFQDN model.FQDN, localSEPPFQDN model.FQDN) (telescopicLabel model.FQDN) {
	return model.FQDN(fmt.Sprintf("%s.%s", remoteFQDN, localSEPPFQDN))
}

// Retrieve the mapping between the FQDN in a foreign PLMN and a telescopic FQDN, or vice versa.
func HandleGetMapping(w http.ResponseWriter, r *http.Request, seppContext *model.SEPPContext) {
	queryParams := r.URL.Query()
	foreignFqdn := queryParams.Get("foreign-fqdn")
	telescopicLabel := queryParams.Get("telescopic-label")

	if (foreignFqdn != "" && telescopicLabel != "") || (foreignFqdn == "" && telescopicLabel == "") {
		http.Error(w, "Either 'foreign-fqdn' or 'telescopic-label' must be provided, but not both.", http.StatusBadRequest)
		return
	}

	var rspData TelescopicMapping

	if telescopicLabel != "" {
		foreignFQDN := generateTelescopicMapping(model.FQDN(telescopicLabel), seppContext.LocalFQDN)
		rspData = TelescopicMapping{
			ForeignFqdn: foreignFQDN,
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
