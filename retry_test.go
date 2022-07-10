package requester

import (
	"bytes"
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/magiconair/properties/assert"
	"net/http"
	"testing"
)

func TestRequestWithExponentialRetry(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	count := 0
	httpmock.RegisterResponder("POST", "someurl", func(req *http.Request) (*http.Response, error) {

		if count < 2 {
			count++
			return httpmock.NewStringResponse(408, "error"), nil
		}
		return httpmock.NewStringResponse(200, "success"), nil

	},
	)
	bty, _ := json.Marshal("some content")

	requestFunc := func() (*http.Response, error) {
		return http.DefaultClient.Post("someurl", "application/json", bytes.NewBuffer(bty))
	}
	RequestWithBackoff(&Requester{requestFunc, "someurl", string(bty)})

	assert.Equal(t, count, 2)
}
