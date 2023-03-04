package gateway

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestCreateCompletion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", BASE_URL+"/chat/completions", func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Content-Type") == "" {
			return httpmock.NewStringResponse(400, ""), nil
		}

		if r.Header.Get("Authorization") == "" {
			return httpmock.NewStringResponse(401, ""), nil
		}

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
