package gateway

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/zalando/go-keyring"
)

func TestCreateCompletion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	keyring.MockInit()

	httpmock.RegisterResponder("POST", BASE_URL+"/chat/completions", func(r *http.Request) (*http.Response, error) {
		if r.Body == nil {
			return httpmock.NewStringResponse(400, ""), nil
		}
		return httpmock.NewStringResponse(200, "{}"), nil
	})

	chat := CreateChatClient()
	_, err := chat.CreateCompletion(CreateCompletionRequest{Description: "test"})
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}
