package n32_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	n32c "github.com/dot-5g/sepp/internal/n32"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"testing"
)

var secNegotiateReqData = `{
	"Sender": "test",
	"SupportedSecCapabilityList": ["TLS", "ALS"]
}`

type SecNegotiateRspData struct {
	Sender                string
	SelectedSecCapability string
}

func TestHandlePostExchangeCapability(t *testing.T) {
	echoServer := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/exchange-capability", strings.NewReader(secNegotiateReqData))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := echoServer.NewContext(request, recorder)
	seppFQDN := "Some.fqdn"

	sepp := n32c.N32C{
		FQDN: n32c.FQDN(seppFQDN),
	}

	if assert.NoError(t, sepp.HandlePostExchangeCapability(context)) {
		assert.Equal(t, http.StatusOK, recorder.Code)

		var rspData SecNegotiateRspData
		err := json.Unmarshal(recorder.Body.Bytes(), &rspData)
		assert.NoError(t, err)

		// Assert fields in rspData
		assert.Equal(t, seppFQDN, rspData.Sender)
		assert.Equal(t, "TLS", rspData.SelectedSecCapability)

	}
}
